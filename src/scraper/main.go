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
	apiKey     string
	cx         string
	query      string
	searchType string
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
	c.searchType = "image"
	return c
}

func main() {
	conf := newConfig()
	ctx := context.Background()
	svc, err := customsearch.NewService(ctx, option.WithAPIKey(conf.apiKey))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := svc.Cse.List().Cx(conf.cx).Q(conf.query).SearchType(conf.searchType).Do()
	if err != nil {
		log.Fatal(err)
	}

	for i, result := range resp.Items {
		fmt.Printf("#%d: %s\n", i+1, result.Title)
		fmt.Printf("\t%s\n", result.Link)
	}
}
