# hosts

## Installation

```bash
go install github.com/no-src/hosts/...@latest
```

## Usage

Print all host items.

```bash
$ hosts
```

Search for the specified hostname or IP, and multiple keywords are split with spaces.

```bash
$ hosts 192.168.100.1 github.com
```