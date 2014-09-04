package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"code.google.com/p/gcfg"
)

func expandPath(path string) (rpath string) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path[:2] == "~/" {
		path = strings.Replace(path, "~/", dir+"/", 1)
	}
	return path
}

func getValueFromConfigrc(key string) (value string) {
	config_file := expandPath(CONFIG_FILE)
	if _, err := os.Stat(config_file); err != nil {
		if os.IsNotExist(err) {
			return ""
		} else {
			printErr("Error looking for config file:", err)
			return ""
		}
	}
	var cfg configrc
	err := gcfg.ReadFileInto(&cfg, config_file)
	if err != nil {
		printErr("Wrong config file:", err)
		return
	}

	if key == "downloads" {
		downloads_folder := cfg.Global.Downloads
		downloads_folder = expandPath(downloads_folder)
		downloads_folder, err = filepath.Abs(downloads_folder)
		if err != nil {
			printErr("Cannot get absolute path:", err)
			return
		}
		if _, err := os.Stat(downloads_folder); os.IsNotExist(err) {
			err := os.MkdirAll(downloads_folder, 0770)
			if err != nil {
				printErr("Could not create downloads folder:", err)
				return
			}
		}
		return downloads_folder
	} else if key == "lang" {
		lang := cfg.Global.Lang
		return lang
	}
	return ""
}
