package s3_test

import (
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/testutils"
	"github.com/abemedia/appcast/source/blob/s3"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestBlobS3(t *testing.T) {
	host := testutils.TestContainer(t, testutils.Container{
		Image: "adobe/s3mock:latest",
		Port:  9090,
		Env:   map[string]string{"initialBuckets": "bucket"},
		Wait:  wait.ForHTTP("/").WithPort("9090").WithStatusCodeMatcher(nil),
	})

	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	dir := "downloads/test"
	repo := "bucket/" + dir

	s, err := s3.New(source.Config{Repo: repo + "?endpoint=" + host + "&disableSSL=true&s3ForcePathStyle=true"})
	if err != nil {
		t.Fatal(err)
	}

	makeURL := func(version, asset string) string {
		return "http://" + host + "/" + filepath.Join(repo, version, asset)
	}

	testutils.TestBlob(t, s, makeURL)
}
