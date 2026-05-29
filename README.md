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

- Linux operating system,
- x86 architecture

## Installation

## Usage

## Legacy

Legacy Mode emulates selected Windows Mode flags on Linux.
The System flag is currently not mapped because Linux has no direct equivalent.

d = directory
a = not directory
r = no write bits for owner/group/others
h = basename starts with "."
s = unsupported, always false
l = symlink

## Architecture

TODO