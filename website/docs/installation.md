---
sidebar_position: 2
---

# Installation

Kubri offers many installation methods. Check out the available methods below.

## Package Managers

### Homebrew

```sh
brew install --cask kubri/tap/kubri
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

### Pacman

```sh
curl -fsSL https://pkg.kubri.dev/arch/key.asc | sudo pacman-key --add -
sudo pacman-key --lsign-key 2565EFDCE78841C6EDCB31E3F0BA49B29B0548B4
echo '[kubri]
Server = https://pkg.kubri.dev/arch/$arch' | sudo tee -a /etc/pacman.conf
sudo pacman -Syu kubri
```

### Winget

```sh
winget install Kubri.Kubri
```

### Scoop

```sh
scoop bucket add kubri https://github.com/kubri/scoop-bucket
scoop install kubri
```

## Docker

Run Kubri from docker.

```sh
docker run --rm -v $(pwd):/app -w /app kubri/kubri <command>
```

If signing releases you will also need to pass in your keys.  
See the following example for passing in keys via environment variables.

```sh
docker run --rm -v $(pwd):/app -w /app \
  -e KUBRI_PGP_KEY \
  -e KUBRI_RSA_KEY \
  -e KUBRI_ED25519_KEY \
  -e KUBRI_DSA_KEY \
  kubri/kubri build
```

Alternatively, mount the key files and reference them by path.

```sh
docker run --rm -v $(pwd):/app -w /app \
  -v /path/to/keys:/keys:ro \
  -e KUBRI_PGP_KEY_PATH=/keys/pgp.key \
  -e KUBRI_RSA_KEY_PATH=/keys/rsa.key \
  -e KUBRI_ED25519_KEY_PATH=/keys/ed25519.key \
  -e KUBRI_DSA_KEY_PATH=/keys/dsa.key \
  kubri/kubri build
```

## Binary

Download the latest binary from https://github.com/kubri/kubri/releases and copy it to a folder in
your `$PATH`.

## Build From Source

```sh
go install github.com/kubri/kubri@latest
```
