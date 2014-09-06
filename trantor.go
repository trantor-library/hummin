package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hailiang/gosocks"
)

type trantor struct {
	cfg        *Config
	useWorkers bool
	client     *http.Client
	download   chan book
}

type book struct {
	Id          string
	Title       string
	Author      []string
	Contributor string
	Publisher   string
	Description string
	Subject     []string
	Date        string
	Lang        []string
	Isbn        string
	Size        int
	Cover       string
	Cover_small string
	Download    string
	Read        string
}

type search struct {
	Found int
	Page  int
	Items int
	Books []book
}

type news struct {
	Date string
	Text string
}

type index struct {
	Title      string
	Url        string
	Count      int
	News       []news
	Tags       []string
	Last_added []book
}

func Trantor(cfg *Config, useWorkers bool) *trantor {
	var t trantor
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, PROXY)
	transport := &http.Transport{Dial: dialSocksProxy}
	t.client = &http.Client{Transport: transport}
	t.cfg = cfg
	t.useWorkers = useWorkers

	if useWorkers {
		t.spanWorkers()
	}
	return &t
}

func (t *trantor) spanWorkers() {
	t.download = make(chan book, 20)
	for i := 0; i < DOWNLOAD_WORKERS; i++ {
		go t.downloadWorker()
	}
}

func (t trantor) Index() (index, error) {
	var i index
	err := t.get(BASE_URL+"?fmt=json", &i)
	return i, err
}

func (t trantor) Book(id string) (book, error) {
	var b book
	err := t.get(BASE_URL+"book/"+id+"?fmt=json", &b)
	return b, err
}

func (t trantor) Download(id string) error {
	b, err := t.Book(id)
	if err != nil {
		return err
	}
	if t.useWorkers {
		t.download <- b
	} else {
		t.downloadBook(b)
	}
	return nil
}

func (t trantor) Search(query string, page int) (search, error) {
	var s search
	if hasSubString(query, "lang:any") {
		query = strings.Replace(query, "lang:any", "", -1)
	} else if !hasSubString(query, "lang:") {
		if t.cfg.Lang != "" {
			query = query + " lang:" + t.cfg.Lang
		}
	}
	pageStr := strconv.Itoa(page)
	numStr := strconv.Itoa(t.cfg.Num)
	escaped_query := url.QueryEscape(query)
	err := t.get(BASE_URL+"search/"+"?q="+escaped_query+"&p="+pageStr+"&fmt=json"+"&num="+numStr, &s)
	return s, err
}

func (t trantor) News() ([]news, error) {
	var n []news
	err := t.get(BASE_URL+"news/?fmt=json", &n)
	return n, err
}

func (t trantor) get(url string, v interface{}) error {
	resp, err := t.client.Get(url)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(v)
}

func (t trantor) downloadWorker() {
	for b := range t.download {
		t.downloadBook(b)
	}
}

func (t trantor) downloadBook(b book) {
	resp, err := t.client.Get(BASE_URL + b.Download)
	if err != nil {
		printErr("There was a problem with the download:", err)
		return
	}
	err = store(resp.Body, filepath.Join(t.cfg.Downloads, b.Title+".epub"))
	if err != nil {
		printErr("There was a problem storing:", err)
		return
	}
	printDownloadFinished(b.Title)
}

func store(src io.Reader, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, src)
	return err
}

func hasSubString(str string, substr string) bool {
	matches := strings.Count(str, substr)
	return matches >= 1
}
