package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type SearchConfig struct {
	apiKey       string
	cx           string
	folder       string
	searchType   string
	searchSite   string
	searchPerson string
	to_db        bool
}

type dbElement struct {
	politican int    `db:"politican"`
	site      string `db:"site"`
	data      string `db:"data"`
}

func newConfig() *SearchConfig {
	c := &SearchConfig{}
	// Secret data from env vars
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
	// ... search data from args
	flag.StringVar(&c.searchPerson, "person", "", "Person to search for")
	flag.StringVar(&c.searchSite, "site", "", "Site to use")
	flag.StringVar(&c.folder, "dir", "/data/", "Directory to save pictures to")
	flag.BoolVar(&c.to_db, "db", true, "Write to database")
	needs_help := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *needs_help {
		flag.Usage()
		os.Exit(0)
	}
	if len(c.searchPerson) == 0 || len(c.searchSite) == 0 {
		log.Fatal("Please provide -person and -site arguments")
	}
	c.searchType = "image"
	return c
}

func queryBuilder(svc *customsearch.Service, conf *SearchConfig, startIndex int64) *customsearch.CseListCall {
	// todo: defaults auslagern
	lst := svc.Cse.List().Cx(conf.cx).SearchType(conf.searchType).Sort("date:d:s").Num(10).Start(1 + startIndex*10).ImgType("face")
	query := ""
	query += "site:" + conf.searchSite
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

func fetchAndSaveImage(urls []string, conf *SearchConfig) error {
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
	// TODO das splitt und in funktionen packen
	if conf.to_db {
		toDB(conf, done, len(urls))
	} else {
		var errStr string
		picNamesBase := path.Join(conf.folder, strings.ReplaceAll(conf.searchPerson, " ", "-")+"_"+conf.searchSite+"_")
		for i := 0; i < len(urls); i++ {
			if err := <-errch; err != nil {
				errStr = errStr + " " + err.Error()
			}
			now := time.Now().UnixNano()
			file, err := os.Create(picNamesBase + fmt.Sprint(now) + ".jpg")
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
	// todo return besser
	return nil
}

func toDB(conf *SearchConfig, receive chan []byte, qlen int) error {
	// todo das hier zu picElement vielleicht?
	tmp := [][]interface{}{}

	// todo vielleicht anders mit q?
	for i := 0; i < qlen; i++ {
		tmp = append(tmp, []interface{}{conf.searchPerson, conf.searchSite, <-receive})
	}

	fmt.Println("Before DB connect")
	//conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	copyCount, queryErr := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"items"},
		[]string{"politican", "site", "data"},
		pgx.CopyFromRows(tmp),
	)
	fmt.Println(copyCount)
	fmt.Println(queryErr)
	return err
}

type picElement struct {
	politican string
	site      string
	data      []byte
}

func main() {
	conf := newConfig()
	ctx := context.Background()
	svc, err := customsearch.NewService(ctx, option.WithAPIKey(conf.apiKey))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 1; i++ {
		svcQuery := queryBuilder(svc, conf, int64(i))

		resp, err := svcQuery.Do()
		if err != nil {
			log.Fatal(err)
		}

		urlList := []string{}
		for _, result := range resp.Items {
			//fmt.Printf("#%d: %s\n", i+1, result.Title)
			//fmt.Printf("\t%s\n", result.Link)
			urlList = append(urlList, result.Link)
		}
		fetchAndSaveImage(urlList, conf)
	}
}
