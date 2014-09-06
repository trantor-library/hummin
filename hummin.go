package main

import (
	"os"
	"strings"
)

func main() {
	cfg, err := ParseConfig(CONFIG_FILE)
	if err != nil {
		printErr("Error parsing config file: ", err)
		return
	}

	if len(os.Args) > 1 {
		comm(cfg)
	} else {
		shell(cfg)
	}
}

func comm(cfg *Config) {
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage(os.Args[0])
		return
	}

	t := Trantor(cfg)
	cmd := Cmd(t)
	cmd.OneCmd(strings.Join(os.Args[1:], " "))
}

func shell(cfg *Config) {
	t := Trantor(cfg)
	printLoading()
	idx, err := t.Index()
	if err != nil {
		printErr("Problem getting index: ", err)
		return
	}
	printWelcome(idx)

	cmd := Cmd(t)
	cmd.SetBooks(idx.Last_added)
	cmd.Loop()
}
