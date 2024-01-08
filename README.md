# AppCast

[![Go Reference](https://pkg.go.dev/badge/github.com/abemedia/appcast.svg)](https://pkg.go.dev/github.com/abemedia/appcast)
[![Codecov](https://codecov.io/gh/abemedia/appcast/branch/master/graph/badge.svg)](https://codecov.io/gh/abemedia/appcast)

AppCast signs and publishes software for common package managers and software update frameworks.

## Supported platforms

- [APT](<https://en.wikipedia.org/wiki/APT_(software)>) (Debian, Ubuntu etc.)
- [YUM](<https://en.wikipedia.org/wiki/Yum_(software)>) (RHEL, Fedora, CentOS, OpenSUSE etc.)
- [APK](https://wiki.alpinelinux.org/wiki/Alpine_Package_Keeper) (Alpine Linux)
- [App Installer](https://en.wikipedia.org/wiki/App_Installer) (Windows)
- [Sparkle](https://sparkle-project.org/) / [WinSparkle](https://winsparkle.org/) (MacOS, Windows)

## Installation

### Homebrew

```sh
brew install abemedia/tap/appcast
```

### Apt

```sh
echo 'deb [trusted=yes] https://apt.fury.io/abemedia/ /' | sudo tee /etc/apt/sources.list.d/appcast.list
sudo apt update
sudo apt install appcast
```

### Yum

```sh
echo '[appcast]
name=AppCast
baseurl=https://yum.fury.io/abemedia/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/appcast.repo
sudo yum install appcast
```

### Binary

Download the latest binary from <https://github.com/abemedia/appcast/releases> and copy it to a
folder in your `$PATH`.
