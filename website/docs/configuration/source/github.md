# GitHub Source

Automatically gets your GitHub releases to generate your repositories.

Only releases tagged with a semantic version with a leading `v`, such as `v1.2.3`, are used. Other
releases are skipped.

## Environment variables

| Name           | Description                                                                                    |
| -------------- | ---------------------------------------------------------------------------------------------- |
| `GITHUB_TOKEN` | A personal access token for accessing your GitHub releases. Required for private repositories. |

## Configuration

| Name    | Description                                               |
| ------- | --------------------------------------------------------- |
| `type`  | Must be `github`.                                         |
| `owner` | Required. Repository owner i.e. username or organisation. |
| `repo`  | Required. Repository name.                                |

## Example

```yaml
source:
  type: github
  owner: my-org
  repo: my-repo
```
