package secret_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/pkg/secret"
)

func TestSecrets(t *testing.T) {
	tests := []string{"ConfigDir", "Path", "SecretPath"}

	key := "my_secret"
	data := []byte("secret data")

	dir := t.TempDir()
	for _, v := range []string{"XDG_CONFIG_HOME", "HOME", "AppData"} {
		t.Setenv(v, dir) // Override config dir.
	}
	os.MkdirAll(filepath.Join(dir, "Library", "Application Support"), 0o755) // Create MacOS/iOS config dir.

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			switch test {
			case "Path":
				t.Setenv("APPCAST_PATH", dir)
			case "SecretPath":
				t.Setenv("APPCAST_MY_SECRET_PATH", filepath.Join(dir, "my-secret-file"))
			}

			if _, err := secret.Get(key); err == nil {
				t.Fatalf("want error %q got %q", secret.ErrKeyNotFound, err)
			}

			if err := secret.Delete(key); err == nil {
				t.Fatalf("want error %q got %q", secret.ErrKeyNotFound, err)
			}

			if err := secret.Put(key, data); err != nil {
				t.Fatal(err)
			}

			if err := secret.Put(key, data); !errors.Is(err, secret.ErrKeyExists) {
				t.Fatalf("want error %q got %q", secret.ErrKeyExists, err)
			}

			got, err := secret.Get(key)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(data, got) {
				t.Fatal("should be equal")
			}

			if err := secret.Delete(key); err != nil {
				t.Fatal(err)
			}
		})
	}

	t.Run("Env", func(t *testing.T) {
		t.Setenv("APPCAST_MY_SECRET", string(data))

		got, err := secret.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, got) {
			t.Fatal("should be equal")
		}

		if err = secret.Put(key, data); !errors.Is(err, secret.ErrEnvironment) {
			t.Fatalf("want error %q got %q", secret.ErrEnvironment, err)
		}

		if err = secret.Delete(key); !errors.Is(err, secret.ErrEnvironment) {
			t.Fatalf("want error %q got %q", secret.ErrEnvironment, err)
		}
	})
}
