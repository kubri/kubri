---
sidebar_position: 1
sidebar_label: GitLab
---

# GitLab Source

Automatically gets your GitLab releases to generate your repositories.

## Environment variables

| Name           | Description                                                                                    |
| -------------- | ---------------------------------------------------------------------------------------------- |
| `GITLAB_TOKEN` | A personal access token for accessing your GitLab releases. Required for private repositories. |

## Configuration

| Name    | Description                                                       |
| ------- | ----------------------------------------------------------------- |
| `type`  | Must be `gitlab`.                                                 |
| `owner` | Repository owner i.e. username or organisation.                   |
| `repo`  | Repository name.                                                  |
| `url`   | The URL of your GitLab instance. Defaults to `https://gitlab.com` |

## Example

```yaml
source:
  type: gitlab
  owner: my-org
  repo: my-repo
  url: https://repo.example.com
```
