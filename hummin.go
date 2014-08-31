package main

func main() {
	printLoading()
	t := Trantor()
	idx, err := t.Index()
	if err != nil {
		printErr("Problem getting index: ", err)
		return
	}
	printWelcome(idx)

	cmd := Cmd(t)
	cmd.Loop()
}
