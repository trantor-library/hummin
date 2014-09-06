package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"code.google.com/p/gcfg"
)

type configrc struct {
	Global struct {
		Downloads string
		Lang      string
		Num       int
	}
}

type Config struct {
	Downloads string
	Lang      string
	Num       int
}

func ParseConfig(path string) (*Config, error) {
	confrc := dummyConfigrc()
	config_file := expandPath(path)
	if _, err := os.Stat(config_file); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		err := gcfg.ReadFileInto(confrc, config_file)
		if err != nil {
			return nil, err
		}
	}

	return parseConfig(confrc)
}

func dummyConfigrc() *configrc {
	var confrc configrc
	confrc.Global.Downloads = "."
	confrc.Global.Num = 20
	return &confrc
}

func parseConfig(confrc *configrc) (*Config, error) {
	var cfg Config

	downloads_expanded := expandPath(confrc.Global.Downloads)
	downloads_folder, err := filepath.Abs(downloads_expanded)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(downloads_folder); os.IsNotExist(err) {
		err := os.MkdirAll(downloads_folder, 0770)
		if err != nil {
			return nil, err
		}
	}
	cfg.Downloads = downloads_folder

	cfg.Lang = confrc.Global.Lang
	cfg.Num = confrc.Global.Num
	return &cfg, nil
}

func expandPath(path string) (rpath string) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path[:2] == "~/" {
		path = strings.Replace(path, "~/", dir+"/", 1)
	}
	return path
}
