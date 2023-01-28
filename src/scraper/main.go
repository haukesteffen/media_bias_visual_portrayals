package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	c.searchPerson = "Armin Laschet"
	return c
}

func queryBuilder(svc *customsearch.Service, conf *SearchConfig) *customsearch.CseListCall {
	//resp, err := svc.Cse.List().Cx(conf.cx).Q(conf.query).SearchType(conf.searchType)
	lst := svc.Cse.List().Cx(conf.cx).SearchType(conf.searchType)
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

	for i, result := range resp.Items {
		fmt.Printf("#%d: %s\n", i+1, result.Title)
		fmt.Printf("\t%s\n", result.Link)
	}
}
