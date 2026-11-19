# Local Filesystem Target

Stores your repositories in a local folder.

## Configuration

| Name   | Description                                                                                          |
| ------ | ---------------------------------------------------------------------------------------------------- |
| `type` | Must be `file`.                                                                                      |
| `path` | Required. Relative or absolute path to your destination folder. Will be created if it doesn't exist. |
| `url`  | The public URL your repositories are served from, used for links in Sparkle and App Installer feeds. |

## Example

```yaml
target:
  type: file
  path: path/to/destination
  url: https://download.example.com
```
