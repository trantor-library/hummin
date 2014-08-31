package main

import (
	"strconv"

	"github.com/gobs/cmd"
)

type command struct {
	cmd   *cmd.Cmd
	t     *trantor
	ids   []string
	last  string
	query string
	page  int
	more  bool
}

func Cmd(t *trantor) *command {
	commander := &cmd.Cmd{Prompt: "> "}
	commander.Init()
	c := &command{commander, t, nil, "", "", 0, false}

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
		"more",
		`more books from previows search`,
		c.More,
	})

	commander.Add(cmd.Command{
		"m",
		`more books from previows search`,
		c.More,
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
			printErr("Not valid search index: "+line, nil)
			return false
		}
		id = c.ids[n]
	}

	b, err := c.t.Book(id)
	if err != nil {
		printErr("An error ocurred fetching the book info:", err)
		return false
	}
	c.last = id

	printBook(b)
	return false
}

func (c *command) Get(line string) (stop bool) {
	id := c.last
	if len(line) == 16 {
		id = line
	} else if len(line) > 0 {
		n, err := strconv.Atoi(line)
		if err != nil || len(c.ids) <= n {
			printErr("Not valid search index: "+line, nil)
			return false
		}
		id = c.ids[n]
	}
	err := c.t.Download(id)
	if err != nil {
		printErr("An error ocurred downloading the book:", err)
		return false
	}
	return false
}

func (c *command) Search(line string) (stop bool) {
	c.query = line
	c.page = 0
	c.ids = []string{}
	c.doSearch()
	return false
}

func (c *command) More(line string) (stop bool) {
	if c.more {
		c.page++
		c.doSearch()
	}
	return false
}

func (c *command) doSearch() {
	s, err := c.t.Search(c.query, c.page)
	if err != nil {
		printErr("An error ocurred searching:", err)
		return
	}

	idx := len(c.ids)
	c.page = s.Page
	for _, b := range s.Books {
		c.ids = append(c.ids, b.Id)
	}
	c.more = s.Found > (s.Page+1)*s.Items
	printSearch(s, idx, c.more)
}

func (c *command) Exit(line string) (stop bool) {
	printExit()
	return true
}
