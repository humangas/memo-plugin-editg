package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type config struct {
	MemoDir          string `toml:"memodir"`
	Editor           string `toml:"editor"`
	Column           int    `toml:"column"`
	SelectCmd        string `toml:"selectcmd"`
	GrepCmd          string `toml:"grepcmd"`
	AssetsDir        string `toml:"assetsdir"`
	PluginsDir       string `toml:"pluginsdir"`
	TemplateDirFile  string `toml:"templatedirfile"`
	TemplateBodyFile string `toml:"templatebodyfile"`
}

func (cfg *config) load() error {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "memo")
		}
		dir = filepath.Join(dir, "memo")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "memo")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}
	file := filepath.Join(dir, "config.toml")

	confDir := dir

	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		cfg.MemoDir = expandPath(cfg.MemoDir)
		cfg.AssetsDir = expandPath(cfg.AssetsDir)
		if cfg.PluginsDir == "" {
			cfg.PluginsDir = filepath.Join(confDir, "plugins")
		}
		cfg.PluginsDir = expandPath(cfg.PluginsDir)

		dir := os.Getenv("MEMODIR")
		if dir != "" {
			cfg.MemoDir = dir
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	dir = filepath.Join(dir, "_posts")
	os.MkdirAll(dir, 0700)
	cfg.MemoDir = filepath.ToSlash(dir)
	cfg.Editor = os.Getenv("EDITOR")
	if cfg.Editor == "" {
		cfg.Editor = "vim"
	}
	cfg.Column = 20
	cfg.SelectCmd = "peco"
	cfg.GrepCmd = "grep -nH ${PATTERN} ${FILES}"
	cfg.AssetsDir = "."
	dir = filepath.Join(confDir, "plugins")
	os.MkdirAll(dir, 0700)
	cfg.PluginsDir = dir

	dir = os.Getenv("MEMODIR")
	if dir != "" {
		cfg.MemoDir = dir
	}
	return toml.NewEncoder(f).Encode(cfg)
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}
