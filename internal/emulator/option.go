package emulator

import "github.com/testcontainers/testcontainers-go"

// Option allows configuring the container request.
type Option func(*testcontainers.ContainerRequest)

// WithEnv sets an environment variable for the container.
func WithEnv(key, value string) Option {
	return func(r *testcontainers.ContainerRequest) {
		if r.Env == nil {
			r.Env = map[string]string{}
		}
		r.Env[key] = value
	}
}
