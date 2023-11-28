---
sidebar_position: 4
---

# Alpine Linux Repository

Generate and publish an Alpine Linux repository from your `.apk` files.

## Configuration

### `disabled`

- Type: `boolean`
- Default: `false`

Disable Alpine Linux.

### `folder`

- Type: `string`
- Default: `apt`

Path to the directory on your target.

### `key-name`

- Type: `string`
- Default: ``

The name of the ed25519 key used to sign the metadata. Required if signing is enabled.

## Example

```yaml
apt:
  folder: alpine
  key-name: alpine@example.com
```
