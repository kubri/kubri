# CLI Reference

## Options

| Flag       | Short | Default | Description            |
| ---------- | ----- | ------- | ---------------------- |
| `--help`   | `-h`  | `false` | Help for kubri.        |
| `--silent` | `-s`  | `false` | Only log fatal errors. |

## Kubri CLI Commands

### `kubri build`

generates and publishes your repositories.

#### Options

| Flag       | Short | Default                                                   | Description               |
| ---------- | ----- | --------------------------------------------------------- | ------------------------- |
| `--config` | `-c`  | `.kubri.yml` / `.kubri.yaml` / `kubri.yml` / `kubri.yaml` | Path to your config file. |

### `kubri completion (bash|fish|powershell|zsh)`

Output a shell completion script. Run `kubri completion <shell> --help` for how to load it.

### `kubri jsonschema`

Output the JSON schema for the config file.

### `kubri keys create`

Create private keys for signing update packages. If keys already exist, this is a no-op.

#### Options

| Flag      | Short | Default | Description        |
| --------- | ----- | ------- | ------------------ |
| `--email` |       |         | Email for PGP key. |
| `--name`  |       |         | Name for PGP key.  |

### `kubri keys export (dsa|ed25519|pgp|rsa)`

Output a private key. Ed25519 keys are output in PEM format.

### `kubri keys import (dsa|ed25519|pgp|rsa) <path>`

Import private keys for signing update packages. If keys already exist, this is a no-op.

#### Options

| Flag      | Short | Default | Description             |
| --------- | ----- | ------- | ----------------------- |
| `--force` | `-f`  | `false` | Overwrite existing key. |

### `kubri keys public (dsa|ed25519|pgp|rsa)`

Output public key.

### `kubri version`

Output the version of Kubri along with your OS and architecture.
