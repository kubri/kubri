---
sidebar_position: 1
---

# Introduction

Kubri is a tool which automates signing and releasing software for common package managers and
software update frameworks.

## Why we made it

Kubri was created to solve the common problem of releasing your software for multiple different
platforms. This is usually done with custom scripts, using various different tools which often even
need to run on different operating systems, making the release process complex and disjointed.

For example, to create a YUM repository you'd use the createrepo tool which only runs on Linux,
while creating an App Installer file is commonly done in Visual Studio, which only runs on Windows.
To make matters worse, many package managers require creating pull requests to public repositories
using various different file formats.

Kubri was created to simplify this process, allowing you to just write a small YAML configuration
file and then release your software for all common package managers and software update frameworks
with just a single command.

It runs on Windows, Mac & Linux with zero dependencies.

## Usage

Your release process is configured through a YAML file called `.kubri.yml`.  
Once you've set it up you can publish your new releases anytime by running the command
`kubri build`.

See the example below for a configuration that would fetch your releases from GitHub and generate an
APT repository, a YUM repository, a Windows App Installer file and a Sparkle feed and publish them
to Amazon S3.

```yaml
source:
  type: github
  owner: my-org
  repo: my-repo

target:
  type: s3
  bucket: my_bucket
  region: us-east-1
  url: https://download.example.com

# Use an empty object to enable the integration with default settings.
apt: {}
yum: {}
apk: {}
appinstaller: {}

arch:
  repo-name: my-repo

sparkle:
  title: My app feed title
  description: My app feed description
  params:
    - os: macos
      minimum-system-version: '10.13.0'
```

## Why the name "Kubri"?

The word "Kubri" derives from Arabic and translates to "bridge". We chose this name because it
perfectly represents our mission to bridge the gap between software release and distribution,
providing a seamless and efficient process.

## Need more?

If you need further help or would like to request a new feature, feel free to
[file an issue](https://github.com/kubri/kubri).
