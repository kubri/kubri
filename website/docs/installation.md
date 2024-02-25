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

## Docker

Run Kubri from docker.

```sh
docker run --rm -v $(pwd):/app -w /app kubri/kubri:latest kubri <command>
```

If signing releases you will also need to pass in your keys.  
See the following example for passing in keys via an environment variable.

```sh
docker run --rm -v $(pwd):/app -w /app \
  -e KUBRI_PGP_KEY="$(cat /path/to/pgp.key)" \
  -e KUBRI_RSA_KEY="$(cat /path/to/rsa.key)" \
  -e KUBRI_ED25519_KEY="$(cat /path/to/ed25519.key)" \
  -e KUBRI_DSA_KEY="$(cat /path/to/dsa.key)" \
  kubri/kubri:latest kubri build
```

Alternatively you can also use a volume to persist the keys.

```sh
# import keys
docker run --rm -v $(pwd):/app -w /app -v ~/.config/kubri -v path/to/pgp.key:/pgp.key kubri/kubri:latest kubri keys import pgp /pgp.key
docker run --rm -v $(pwd):/app -w /app -v ~/.config/kubri -v path/to/rsa.key:/rsa.key kubri/kubri:latest kubri keys import rsa /rsa.key

# build
docker run --rm -v $(pwd):/app -v ~/.config/kubri -w /app kubri/kubri:latest kubri build
```

## Binary

Download the latest binary from https://github.com/kubri/kubri/releases and copy it to a folder in
your `$PATH`.

## Build From Source

```sh
go install github.com/kubri/kubri@latest
```
