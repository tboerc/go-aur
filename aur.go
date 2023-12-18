package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type pkgbuild struct {
	_pkgname string
	pkgname  string
	pkgver   string
	pkgrel   string
	epoch    string
}

func (p *pkgbuild) get(path string) error {
	file, err := os.Open(filepath.Join(path, "PKGBUILD"))
	if err != nil {
		return errors.New("failed to open pkgbuild")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.Split(line, "=")

		switch slice[0] {
		case "_pkgname":
			p._pkgname = slice[1]
		case "pkgname":
			p.pkgname = slice[1]
		case "pkgver":
			p.pkgver = slice[1]
		case "pkgrel":
			p.pkgrel = slice[1]
		case "epoch":
			p.epoch = slice[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.New("failed to scan pkgbuild")
	}

	return nil
}

func (p *pkgbuild) name() (name string) {
	return strings.Replace(p.pkgname, "${_pkgname}", p._pkgname, 1)
}

func (p *pkgbuild) version() (version string) {
	if p.epoch != "" {
		version = fmt.Sprintf("%v:%v-%v", p.epoch, p.pkgver, p.pkgrel)
	} else {
		version = fmt.Sprintf("%v-%v", p.pkgver, p.pkgrel)
	}
	return
}

type entry struct {
	pkg        *pkgbuild
	path       string
	installver string
}

func (e *entry) pull() error {
	if err := os.Chdir(e.path); err != nil {
		return errors.New("failed to change dir")
	}

	fmt.Println("Running `git pull` inside ", e.path)

	cmd := exec.Command("git", "pull")
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to check for update")
	}

	return nil
}

func (e *entry) get() error {
	if err := e.pkg.get(e.path); err != nil {
		return err
	}

	cmd := exec.Command("pacman", "-Qi", e.pkg.name())
	out, err := cmd.Output()
	if err != nil {
		return errors.New("failed to check package info")
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))

	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.SplitN(line, ":", 2)
		key, installver := strings.Trim(slice[0], " "), strings.Trim(slice[1], " ")

		if key == "Version" {
			e.installver = installver
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.New("failed to query package")
	}

	return nil
}

func (e *entry) update() error {
	if err := os.Chdir(e.path); err != nil {
		return errors.New("failed to change dir")
	}

	fmt.Printf("\nUpdating %v...\n\n", e.pkg.name())

	cmd := exec.Command("makepkg", "-sirc")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return errors.New("failed to update package")
	}

	return nil
}

type aur struct {
	entries []*entry
}

func (a *aur) get() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New("no user home dir")
	}

	path := filepath.Join(home, ".aur")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return errors.New("unable to create aur dir")
		}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return errors.New("failed to read aur dir")
	}

	for _, dirEntry := range entries {
		if !dirEntry.IsDir() {
			continue
		}

		a.entries = append(a.entries, &entry{
			pkg:        &pkgbuild{},
			path:       filepath.Join(path, dirEntry.Name()),
			installver: "",
		})
	}

	return nil
}
