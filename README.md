# GoGet-ChildItem

Linux clone of Powershell Get-ChildItem commandlet.

## Design

It is designed to work on Linux, and produce similar output as Get-ChildItem.
Because Windows and Linux are not entirely compatible, it supports only those
features that can be supported by linux. With few extra things like some of the
Attributes, that exists and can be derived from item stats. For example,
"Hidden" attribute can be supported because on Unix it means that item name
starts from dot.

## Requirements

- Linux operating system
- x86 architecture

> **Why no Windows?**
>
> Because Windows already has the original Get-ChildItem in the Powershell.
> Its implementation is far superior the this one, mostly because it uses
> native system libraries, and has entire Powershell ecosystem that backs it up.
> Therefore I do not see absolutely a reason why anyone on Windows would like
> to use the clone that does not support Powershell features and not the
> original Get-ChildItem.

### Option 1: Using Go

Requires Go 1.22 or newer.

```bash
go install github.com/ajwalkiewicz/ggci/cmd/ggci@latest
```

Make sure your Go binary directory is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify:
```bash
ggci -Version
```

### Option 2: Download binary

Download the latest binary from the Releases page:

https://github.com/ajwalkiewicz/ggci/releases/latest/

Or by using command line.
```bash
curl -L -o ggci https://github.com/OWNER/REPO/releases/latest/download/ggci-linux-amd64
chmod +x ggci
sudo mv ggci /usr/local/bin/
```

Verify:
```bash
ggci -Version
```

## Usage

## Legacy

Legacy Mode emulates selected Windows Mode flags on Linux.
The System flag is currently not mapped because Linux has no direct equivalent.

- d = directory
- a = not directory
- r = no write bits for owner/group/others
- h = basename starts with "."
- s = unsupported, always false
- l = symlink

## Architecture

TODO