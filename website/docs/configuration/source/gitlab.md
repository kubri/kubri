# GitLab Source

Automatically gets your GitLab releases to generate your repositories.

Only releases tagged with a semantic version with a leading `v`, such as `v1.2.3`, are used. Other
releases are skipped.

## Environment variables

| Name           | Description                                                                                    |
| -------------- | ---------------------------------------------------------------------------------------------- |
| `GITLAB_TOKEN` | A personal access token for accessing your GitLab releases. Required for private repositories. |

## Configuration

| Name    | Description                                                       |
| ------- | ----------------------------------------------------------------- |
| `type`  | Must be `gitlab`.                                                 |
| `owner` | Required. Repository owner i.e. username or organisation.         |
| `repo`  | Required. Repository name.                                        |
| `url`   | The URL of your GitLab instance. Defaults to `https://gitlab.com` |

## Example

```yaml
source:
  type: gitlab
  owner: my-org
  repo: my-repo
  url: https://repo.example.com
```
