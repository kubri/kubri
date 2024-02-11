---
sidebar_label: Local Filesystem
---

# Local Filesystem Target

Automatically gets your GitLab releases to generate your repositories.

### Configuration

| Name   | Description                                 |
| ------ | ------------------------------------------- |
| `type` | Must be `file`.                             |
| `path` | Relative or absolute path to your releases. |

### Example

```yaml
source:
  type: file
  path: path/to/releases
```
