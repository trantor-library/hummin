package main

import (
	"fmt"
	"strconv"

	"github.com/gobs/cmd"
)

type command struct {
	cmd  *cmd.Cmd
	t    *trantor
	ids  map[int]string
	last string
}

func Cmd(t *trantor) *command {
	commander := &cmd.Cmd{Prompt: "> "}
	commander.Init()
	c := &command{commander, t, nil, ""}

	commander.Add(cmd.Command{
		"book",
		"book id|num\nDisplays the book from the id or the index number in the last search",
		c.Book,
	})

	commander.Add(cmd.Command{
		"b",
		"book id|num\nDisplays the book from the id or the index number in the last search",
		c.Book,
	})

	commander.Add(cmd.Command{
		"get",
		"get [id|num]\nDownloads the last seen book, or an id or the index number from the last search",
		c.Get,
	})

	commander.Add(cmd.Command{
		"g",
		"get [id|num]\nDownloads the last seen book, or an id or the index number from the last search",
		c.Get,
	})

	commander.Add(cmd.Command{
		"search",
		`search books`,
		c.Search,
	})

	commander.Add(cmd.Command{
		"s",
		`search books`,
		c.Search,
	})

	commander.Add(cmd.Command{
		"exit",
		`exit program`,
		c.Exit,
	})

	commander.Add(cmd.Command{
		"quit",
		`same as exit`,
		c.Exit,
	})

	commander.Add(cmd.Command{
		"q",
		`same as exit`,
		c.Exit,
	})

	return c
}

func (c *command) Loop() {
	c.cmd.CmdLoop()
}

func (c *command) Book(line string) (stop bool) {
	id := line
	if len(line) != 16 {
		n, err := strconv.Atoi(line)
		if err != nil || len(c.ids) <= n {
			fmt.Println("Not valid search index:", line)
			return false
		}
		id = c.ids[n]
	}
	b, err := c.t.Book(id)
	if err != nil {
		fmt.Println("An error ocurred fetching the book info:", err)
		return false
	}

	c.last = id
	fmt.Println("Title:     ", b.Title)
	fmt.Print("Author:     ")
	for _, a := range b.Author {
		fmt.Print(a, ", ")
	}
	fmt.Println()
	fmt.Println("Publisher: ", b.Publisher)
	if b.Isbn != "" {
		fmt.Println("isbn:      ", b.Isbn)
	}
	fmt.Print("Tags:       ")
	for _, s := range b.Subject {
		fmt.Print(s, ", ")
	}
	fmt.Println()
	fmt.Print("Lang:       ")
	for _, l := range b.Lang {
		fmt.Print(l, ", ")
	}
	fmt.Println()
	fmt.Println()
	fmt.Println(b.Description)
	return false
}

func (c *command) Get(line string) (stop bool) {
	id := c.last
	if len(line) == 16 {
		id = line
	} else if len(line) > 0 {
		n, err := strconv.Atoi(line)
		if err != nil || len(c.ids) <= n {
			fmt.Println("Not valid search index:", line)
			return false
		}
		id = c.ids[n]
	}
	err := c.t.Download(id)
	if err != nil {
		fmt.Println("An error ocurred downloading the book:", err)
		return false
	}
	return false
}

func (c *command) Search(line string) (stop bool) {
	s, err := c.t.Search(line)
	if err != nil {
		fmt.Println("An error ocurred searching:", err)
		return false
	}

	c.ids = make(map[int]string, len(s.Books))
	fmt.Println("Found", s.Found)
	for i, b := range s.Books {
		c.ids[i] = b.Id

		fmt.Print("#", i, " ")
		if i < 10 {
			fmt.Print(" ")
		}
		fmt.Print("=> ")
		if len(b.Lang) > 0 {
			fmt.Print("[", b.Lang[0], "]")
		}
		fmt.Printf("%s (%s) || ", b.Title, b.Publisher)
		for _, a := range b.Author {
			fmt.Print(a, ", ")
		}
		fmt.Println()
	}
	if s.Found > s.Items {
		fmt.Println("(more)") //TODO
	}
	return false
}

func (c *command) Exit(line string) (stop bool) {
	fmt.Println("goodbye!")
	return true
}
