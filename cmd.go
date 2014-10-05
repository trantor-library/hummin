package main

import (
	"strconv"

	"github.com/gobs/cmd"
)

type command struct {
	cmd           *cmd.Cmd
	t             *trantor
	books         []book
	last          string
	query         string
	page          int
	more          bool
	shell         bool
	notifications chan Notification
}

type Notification struct {
	str string
	err error
}

func Cmd(t *trantor, notifications chan Notification) *command {
	commander := &cmd.Cmd{Prompt: "> "}
	commander.Init()
	c := &command{
		cmd:           commander,
		t:             t,
		books:         []book{},
		last:          "",
		query:         "",
		page:          0,
		more:          false,
		shell:         false,
		notifications: notifications,
	}

	commander.PostCmd = c.PostCmd

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
		"news",
		`get the list of news`,
		c.News,
	})

	commander.Add(cmd.Command{
		"n",
		`get the list of news`,
		c.News,
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

func (c *command) OneCmd(line string) {
	c.cmd.OneCmd(line)
}

func (c *command) SetBooks(books []book) {
	c.books = books
}

func (c *command) Loop() {
	c.shell = true
	c.cmd.CmdLoop()
}

func (c *command) PostCmd(line string, stop bool) bool {
	done := false
	for !done {
		select {
		case n := <-c.notifications:
			if n.err != nil {
				printErr(n.str, n.err)
			} else {
				printNotification(n.str)
			}
		default:
			done = true
		}
	}

	return stop
}

func (c *command) Book(line string) (stop bool) {
	var b book
	id, n := c.getId(line, "")
	if n != -1 {
		b = c.books[n]
	} else {
		if id == "" {
			printErr("Not valid id "+line, nil)
			return false
		}

		var err error
		b, err = c.t.Book(id)
		if err != nil {
			printErr("An error ocurred fetching the book info:", err)
			return false
		}
	}

	c.last = b.Id
	printBook(b, !c.shell)
	return false
}

func (c *command) Get(line string) (stop bool) {
	id, _ := c.getId(line, c.last)
	err := c.t.Download(id)
	if err != nil {
		printErr("An error ocurred downloading the book:", err)
		return false
	}
	return false
}

func (c *command) getId(line string, fallBack string) (id string, n int) {
	n = -1
	id = fallBack
	if len(line) == 16 {
		id = line
	} else if len(line) > 0 {
		var err error
		n, err = strconv.Atoi(line)
		if err != nil || len(c.books) <= n {
			return
		}
		id = c.books[n].Id
	}
	return
}

func (c *command) Search(line string) (stop bool) {
	c.query = line
	c.page = 0
	c.books = []book{}
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

func (c *command) News(line string) (stop bool) {
	ns, err := c.t.News()
	if err != nil {
		printErr("An error ocurred getting the list of news:", err)
		return false
	}
	printNews(ns)
	return false
}

func (c *command) doSearch() {
	s, err := c.t.Search(c.query, c.page)
	if err != nil {
		printErr("An error ocurred searching:", err)
		return
	}

	idx := len(c.books)
	c.page = s.Page
	c.SetBooks(s.Books)
	c.more = s.Found > (s.Page+1)*s.Items
	printSearch(s, idx, c.more, !c.shell)
}

func (c *command) Exit(line string) (stop bool) {
	printExit()
	return true
}
