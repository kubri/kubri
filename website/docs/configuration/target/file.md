---
sidebar_label: Local Filesystem
---

# Local Filesystem Target

Stores your repositories in a local folder.

## Configuration

| Name   | Description                                                                                |
| ------ | ------------------------------------------------------------------------------------------ |
| `type` | Must be `file`.                                                                            |
| `path` | Relative or absolute path to your destination folder. Will be created if it doesn't exist. |

## Example

```yaml
source:
  type: file
  path: path/to/destination
```
