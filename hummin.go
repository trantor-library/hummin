package main

import "fmt"

func main() {
	t := Trantor()
	welcome(t)

	cmd := Cmd(t)
	cmd.Loop()
}

func welcome(t *trantor) {
	fmt.Println("Welcome to hummin")
	fmt.Println("    ...connecting")
	fmt.Println()

	idx, err := t.Index()
	if err != nil {
		fmt.Println("Problem getting index: ", err)
		return
	}
	fmt.Println(idx.Title)
	fmt.Println(idx.Url)
	fmt.Println("    Num books:", idx.Count)
	fmt.Println()

	if len(idx.News) > 0 {
		fmt.Println("<<", idx.News[0].Text, ">>")
	}
	fmt.Println()
}
