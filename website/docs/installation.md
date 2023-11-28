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
curl -fsSL https://dl.kubri.dev/deb/key.asc | gpg --dearmor | sudo tee /usr/share/keyrings/kubri.gpg
echo 'deb [signed-by=/usr/share/keyrings/kubri.gpg] https://dl.kubri.dev/deb/ /' | sudo tee /etc/apt/sources.list.d/kubri.list
sudo apt update
sudo apt install kubri
```

### YUM / DNF

```sh
echo '[kubri]
name=Kubri
baseurl=https://dl.kubri.dev/rpm/
enabled=1
gpgcheck=0
repo_gpgcheck=1
gpgkey=https://dl.kubri.dev/rpm/repodata/repomd.xml.key' | sudo tee /etc/yum.repos.d/kubri.repo

# yum
sudo yum install kubri

# dnf
sudo dnf install kubri
```

### Zypper

```sh
sudo zypper addrepo "https://dl.kubri.dev/rpm/" kubri
sudo zypper --gpg-auto-import-keys refresh
sudo zypper install kubri
```

## Binary

Download the latest binary from https://github.com/kubri/kubri/releases and copy it to a folder in
your `$PATH`.

## Build From Source

```sh
go install github.com/kubri/kubri@latest
```
