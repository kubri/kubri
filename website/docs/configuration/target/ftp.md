---
sidebar_label: FTP
---

# FTP Target

Uploads your repositories to an FTP server.

## Environment variables

| Name           | Description                                         |
| -------------- | --------------------------------------------------- |
| `FTP_USER`     | The username for authenticating to your FTP server. |
| `FTP_PASSWORD` | The password for authenticating to your FTP server. |

## Configuration

| Name      | Description                                                      |
| --------- | ---------------------------------------------------------------- |
| `type`    | Must be `ftp`.                                                   |
| `address` | The FTP server host and port.                                    |
| `folder`  | The folder to store your artifacts in. Defaults to the FTP root. |
| `url`     | The public URL of your bucket.                                   |

## Example

```yaml
source:
  type: ftp
  address: ftp.example.com:21
  folder: my-folder
  url: https://downloads.example.com
```
