---
sidebar_position: 1
sidebar_label: Google Cloud Storage
---

# Google Cloud Storage Target

Uploads your repositories to a Google Cloud Storage bucket.

## Environment variables

| Name                             | Description                                                         |
| -------------------------------- | ------------------------------------------------------------------- |
| `GOOGLE_APPLICATION_CREDENTIALS` | The location of a credential JSON file used to authenticate to GCS. |

## Configuration

| Name     | Description                                                                           |
| -------- | ------------------------------------------------------------------------------------- |
| `type`   | Must be `gcs`.                                                                        |
| `bucket` | Storage bucket name.                                                                  |
| `folder` | The folder to store your artifacts in. Defaults to the bucket root.                   |
| `url`    | The public URL of your bucket. Defaults to `https://storage.googleapis.com/{bucket}`. |

## Example

```yaml
target:
  type: gcs
  bucket: my-bucket
  folder: my-folder
  url: https://downloads.example.com
```
