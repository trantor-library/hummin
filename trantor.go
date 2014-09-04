package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"code.google.com/p/gcfg"
	"github.com/hailiang/gosocks"
)

type trantor struct {
	client   *http.Client
	download chan book
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

type configrc struct {
	Global struct {
		Downloads string
		Lang      string
	}
}

func Trantor() *trantor {
	var t trantor
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, PROXY)
	transport := &http.Transport{Dial: dialSocksProxy}
	t.client = &http.Client{Transport: transport}

	t.download = make(chan book, 20)
	for i := 0; i < DOWNLOAD_WORKERS; i++ {
		go t.downloadWorker()
	}
	return &t
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

func (t trantor) Download(id string, useWorker bool) error {
	b, err := t.Book(id)
	if err != nil {
		return err
	}
	if useWorker {
		t.download <- b
	} else {
		t.downloadBook(b)
	}
	return nil
}

func (t trantor) Search(query string, page int) (search, error) {
	var s search
	lang_in_query := strings.Count(query, "lang:")
	if lang_in_query < 1 {
		lang := getValueFromConfigrc("lang")
		if lang != "" {
			query = query + " lang:" + lang
		}
	}
	escaped_query := url.QueryEscape(query)
	err := t.get(BASE_URL+"search/"+"?q="+escaped_query+"&p="+strconv.Itoa(page)+"&fmt=json", &s)
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
	err = store(resp.Body, b.Title+".epub")
	if err != nil {
		printErr("There was a problem storing:", err)
		return
	}
	printDownloadFinished(b.Title)
}

func expandPath(path string) (rpath string) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path[:2] == "~/" {
		path = strings.Replace(path, "~/", dir+"/", 1)
	}
	return path
}

func getValueFromConfigrc(key string) (value string) {
	config_file := expandPath(CONFIG_FILE)
	if _, err := os.Stat(config_file); err != nil {
		if os.IsNotExist(err) {
			return ""
		} else {
			printErr("Error looking for config file:", err)
			return ""
		}
	}
	var cfg configrc
	err := gcfg.ReadFileInto(&cfg, config_file)
	if err != nil {
		printErr("Wrong config file:", err)
		return
	}

	if key == "downloads" {
		downloads_folder := cfg.Global.Downloads
		downloads_folder = expandPath(downloads_folder)
		downloads_folder, err = filepath.Abs(downloads_folder)
		if err != nil {
			printErr("Cannot get absolute path:", err)
			return
		}
		if _, err := os.Stat(downloads_folder); os.IsNotExist(err) {
			err := os.MkdirAll(downloads_folder, 0770)
			if err != nil {
				printErr("Could not create downloads folder:", err)
				return
			}
		}
		return downloads_folder
	} else if key == "lang" {
		lang := cfg.Global.Lang
		return lang
	}
	return ""
}

func store(src io.Reader, dest string) error {
	folder := getValueFromConfigrc("downloads")
	if folder != "" {
		dest = filepath.Join(folder, dest)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, src)
	return err
}
