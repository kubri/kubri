// Package githubclient creates GitHub API clients that wait out rate limits.
package githubclient

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-github/v83/github"
)

const (
	maxRetries  = 5
	maxWait     = 10 * time.Minute
	defaultWait = time.Minute
)

// New returns a GitHub client, authenticated if token is set, that waits and retries when rate
// limited, giving up after 5 retries or once a request would wait more than 10 minutes in total.
func New(token string) *github.Client {
	return newClient(http.DefaultTransport, token)
}

func newClient(base http.RoundTripper, token string) *github.Client {
	client := github.NewClient(&http.Client{Transport: &transport{base: base}})
	if token != "" {
		client = client.WithAuthToken(token)
	}
	client.DisableRateLimitCheck = true
	return client
}

type transport struct {
	base http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var waited time.Duration
	for retries := 0; ; retries++ {
		resp, err := t.base.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		wait, limited := retryAfter(resp)
		if !limited || retries == maxRetries || waited+wait > maxWait || (req.Body != nil && req.GetBody == nil) {
			return resp, nil
		}

		_ = resp.Body.Close()

		timer := time.NewTimer(wait)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}
		waited += wait

		if req.GetBody != nil {
			req = req.Clone(req.Context())
			if req.Body, err = req.GetBody(); err != nil {
				return nil, err
			}
		}
	}
}

// retryAfter follows https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api#exceeding-the-rate-limit
func retryAfter(resp *http.Response) (time.Duration, bool) {
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusTooManyRequests {
		return 0, false
	}

	if seconds, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil {
		return time.Duration(seconds) * time.Second, true
	}

	if resp.Header.Get("X-Ratelimit-Remaining") == "0" {
		if reset, err := strconv.ParseInt(resp.Header.Get("X-Ratelimit-Reset"), 10, 64); err == nil {
			return max(time.Until(time.Unix(reset, 0)), 0) + time.Second, true
		}
	}

	b, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(b))
	if err == nil && bytes.Contains(b, []byte("secondary rate limit")) {
		return defaultWait, true
	}

	return 0, false
}
