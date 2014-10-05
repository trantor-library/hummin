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
	notifications := make(chan Notification, 20)

	if len(os.Args) > 1 {
		comm(cfg, notifications)
	} else {
		shell(cfg, notifications)
	}
}

func comm(cfg *Config, notifications chan Notification) {
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage(os.Args[0])
		return
	}

	t := Trantor(cfg, notifications, false)
	cmd := Cmd(t, notifications)
	cmd.OneCmd(strings.Join(os.Args[1:], " "))
	cmd.PostCmd("", true)
}

func shell(cfg *Config, notifications chan Notification) {
	t := Trantor(cfg, notifications, true)
	printLoading()
	idx, err := t.Index()
	if err != nil {
		printErr("Problem getting index: ", err)
		return
	}
	printWelcome(idx)

	cmd := Cmd(t, notifications)
	cmd.SetBooks(idx.Last_added)
	cmd.Loop()
}
