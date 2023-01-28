package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type SearchConfig struct {
	apiKey       string
	cx           string
	query        string
	searchType   string
	searchSites  []string
	searchPerson string
}

func newConfig() *SearchConfig {
	c := &SearchConfig{}
	val, err := os.LookupEnv("SEARCHAPIKEY")
	if !err {
		fmt.Fprintln(os.Stderr, "Set SEARCHAPIKEY env")
		os.Exit(1)
	}
	c.apiKey = val
	val, err = os.LookupEnv("SEARCHCX")
	if !err {
		fmt.Fprintln(os.Stderr, "Set SEARCHCX env")
		os.Exit(1)
	}
	c.cx = val
	val, err = os.LookupEnv("SEARCHQUERY")
	if !err {
		fmt.Fprintln(os.Stderr, "Set SEARCHQUERY env")
		os.Exit(1)
	}
	c.query = val
	// todo: vars
	c.searchType = "image"
	c.searchSites = []string{"tagesschau.de", "faz.net", "welt.de", "bild.de"}
	c.searchPerson = "erich honecker"
	return c
}

func queryBuilder(svc *customsearch.Service, conf *SearchConfig) *customsearch.CseListCall {
	// todo: defaults auslagern
	lst := svc.Cse.List().Cx(conf.cx).SearchType(conf.searchType).Sort("date")
	query := ""
	// todo: maybe use 'orTerms' here
	for i, site := range conf.searchSites {
		query += "site:" + site
		// Dont add 'OR' keyword at the end
		if i != len(conf.searchSites)-1 {
			query += " OR "
		}
	}
	query += " " + conf.searchPerson
	lst.Q(query)
	log.Println("debug searc: ", *lst)
	return lst
}

func downloadFile(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, response.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func fetchAndSaveImage(urls []string) error {
	done := make(chan []byte, len(urls))
	errch := make(chan error, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			b, err := downloadFile(URL)
			if err != nil {
				errch <- err
				done <- nil
				return
			}
			done <- b
			errch <- nil
		}(URL)
	}

	rand.Seed(time.Now().UnixNano())
	var errStr string
	for i := 0; i < len(urls); i++ {
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
		file, err := os.Create("./data/asd" + fmt.Sprint(rand.Intn(1000000)))
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, bytes.NewReader(<-done))
		if err != nil {
			return err
		}
	}
	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}
	return err
}

func main() {
	conf := newConfig()
	ctx := context.Background()
	svc, err := customsearch.NewService(ctx, option.WithAPIKey(conf.apiKey))
	if err != nil {
		log.Fatal(err)
	}
	svcQuery := queryBuilder(svc, conf)

	resp, err := svcQuery.Do()
	if err != nil {
		log.Fatal(err)
	}

	urlList := []string{}
	for i, result := range resp.Items {
		fmt.Printf("#%d: %s\n", i+1, result.Title)
		fmt.Printf("\t%s\n", result.Link)
		urlList = append(urlList, result.Link)
	}
	fetchAndSaveImage(urlList)
}
