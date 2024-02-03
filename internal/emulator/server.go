package emulator

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

// FileServer creates a new file server and returns its URL.
//
// On MacOS the host is replaced with host.docker.internal to allow docker
// containers to access the service.
// See https://docs.docker.com/desktop/networking/#i-want-to-connect-from-a-container-to-a-service-on-the-host
func FileServer(t *testing.T, dir string) string {
	t.Helper()
	s := httptest.NewServer(http.FileServer(http.Dir(dir)))
	t.Cleanup(s.Close)
	if runtime.GOOS != "darwin" {
		return s.URL
	}
	return strings.Replace(s.URL, "127.0.0.1", "host.docker.internal", 1)
}
