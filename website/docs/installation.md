---
sidebar_position: 2
---

# Installation

Kubri offers many installation methods. Check out the available methods below.

## Package Managers

### Homebrew

```sh
brew install kubri/tap/kubri
```

### APT

```sh
curl -fsSL https://pkg.kubri.dev/deb/key.asc | gpg --dearmor | sudo tee /usr/share/keyrings/kubri.gpg
echo 'deb [signed-by=/usr/share/keyrings/kubri.gpg] https://pkg.kubri.dev/deb/ stable main' | sudo tee /etc/apt/sources.list.d/kubri.list
sudo apt update
sudo apt install kubri
```

### YUM / DNF

```sh
echo '[kubri]
name=Kubri
baseurl=https://pkg.kubri.dev/rpm/
enabled=1
gpgcheck=0
repo_gpgcheck=1
gpgkey=https://pkg.kubri.dev/rpm/repodata/repomd.xml.key' | sudo tee /etc/yum.repos.d/kubri.repo

# yum
sudo yum install kubri

# dnf
sudo dnf install kubri
```

### Zypper

```sh
sudo zypper addrepo "https://pkg.kubri.dev/rpm/" kubri
sudo zypper --gpg-auto-import-keys refresh
sudo zypper install kubri
```

### APK

```sh
curl -fsSL -o /etc/apk/keys/info@kubri.dev.rsa.pub https://pkg.kubri.dev/alpine/info@kubri.dev.rsa.pub
echo 'https://pkg.kubri.dev/alpine' >> /etc/apk/repositories
apk add kubri
```

## Binary

Download the latest binary from https://github.com/kubri/kubri/releases and copy it to a folder in
your `$PATH`.

## Build From Source

```sh
go install github.com/kubri/kubri@latest
```
