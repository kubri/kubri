package appcast_test

import (
	"testing"

	"github.com/abemedia/appcast"
	"github.com/google/go-cmp/cmp"
)

func TestSignaturesUnmarshal(t *testing.T) {
	in := `
file1	dsa	myDsaSignature
file2	ed25519	myEdSignature
`

	want := appcast.Signatures{
		[2]string{"file1", "dsa"}:     "myDsaSignature",
		[2]string{"file2", "ed25519"}: "myEdSignature",
	}

	sig := appcast.Signatures{}
	if err := sig.UnmarshalText([]byte(in)); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, sig); diff != "" {
		t.Error(diff)
	}
}

func TestSignaturesMarshal(t *testing.T) {
	in := appcast.Signatures{
		[2]string{"file1", "dsa"}:     "myDsaSignature",
		[2]string{"file2", "ed25519"}: "myEdSignature",
	}

	b, err := in.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	sig := appcast.Signatures{}
	if err := sig.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(in, sig); diff != "" {
		t.Error(diff)
	}
}
