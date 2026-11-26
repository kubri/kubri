# Alpine Linux Repository

Generate and publish an Alpine Linux repository from your `.apk` files.

## Configuration

### `disabled`

- Type: `boolean`
- Default: `false`

Disable Alpine Linux.

### `folder`

- Type: `string`
- Default: `'apk'`

Path to the directory on your target.

### `key-name`

- Type: `string`

The name of the RSA key used to sign the metadata. Required if signing is enabled.

## Example

```yaml
apk:
  folder: alpine
  key-name: alpine@example.com
```
