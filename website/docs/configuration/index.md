---
description: How a Kubri configuration file is structured and where each section is documented.
---

# Configuration

Kubri is configured through a YAML file in your project root. It looks for `.kubri.yml`,
`.kubri.yaml`, `kubri.yml` and `kubri.yaml`, in that order.

A configuration has three parts:

- [Source](/docs/configuration/source/) defines where your releases come from. It is required.
- [Target](/docs/configuration/target/) defines where the repositories are uploaded to. It is
  required.
- [Generators](/docs/configuration/generators/) define which repositories and update feeds are
  built. Each generator has its own top-level key.

```yaml
source:
  type: github
  owner: my-org
  repo: my-repo

target:
  type: s3
  bucket: my-bucket
  region: us-east-1
  url: https://download.example.com

apt: {}
yum: {}
```

An empty object enables a generator with its default settings. See
[Getting Started](/docs/getting-started/) for a configuration that uses every generator.

## Global options

These options are set at the top level of the file.

### `version`

- Type: `string`

A version constraint limiting which releases are published. See
[Version Constraints](/docs/version-constraints/) for the syntax.

### `prerelease`

- Type: `boolean`
- Default: `false`

Publish pre-release versions, which are skipped by default.

### `upload-packages`

- Type: `boolean`
- Default: `false`

Upload packages to the target for App Installer and Sparkle instead of linking to the source. Set it
to `true` if your source is private.

### `title`

- Type: `string`

The default title of the Sparkle feed.

### `description`

- Type: `string`

The default description of the Sparkle feed.
