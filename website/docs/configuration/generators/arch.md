---
sidebar_position: 4
---

# Arch Linux Repository

Generate and publish an Arch Linux repository from your `.pkg.tar.zst` files.

## Configuration

### `disabled`

- Type: `boolean`
- Default: `false`

Disable Arch Linux.

### `folder`

- Type: `string`
- Default: `'arch'`

Path to the directory on your target.

### `repo-name`

- Type: `string`
- Allowed Values: Alphanumerical with dashes & underscores (`[A-Za-z0-9_-]`).

The name of the repository. Required.

## Example

```yaml
arch:
  folder: archlinux
  repo-name: example-repo
```
