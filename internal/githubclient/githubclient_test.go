package githubclient_test

import (
	"context"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/go-github/v83/github"

	"github.com/kubri/kubri/internal/githubclient"
)

const (
	secondaryLimit = `{"message":"You have exceeded a secondary rate limit.","documentation_url":"https://docs.github.com/rest/using-the-rest-api/rate-limits-for-the-rest-api#about-secondary-rate-limits"}`
	primaryLimit   = `{"message":"API rate limit exceeded.","documentation_url":"https://docs.github.com/rest/using-the-rest-api/rate-limits-for-the-rest-api"}`
	forbidden      = `{"message":"Resource not accessible by integration","documentation_url":"https://docs.github.com/rest"}`
)

type response struct {
	status int
	header map[string]string
	body   string
}

type received struct {
	auth string
	body string
}

type fakeTransport struct {
	responses []response
	received  []received
}

func (f *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		var err error
		if body, err = io.ReadAll(req.Body); err != nil {
			return nil, err
		}
		_ = req.Body.Close()
	}
	f.received = append(f.received, received{auth: req.Header.Get("Authorization"), body: string(body)})

	res := f.responses[min(len(f.received), len(f.responses))-1]
	header := make(http.Header)
	for k, v := range res.header {
		header.Set(k, v)
	}

	return &http.Response{
		StatusCode: res.status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(res.body)),
		Request:    req,
	}, nil
}

//nolint:gochecknoglobals
var now = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

func TestClient(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		timeout   time.Duration
		responses []response
		wantWait  time.Duration
		wantErr   error
	}{
		{
			name:  "SecondaryRetryAfter",
			token: "token",
			responses: []response{
				{status: http.StatusForbidden, header: map[string]string{"Retry-After": "1"}, body: secondaryLimit},
				{status: http.StatusCreated, body: "{}"},
			},
			wantWait: time.Second,
		},
		{
			name: "SecondaryNoHeaders",
			responses: []response{
				{status: http.StatusTooManyRequests, body: secondaryLimit},
				{status: http.StatusCreated, body: "{}"},
			},
			wantWait: time.Minute,
		},
		{
			name: "Primary",
			responses: []response{
				{
					status: http.StatusForbidden,
					header: map[string]string{
						"X-RateLimit-Limit":     "5000",
						"X-RateLimit-Remaining": "0",
						"X-RateLimit-Reset":     strconv.FormatInt(now.Add(30*time.Second).Unix(), 10),
					},
					body: primaryLimit,
				},
				{status: http.StatusCreated, body: "{}"},
			},
			wantWait: 31 * time.Second,
		},
		{
			name: "NotRateLimited",
			responses: []response{
				{status: http.StatusForbidden, body: forbidden},
			},
			wantErr: &github.ErrorResponse{
				Message:          "Resource not accessible by integration",
				DocumentationURL: "https://docs.github.com/rest",
			},
		},
		{
			name: "SecondaryWaitTooLong",
			responses: []response{
				{status: http.StatusForbidden, header: map[string]string{"Retry-After": "601"}, body: secondaryLimit},
			},
			wantErr: &github.AbuseRateLimitError{
				Message:    "You have exceeded a secondary rate limit.",
				RetryAfter: new(601 * time.Second),
			},
		},
		{
			name: "PrimaryWaitTooLong",
			responses: []response{
				{
					status: http.StatusForbidden,
					header: map[string]string{
						"X-RateLimit-Limit":     "5000",
						"X-RateLimit-Remaining": "0",
						"X-RateLimit-Reset":     strconv.FormatInt(now.Add(601*time.Second).Unix(), 10),
					},
					body: primaryLimit,
				},
			},
			wantErr: &github.RateLimitError{
				Rate:    github.Rate{Limit: 5000, Reset: github.Timestamp{Time: now.Add(601 * time.Second)}},
				Message: "API rate limit exceeded.",
			},
		},
		{
			name: "TooManyRetries",
			responses: slices.Repeat([]response{
				{status: http.StatusForbidden, header: map[string]string{"Retry-After": "0"}, body: secondaryLimit},
			}, 6),
			wantErr: &github.AbuseRateLimitError{
				Message:    "You have exceeded a secondary rate limit.",
				RetryAfter: new(time.Duration(0)),
			},
		},
		{
			name:    "Cancelled",
			timeout: time.Second,
			responses: []response{
				{status: http.StatusForbidden, header: map[string]string{"Retry-After": "60"}, body: secondaryLimit},
			},
			wantWait: time.Second,
			wantErr:  context.DeadlineExceeded,
		},
	}

	opts := cmp.Options{
		cmpopts.IgnoreFields(github.ErrorResponse{}, "Response"),
		cmpopts.IgnoreFields(github.AbuseRateLimitError{}, "Response"),
		cmpopts.IgnoreFields(github.RateLimitError{}, "Response"),
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := t.Context()
				if tc.timeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tc.timeout)
					defer cancel()
				}

				base := &fakeTransport{responses: tc.responses}
				client := githubclient.NewClient(base, tc.token)

				start := time.Now()
				_, _, err := client.Repositories.Create(ctx, "", &github.Repository{Name: new("repo")})

				if diff := cmp.Diff(tc.wantErr, err, opts); diff != "" {
					t.Errorf("error (-want +got):\n%s", diff)
				}
				if got := time.Since(start); got != tc.wantWait {
					t.Errorf("waited %s, want %s", got, tc.wantWait)
				}
				if len(base.received) != len(tc.responses) {
					t.Errorf("calls = %d, want %d", len(base.received), len(tc.responses))
				}

				want := received{body: `{"name":"repo"}` + "\n"}
				if tc.token != "" {
					want.auth = "Bearer " + tc.token
				}
				for i, got := range base.received {
					if got != want {
						t.Errorf("request %d = %+v, want %+v", i, got, want)
					}
				}
			})
		})
	}
}

func TestClientRateLimitCheckDisabled(t *testing.T) {
	base := &fakeTransport{responses: []response{
		{
			status: http.StatusOK,
			header: map[string]string{
				"X-RateLimit-Limit":     "5000",
				"X-RateLimit-Remaining": "0",
				"X-RateLimit-Reset":     strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10),
			},
			body: "{}",
		},
	}}
	client := githubclient.NewClient(base, "")

	for range 2 {
		if _, _, err := client.Users.Get(t.Context(), ""); err != nil {
			t.Fatal(err)
		}
	}

	if len(base.received) != 2 {
		t.Errorf("calls = %d, want 2", len(base.received))
	}
}
