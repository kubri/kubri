---
sidebar_position: 3
---

# Target

Your target defines where the repositories are uploaded to.

## Azure Blob Storage

Uploads your repositories to an Azure Blob Storage bucket.

### Environment variables

| Name                      | Description                                                                                       |
| ------------------------- | ------------------------------------------------------------------------------------------------- |
| `AZURE_STORAGE_ACCOUNT`   | Storage account name. Must be used in conjunction with either storage account key or a SAS token. |
| `AZURE_STORAGE_KEY`       | Storage account key. Must be used in conjunction with storage account name.                       |
| `AZURE_STORAGE_SAS_TOKEN` | A Shared Access Signature (SAS). Must be used in conjunction with storage account name.           |

### Configuration

| Name     | Description                                                         |
| -------- | ------------------------------------------------------------------- |
| `type`   | Must be `azureblob`.                                                |
| `bucket` | Storage bucket name.                                                |
| `folder` | The folder to store your artifacts in. Defaults to the bucket root. |
| `url`    | The public URL of your bucket.                                      |

### Example

```yaml
target:
  type: azureblob
  bucket: my-bucket
  folder: my-folder
  url: https://downloads.example.com
```

## Google Cloud Storage

Uploads your repositories to a Google Cloud Storage bucket.

### Environment variables

| Name                             | Description                                                         |
| -------------------------------- | ------------------------------------------------------------------- |
| `GOOGLE_APPLICATION_CREDENTIALS` | The location of a credential JSON file used to authenticate to GCS. |

### Configuration

| Name     | Description                                                                           |
| -------- | ------------------------------------------------------------------------------------- |
| `type`   | Must be `gcs`.                                                                        |
| `bucket` | Storage bucket name.                                                                  |
| `folder` | The folder to store your artifacts in. Defaults to the bucket root.                   |
| `url`    | The public URL of your bucket. Defaults to `https://storage.googleapis.com/{bucket}`. |

### Example

```yaml
target:
  type: gcs
  bucket: my-bucket
  folder: my-folder
  url: https://downloads.example.com
```

## Amazon S3

Uploads your repositories to an Amazon S3 or S3-compatible bucket.

### Environment variables

| Name                    | Description                                                                                                                  |
| ----------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `AWS_ACCESS_KEY_ID`     | The access key for your AWS account or S3-compatible service.                                                                |
| `AWS_SECRET_ACCESS_KEY` | The secret key for your AWS account or S3-compatible service.                                                                |
| `AWS_SESSION_TOKEN`     | The session key for your AWS account or S3-compatible service. This is only needed when you are using temporary credentials. |

### Configuration

| Name          | Description                                                                           |
| ------------- | ------------------------------------------------------------------------------------- |
| `type`        | Must be `s3`.                                                                         |
| `bucket`      | Storage bucket name.                                                                  |
| `endpoint`    | The endpoint of your S3-compatible server. Not required if using Amazon S3.           |
| `region`      | The bucket region.                                                                    |
| `disable-ssl` | Disable SSL.                                                                          |
| `folder`      | The folder to store your artifacts in. Defaults to the bucket root.                   |
| `url`         | The public URL of your bucket. Defaults to `https://storage.googleapis.com/<bucket>`. |

### Examples

#### Amazon S3

```yaml
target:
  type: s3
  bucket: my-bucket
  region: us-east-1
  folder: my-folder
  url: https://downloads.example.com
```

#### Cloudflare R2

```yaml
target:
  type: s3
  bucket: my-bucket
  region: auto
  folder: my-folder
  endpoint: https://<accountId>.r2.cloudflarestorage.com
  url: https://downloads.example.com
```

## GitHub

Pushes your repositories to a GitHub repository.

### Environment variables

| Name           | Description                                                   |
| -------------- | ------------------------------------------------------------- |
| `GITHUB_TOKEN` | A personal access token for accessing your GitHub repository. |

### Configuration

| Name     | Description                                                              |
| -------- | ------------------------------------------------------------------------ |
| `type`   | Must be `github`.                                                        |
| `owner`  | Repository owner i.e. username or organisation.                          |
| `repo`   | Repository name.                                                         |
| `branch` | The git branch to push the artifacts to. Defaults to the default branch. |
| `folder` | The folder to store your artifacts in. Defaults to the repo root.        |

### Example

```yaml
target:
  type: github
  owner: my-org
  repo: my-repo
```

## File

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
