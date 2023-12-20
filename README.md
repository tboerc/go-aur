# go-aur

![Github license](https://img.shields.io/github/license/tboerc/go-aur)

`go-aur` as its name suggests, is a CLI tool written in [Go](https://golang.org/) for managing AUR packages. I created this for personal use and to practice Go, but feel free to give it a try.

## How it works?

This tool was created to automate the manual process I typically use for managing my few AUR packages into a CLI tool. All the packages will be cloned into a hidden folder called `.aur` inside your home directory. You can install, list, update, and remove them.

When executing `pacman`` commands, sometimes it will prompt the user to proceed or not. These questions don't display correctly when running the command through Go. Simply pressing ENTER continues with the command execution.

## Installation

- Use the Go command tool: `go install github.com/tboerc/go-aur@latest`.
- Alternatively, download the executable from the [Releases](https://github.com/tboerc/go-aur/releases) and add it to your PATH.

## Reference

### install

```bash
go-aur install AUR_GIT_URL
```

Clones into `~/.aur` and installs the package along with necessary dependencies.

### update

```bash
go-aur update
```

Updates every package inside `~/.aur`.

### list

```bash
go-aur list
```

Lists every package inside `~/.aur` along with its current version.


### remove

```bash
go-aur PACKAGE_NAME
```

Uninstalls the package from the system and deletes its folder from `~/.aur`
