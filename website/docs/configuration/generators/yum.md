---
sidebar_position: 0
---

# YUM Repository

Generate and publish a YUM repository from your `.rpm` files.

## Configuration

### `disabled`

- Type: `boolean`
- Default: `false`

Disable YUM.

### `folder`

- Type: `string`
- Default: `'yum'`

Path to the directory on your target.

## Example

```yaml
yum:
  folder: rpm
```
