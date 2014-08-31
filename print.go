package main

import (
	"fmt"

	"github.com/ttacon/chalk"
)

func printErr(msg string, err error) {
	fmt.Println(chalk.Red, msg, err, chalk.Reset)
}

func printLoading() {
	fmt.Println(chalk.Green, "Welcome to hummin", chalk.Reset)
	fmt.Println("     ...connecting")
	fmt.Println()
}

func printWelcome(idx index) {
	fmt.Println(chalk.Blue, idx.Title)
	fmt.Println(idx.Url, chalk.Reset)
	fmt.Println("      Num books:", idx.Count)
	fmt.Println()

	if len(idx.News) > 0 {
		fmt.Println(chalk.Yellow, "<<", idx.News[0].Text, ">>", chalk.Reset)
	}
	fmt.Println()
}

func printBook(b book) {
	fieldStyle := chalk.Bold.NewStyle().WithForeground(chalk.Cyan)
	fmt.Println(fieldStyle, "Title:    ", chalk.Reset, b.Title)
	fmt.Print(fieldStyle, " Author:     ", chalk.Reset)
	for _, a := range b.Author {
		fmt.Print(a, ", ")
	}
	fmt.Println()
	if b.Publisher != "" {
		fmt.Println(fieldStyle, "Publisher:", chalk.Reset, b.Publisher)
	}
	if b.Isbn != "" {
		fmt.Println(fieldStyle, "isbn:     ", chalk.Reset, b.Isbn)
	}
	fmt.Print(fieldStyle, " Tags:       ", chalk.Reset)
	for _, s := range b.Subject {
		fmt.Print(s, ", ")
	}
	fmt.Println()
	fmt.Print(fieldStyle, " Lang:       ", chalk.Reset)
	for _, l := range b.Lang {
		fmt.Print(l, ", ")
	}
	fmt.Println()
	fmt.Println()
	fmt.Println(b.Description)
}

func printSearch(s search) {
	fmt.Println("   Found", s.Found)
	for i, b := range s.Books {
		fmt.Print(chalk.Magenta, "#", i, " ", chalk.Reset)
		if i < 10 {
			fmt.Print(" ")
		}
		fmt.Print("=> ")
		if len(b.Lang) > 0 {
			fmt.Print(chalk.Cyan, "[", b.Lang[0][:2], "]", chalk.Reset)
		}
		fmt.Print(b.Title, " ")
		if b.Publisher != "" {
			fmt.Printf("%s(%s)%s ", chalk.Yellow, b.Publisher, chalk.Reset)
		}
		fmt.Print(chalk.Green)
		for _, a := range b.Author {
			fmt.Print(a, ", ")
		}
		fmt.Println(chalk.Reset)
	}
	if s.Found > s.Items {
		fmt.Println("(more)")
	}
}

func printExit() {
	fmt.Println(chalk.Green, "goodbye!", chalk.Reset)
}

func printDownloadFinished(title string) {
	fmt.Println(chalk.Yellow, "Download of", title, "finished", chalk.Reset)
}
