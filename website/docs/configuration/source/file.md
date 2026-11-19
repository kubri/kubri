# Local Filesystem Source

Generate your repositories from local files.

Expects files to be sorted into folders named after the version e.g. `./v1.0.0/my-app.dmg`. Folder
names must be semantic versions with a leading `v`. Other folders are skipped.

## Configuration

| Name   | Description                                                                                              |
| ------ | -------------------------------------------------------------------------------------------------------- |
| `type` | Must be `file`.                                                                                          |
| `path` | Required. Relative or absolute path to your releases.                                                    |
| `url`  | The public URL your releases are served from, used for package links in Sparkle and App Installer feeds. |

## Example

```yaml
source:
  type: file
  path: path/to/releases
  url: https://download.example.com
```
