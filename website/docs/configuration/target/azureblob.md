---
sidebar_position: 1
sidebar_label: Azure Blob Storage
---

# Azure Blob Storage Target

Uploads your repositories to an Azure Blob Storage bucket.

## Environment variables

Kubri authenticates using the environment variables below. If none of the key, SAS token or
connection string is set, it falls back to
[`DefaultAzureCredential`](https://learn.microsoft.com/en-us/azure/developer/go/sdk/authentication/credential-chains).

| Name                              | Description                                                                             |
| --------------------------------- | --------------------------------------------------------------------------------------- |
| `AZURE_STORAGE_ACCOUNT`           | Storage account name. Required unless a connection string is set.                       |
| `AZURE_STORAGE_KEY`               | Storage account key. Must be used in conjunction with storage account name.             |
| `AZURE_STORAGE_SAS_TOKEN`         | A Shared Access Signature (SAS). Must be used in conjunction with storage account name. |
| `AZURE_STORAGE_CONNECTION_STRING` | A storage account connection string. Used when no key or SAS token is set.              |

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
