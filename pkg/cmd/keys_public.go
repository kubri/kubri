package cmd

import (
	"github.com/spf13/cobra"

	"github.com/kubri/kubri/pkg/crypto/dsa"
	"github.com/kubri/kubri/pkg/crypto/ed25519"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/crypto/rsa"
	"github.com/kubri/kubri/pkg/secret"
)

func keysPublicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "public (dsa|ed25519|pgp|rsa)",
		Short:     "Output public key",
		Aliases:   []string{"p"},
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"dsa", "ed25519", "pgp", "rsa"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				pub []byte
				err error
			)
			switch args[0] {
			case "dsa":
				pub, err = getPublicKey("dsa_key", dsa.UnmarshalPrivateKey, dsa.Public, dsa.MarshalPublicKey)
			case "ed25519":
				pub, err = getPublicKey("ed25519_key", ed25519.UnmarshalPrivateKey, ed25519.Public, ed25519.MarshalPublicKey)
			case "pgp":
				pub, err = getPublicKey("pgp_key", pgp.UnmarshalPrivateKey, pgp.Public, pgp.MarshalPublicKey)
			case "rsa":
				pub, err = getPublicKey("rsa_key", rsa.UnmarshalPrivateKey, rsa.Public, rsa.MarshalPublicKey)
			}
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(pub)
			return err
		},
	}

	return cmd
}

func getPublicKey[Private, Public any](
	name string,
	unmarshal func([]byte) (Private, error),
	public func(Private) Public,
	marshal func(Public) ([]byte, error),
) ([]byte, error) {
	priv, err := secret.Get(name)
	if err != nil {
		return nil, err
	}
	key, err := unmarshal(priv)
	if err != nil {
		return nil, err
	}
	return marshal(public(key))
}
