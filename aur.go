package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const NOT_INSTALLED = "not installed"

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
		e.installver = NOT_INSTALLED
		return nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))

	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.SplitN(line, ":", 2)
		key, installver := strings.TrimSpace(slice[0]), strings.TrimSpace(slice[1])

		if key == "Version" {
			e.installver = installver
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.New("failed to scan package version")
	}

	return nil
}

func (e *entry) install() error {
	if err := os.Chdir(e.path); err != nil {
		return errors.New("failed to change dir")
	}

	if e.installver == NOT_INSTALLED {
		fmt.Printf("\nInstalling %v...\n\n", e.pkg.name())
	} else {
		fmt.Printf("\nUpdating %v...\n\n", e.pkg.name())
	}

	cmd := exec.Command("makepkg", "-sirc")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return errors.New("failed to install package")
	}

	return nil
}

type aur struct {
	path    string
	entries map[string]*entry
}

func (a *aur) start() error {
	if a.path != "" {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New("no user home dir")
	}

	a.path = filepath.Join(home, ".aur")

	if _, err := os.Stat(a.path); os.IsNotExist(err) {
		if err := os.Mkdir(a.path, os.ModePerm); err != nil {
			return errors.New("unable to create aur dir")
		}
	}

	return nil
}

func (a *aur) get() error {
	if err := a.start(); err != nil {
		return err
	}

	entries, err := os.ReadDir(a.path)
	if err != nil {
		return errors.New("failed to read aur dir")
	}

	if a.entries == nil {
		a.entries = make(map[string]*entry)
	}

	for _, dirEntry := range entries {
		if !dirEntry.IsDir() {
			continue
		}

		dirName := dirEntry.Name()

		if a.entries[dirName] != nil {
			continue
		}

		a.entries[dirName] = &entry{
			pkg:        &pkgbuild{},
			path:       filepath.Join(a.path, dirName),
			installver: "",
		}
	}

	return nil
}

func (a *aur) clone(url string) (string, error) {
	if err := os.Chdir(a.path); err != nil {
		return "", errors.New("failed to change dir")
	}

	repo, err := getUrlRepoName(url)
	if err != nil {
		return "", err
	}

	if a.entries[repo] != nil {
		return "", errors.New("repository already exists")
	}

	fmt.Println("Running `git clone` for ", url, repo)

	cmd := exec.Command("git", "clone", url)
	if err := cmd.Run(); err != nil {
		return "", errors.New("failed to clone repository")
	}

	return repo, nil
}

func getUrlRepoName(url string) (string, error) {
	regex := regexp.MustCompile(`\/([^\/]+)\.git$`)
	match := regex.FindStringSubmatch(url)

	if len(match) == 0 {
		return "", errors.New("invalid url")
	}

	name := match[1]

	return name, nil
}
