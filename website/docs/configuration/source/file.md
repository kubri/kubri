---
sidebar_label: Local Filesystem
---

# Local Filesystem Source

Generate your repositories from local files.

Expects files to be sorted into folders named after the version e.g. `./v1.0.0/my-app.dmg`

## Configuration

| Name   | Description                                 |
| ------ | ------------------------------------------- |
| `type` | Must be `file`.                             |
| `path` | Relative or absolute path to your releases. |

## Example

```yaml
source:
  type: file
  path: path/to/releases
```
