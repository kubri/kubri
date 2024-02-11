---
sidebar_position: 3
---

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

### `kubri keys create`

Create private keys for signing update packages. If keys already exist, this is a no-op.

#### Options

| Flag      | Short | Default | Description        |
| --------- | ----- | ------- | ------------------ |
| `--email` |       |         | Email for PGP key. |
| `--name`  |       |         | Name for PGP key.  |

### `kubri keys import (dsa|ed25519|pgp) <path>`

Import private keys for signing update packages. If keys already exist, this is a no-op.

### `kubri keys public (dsa|ed25519|pgp)`

Output public key.
