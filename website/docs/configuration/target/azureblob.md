---
sidebar_position: 1
sidebar_label: Azure Blob Storage
---

# Azure Blob Storage Target

Uploads your repositories to an Azure Blob Storage bucket.

## Environment variables

| Name                      | Description                                                                                       |
| ------------------------- | ------------------------------------------------------------------------------------------------- |
| `AZURE_STORAGE_ACCOUNT`   | Storage account name. Must be used in conjunction with either storage account key or a SAS token. |
| `AZURE_STORAGE_KEY`       | Storage account key. Must be used in conjunction with storage account name.                       |
| `AZURE_STORAGE_SAS_TOKEN` | A Shared Access Signature (SAS). Must be used in conjunction with storage account name.           |

## Configuration

| Name     | Description                                                         |
| -------- | ------------------------------------------------------------------- |
| `type`   | Must be `azureblob`.                                                |
| `bucket` | Storage bucket name.                                                |
| `folder` | The folder to store your artifacts in. Defaults to the bucket root. |
| `url`    | The public URL of your bucket.                                      |

## Example

```yaml
target:
  type: azureblob
  bucket: my-bucket
  folder: my-folder
  url: https://download.example.com
```
