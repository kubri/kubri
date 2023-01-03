package cmd_test

import (
	"bytes"
	"encoding/pem"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/secret"
)

func TestKeysCreateCmd(t *testing.T) {
	t.Setenv("APPCAST_PATH", t.TempDir())

	err := cmd.Execute("", []string{"keys", "create"})
	if err != nil {
		t.Fatal(err)
	}

	dsaKey, err := secret.Get("dsa_key")
	if err != nil {
		t.Fatalf("should have created dsa key: %s", err)
	}
	p, _ := pem.Decode(dsaKey)
	_, err = dsa.UnmarshalPrivateKey(p.Bytes)
	if err != nil {
		t.Fatalf("should be valid dsa key: %s", err)
	}

	edKey, err := secret.Get("ed25519_key")
	if err != nil {
		t.Fatalf("should have created ed25519 key: %s", err)
	}
	p, _ = pem.Decode(edKey)
	_, err = ed25519.UnmarshalPrivateKey(p.Bytes)
	if err != nil {
		t.Fatalf("should be valid ed25519 key: %s", err)
	}

	// Run again to ensure existing keys aren't overwritten.
	err = cmd.Execute("", []string{"keys", "create"})
	if err != nil {
		t.Fatal(err)
	}

	if k, _ := secret.Get("dsa_key"); !bytes.Equal(dsaKey, k) {
		t.Fatal("should not have regenerated dsa key")
	}

	if k, _ := secret.Get("ed25519_key"); !bytes.Equal(edKey, k) {
		t.Fatal("should not have regenerated ed25519 key")
	}
}
