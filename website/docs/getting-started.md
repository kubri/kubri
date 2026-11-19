---
description: Write a config, create signing keys, publish your repositories and run Kubri in CI.
---

# Getting Started

This page walks you through publishing your releases as signed repositories and update feeds and
running the same setup in CI.

Before you begin, [install Kubri](/docs/installation/).

## Write a config

Create `.kubri.yml` in your project root. Kubri also looks for `.kubri.yaml`, `kubri.yml` and
`kubri.yaml`, in that order.

A config needs a `source` that your releases are fetched from, a `target` that your repositories are
published to, and at least one generator. Top-level settings such as `version` and `prerelease` are
optional and apply to every generator.

```yaml
# version is a version constraint limiting which releases are published.
version: '>=1.0.0'

# prerelease publishes pre-release versions, which are skipped by default.
prerelease: false

# upload-packages uploads packages to the target for appinstaller and sparkle
# instead of linking to the source. Set it to true if your source is private.
upload-packages: true

# source contains the configuration for the source of your releases.
source:
  type: github
  owner: my-org
  repo: my-repo

# target contains the configuration for the target of your releases.
target:
  type: s3
  bucket: my-bucket
  region: us-west-1
  url: https://download.example.com

# apk contains the configuration for the Alpine Linux repository generator.
apk:
  folder: alpine
  key-name: alpine@example.com

# apt contains the configuration for the Debian repository generator.
apt:
  folder: deb
  compress: [gzip, xz]

# arch contains the configuration for the Arch Linux repository generator.
arch:
  folder: archlinux
  repo-name: my-repo

# yum contains the configuration for the YUM repository generator.
yum:
  folder: rpm

# appinstaller contains the configuration for the AppInstaller generator.
appinstaller:
  folder: windows
  on-launch:
    hours-between-update-checks: 12
    show-prompt: true
    update-blocks-activation: true
  automatic-background-task: true
  force-update-from-any-version: true

# sparkle contains the configuration for the Sparkle appcast feed generator.
sparkle:
  folder: appcast
  title: My app feed title
  description: My app feed description
  params:
    - os: windows
      installer-arguments: /passive
    - os: macos
      minimum-system-version: '10.13.0'
```

The pages under Configuration in the sidebar cover every option for each source, target and
generator.

## Create signing keys

```sh
kubri keys create --name "My Company" --email releases@example.com
```

This creates a PGP, an RSA, an Ed25519 and a DSA key. The PGP key needs a name, an email or both.
Keys that already exist are left unchanged.

| Key          | Signs                            |
| ------------ | -------------------------------- |
| PGP          | APT, YUM and Pacman repositories |
| RSA          | APK repositories                 |
| Ed25519, DSA | Sparkle update feeds             |

When an RSA key exists, `apk` requires `key-name`, which sets the name of the public key file.

Kubri stores keys in its config directory: `~/.config/kubri` on Linux,
`~/Library/Application Support/kubri` on macOS and `%AppData%\kubri` on Windows. Set `KUBRI_PATH` to
use a different directory. To use keys you already have, import them with `kubri keys import`. See
the [CLI reference](/docs/cli/).

## Publish

```sh
kubri build
```

Kubri fetches your releases from the source, builds each configured repository and feed in parallel
and uploads them to the target. Run it after each release to publish the new version.

Sources and targets that need credentials read them from environment variables, such as
`GITHUB_TOKEN` for GitHub or `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` for S3. Each source and
target page lists the variables it uses.

## Run it in CI

In CI, store each key as a secret and pass it in a `KUBRI_<TYPE>_KEY` variable, such as
`KUBRI_PGP_KEY`. `kubri keys export <type>` prints the value to store. On GitHub Actions, add the
Kubri action after the step that creates your release:

```yaml
- uses: kubri/kubri-action@v1
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
    AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    KUBRI_PGP_KEY: ${{ secrets.KUBRI_PGP_KEY }}
    KUBRI_RSA_KEY: ${{ secrets.KUBRI_RSA_KEY }}
    KUBRI_ED25519_KEY: ${{ secrets.KUBRI_ED25519_KEY }}
```

On other CI systems, [install Kubri](/docs/installation/) or use the `kubri/kubri` Docker image, and
run `kubri build` with the same variables.

## Where to go next

- The Configuration section in the sidebar, for every source, target and generator option
- [CLI reference](/docs/cli/), to look up a command or flag
- [Version constraints](/docs/version-constraints/), for the version syntax used in the config
