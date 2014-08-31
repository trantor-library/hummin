package main

import (
	"fmt"

	"github.com/ttacon/chalk"
)

func printUsage(prog string) {
	fmt.Println(prog, "[command [args]]")
	fmt.Println("Commands: book, get, search, news")
}

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

	fmt.Println("Last added books:")
	printListBooks(idx.Last_added, 0, false)
	fmt.Println()

	if len(idx.News) > 0 {
		fmt.Println(chalk.Yellow, "<<", idx.News[0].Date, "-", idx.News[0].Text, ">>", chalk.Reset)
		fmt.Println()
	}
}

func printBook(b book) {
	fieldStyle := chalk.Bold.NewStyle().WithForeground(chalk.Cyan)
	fmt.Println(fieldStyle, "Id:       ", chalk.Reset, b.Id)
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

func printSearch(s search, startIdx int, more bool, fullId bool) {
	fmt.Println("   Found", s.Found)
	printListBooks(s.Books, startIdx, fullId)
	if more {
		fmt.Println("(more)")
	}
}

func printListBooks(books []book, startIdx int, fullId bool) {
	i := startIdx
	for _, b := range books {
		if fullId {
			fmt.Print(chalk.Magenta, b.Id, " ", chalk.Reset)
		} else {
			fmt.Print(chalk.Magenta, "#", i, " ", chalk.Reset)
			if i < 10 {
				fmt.Print(" ")
			}
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
		i++
	}
}

func printNews(ns []news) {
	for _, n := range ns {
		fmt.Println(chalk.Yellow, n.Date, chalk.Reset, n.Text)
	}
}

func printExit() {
	fmt.Println(chalk.Green, "goodbye!", chalk.Reset)
}

func printDownloadFinished(title string) {
	fmt.Println(chalk.Yellow, "Download of", title, "finished", chalk.Reset)
}
