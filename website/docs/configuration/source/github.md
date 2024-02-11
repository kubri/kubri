---
sidebar_position: 1
sidebar_label: GitHub
---

# GitHub Source

Automatically gets your GitHub releases to generate your repositories.

## Environment variables

| Name           | Description                                                                                    |
| -------------- | ---------------------------------------------------------------------------------------------- |
| `GITHUB_TOKEN` | A personal access token for accessing your GitHub releases. Required for private repositories. |

## Configuration

| Name    | Description                                     |
| ------- | ----------------------------------------------- |
| `type`  | Must be `github`.                               |
| `owner` | Repository owner i.e. username or organisation. |
| `repo`  | Repository name.                                |

## Example

```yaml
source:
  type: github
  owner: my-org
  repo: my-repo
```
