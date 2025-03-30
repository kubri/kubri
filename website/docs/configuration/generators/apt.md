---
sidebar_position: 0
---

# APT Repository

Generate and publish an APT repository from your `.deb` files.

## Configuration

### `disabled`

- Type: `boolean`
- Default: `false`

Disable APT.

### `folder`

- Type: `string`
- Default: `'apt'`

Path to the directory on your target.

### `compress`

- Type: `string[]`
- Default: `['gzip', 'xz']`
- Allowed Values: `'gzip'`, `'bzip2'`, `'xz'`, `'lzma'`, `'zstd'`

Compression algorithms to compress your package metadata with.

## Example

```yaml
apt:
  folder: deb
  compress:
    - gzip
    - bzip2
    - xz
    - lzma
    - zstd
```
