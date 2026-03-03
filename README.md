# sshconfctl

CLI to manage `~/.ssh/config` quickly for tools like lazyssh.

## Features

- `add` interactive host creation with strict alias/host validation
- `list` host aliases
- `edit <host>` update existing host
- `remove <host>` with safer confirmation flow
- `backup` timestamped backup of SSH config
- `format` normalize `~/.ssh/config` with clean 2-space indentation
- `version` show build/version

## Install

### From release (recommended)

Download latest Linux binary from GitHub Releases:

```bash
curl -L -O https://github.com/chalvinwz/sshconfctl/releases/download/v0.1.0/sshconfctl-linux-amd64
chmod +x sshconfctl-linux-amd64
sudo install -m 0755 sshconfctl-linux-amd64 /usr/local/bin/sshconfctl
```

### Build from source

```bash
go build -o sshconfctl ./cmd/sshconfctl
```

or:

```bash
make build
```

For a local release build with embedded version metadata:

```bash
make release
./sshconfctl version
```

## Docker usage

Build image:

```bash
docker build -t sshconfctl:latest .
```

Run with mounted SSH/config directories:

```bash
docker run --rm -it \
  -v ~/.ssh:/root/.ssh \
  -v ~/.config/sshconfctl:/root/.config/sshconfctl \
  sshconfctl:latest list
```

Replace `list` with any command (`add`, `edit`, `remove`, `format`, `version`, etc).

## Usage

```bash
./sshconfctl add
./sshconfctl list
./sshconfctl edit my-host
./sshconfctl remove my-host
./sshconfctl backup
./sshconfctl format
./sshconfctl version
```

## Config

Defaults are stored in:

`~/.config/sshconfctl/config.yaml`

```yaml
defaults:
  user: your.username
  identity_file: ~/.ssh/your-key
```

## Development

```bash
make fmt
make test
```

Override version manually (optional):

```bash
make release VERSION=v0.1.0
```
