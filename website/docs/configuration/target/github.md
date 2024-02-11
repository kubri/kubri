---
sidebar_position: 1
sidebar_label: GitHub
---

# GitHub Target

Pushes your repositories to a GitHub repository.

## Environment variables

| Name           | Description                                                   |
| -------------- | ------------------------------------------------------------- |
| `GITHUB_TOKEN` | A personal access token for accessing your GitHub repository. |

## Configuration

| Name     | Description                                                              |
| -------- | ------------------------------------------------------------------------ |
| `type`   | Must be `github`.                                                        |
| `owner`  | Repository owner i.e. username or organisation.                          |
| `repo`   | Repository name.                                                         |
| `branch` | The git branch to push the artifacts to. Defaults to the default branch. |
| `folder` | The folder to store your artifacts in. Defaults to the repo root.        |

## Example

```yaml
target:
  type: github
  owner: my-org
  repo: my-repo
```
