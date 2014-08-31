package main

import (
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		comm()
	} else {
		shell()
	}
}

func comm() {
	t := Trantor()
	cmd := Cmd(t)
	cmd.OneCmd(strings.Join(os.Args[1:], " "))
}

func shell() {
	t := Trantor()
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
