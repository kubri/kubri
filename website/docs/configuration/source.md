---
sidebar_position: 3
---

# Source

Your source defines where the releases come from.

## GitHub

Automatically gets your GitHub releases to generate your repositories.

### Environment variables

| Name           | Description                                                                                    |
| -------------- | ---------------------------------------------------------------------------------------------- |
| `GITHUB_TOKEN` | A personal access token for accessing your GitHub releases. Required for private repositories. |

### Configuration

| Name    | Description                                     |
| ------- | ----------------------------------------------- |
| `type`  | Must be `gitlab`.                               |
| `owner` | Repository owner i.e. username or organisation. |
| `repo`  | Repository name.                                |

### Example

```yaml
source:
  type: github
  owner: my-org
  repo: my-repo
```

## GitLab

Automatically gets your GitLab releases to generate your repositories.

### Environment variables

| Name           | Description                                                                                    |
| -------------- | ---------------------------------------------------------------------------------------------- |
| `GITLAB_TOKEN` | A personal access token for accessing your GitLab releases. Required for private repositories. |

### Configuration

| Name    | Description                                                       |
| ------- | ----------------------------------------------------------------- |
| `type`  | Must be `gitlab`.                                                 |
| `owner` | Repository owner i.e. username or organisation.                   |
| `repo`  | Repository name.                                                  |
| `url`   | The URL of your GitLab instance. Defaults to `https://gitlab.com` |

### Example

```yaml
source:
  type: gitlab
  owner: my-org
  repo: my-repo
  url: https://repo.example.com
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
