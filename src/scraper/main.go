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
	dbhostname   string
	dbport       string
	dbuser       string
	dbpassword   string
	dbdbname     string
}

type dbElement struct {
	politican int    `db:"politican"`
	site      string `db:"site"`
	data      string `db:"data"`
}

func newConfig() *SearchConfig {
	c := &SearchConfig{}
	// Secret data from env vars
	// ... search data from args
	flag.StringVar(&c.searchPerson, "person", "", "Person to search for")
	flag.StringVar(&c.searchSite, "site", "", "Search for picture on this site")
	flag.StringVar(&c.folder, "dir", "/data/", "Directory to save pictures to")
	flag.BoolVar(&c.to_db, "db", false, "Write to database")
	needs_help := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *needs_help {
		flag.Usage()
		os.Exit(0)
	}
	val, err := os.LookupEnv("SEARCHAPIKEY")
	if !err {
		log.Fatal("Set SEARCHAPIKEY env")
	}
	c.apiKey = val
	val, err = os.LookupEnv("SEARCHCX")
	if !err {
		log.Fatal("Set SEARCHCX env")
	}
	c.cx = val
	if len(c.searchPerson) == 0 || len(c.searchSite) == 0 {
		log.Fatal("Please provide --person and --site arguments")
	}
	c.searchType = "image"
	if c.to_db {
		val, err = os.LookupEnv("DBHOSTNAME")
		if !err {
			log.Fatal("Set DBHOSTNAME env")
		}
		c.dbhostname = val
		val, err = os.LookupEnv("DBPORT")
		if !err {
			log.Fatal("Set DBPORT env")
		}
		c.dbport = val
		val, err = os.LookupEnv("DBUSER")
		if !err {
			log.Fatal("Set DBUSER env")
		}
		c.dbuser = val
		val, err = os.LookupEnv("DBPASSWORD")
		if !err {
			log.Fatal("Set DBPASSWORD env")
		}
		c.dbpassword = val
		val, err = os.LookupEnv("DBDBNAME")
		if !err {
			log.Fatal("Set DBDBNAME env")
		}
		c.dbdbname = val
	}
	return c
}

func queryBuilder(svc *customsearch.Service, conf *SearchConfig, startIndex int64) *customsearch.CseListCall {
	// todo: defaults auslagern
	//lst := svc.Cse.List().Cx(conf.cx).SearchType(conf.searchType).Sort("date:d:s").Num(10).Start(1 + startIndex*10).ImgType("face")
	lst := svc.Cse.List().Cx(conf.cx)
	// default: search for images
	lst = lst.SearchType(conf.searchType)
	// Bias results strongly towards newer dates
	lst = lst.Sort("date:d:s")
	// 10 results per page
	resultsPerPage := 10
	lst = lst.Num(int64(resultsPerPage))
	// Start at this offset
	lst = lst.Start(1 + startIndex*int64(resultsPerPage))
	// Search for faces
	lst = lst.ImgType("face")
	query := ""
	query += "site:" + conf.searchSite
	query += " " + conf.searchPerson
	// Query layout: `site:${site} ${person}`
	lst = lst.Q(query)
	log.Println("debug search: ", *lst)
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
		log.Println("Writing to Database")
		connectDB(conf)
		err := toDB(conf, done, len(urls))
		return err
	} else {
		log.Println("Writing to filesystem to folder", conf.folder)
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
}

var dbconn *pgx.Conn

func connectDB(conf *SearchConfig) {
	var connErr error
	databaseString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.dbuser, conf.dbpassword, conf.dbhostname, conf.dbport, conf.dbdbname)
	//dbconn, connErr = pgx.Connect(context.Background(), os.Getenv("DATABASE_STRING"))
	dbconn, connErr = pgx.Connect(context.Background(), databaseString)
	if connErr != nil {
		log.Fatalf("Unable to connect to database: %v\n", connErr)
	}
	//defer dbconn.Close(context.Background())
}

func toDB(conf *SearchConfig, receive chan []byte, qlen int) error {
	// todo das hier zu picElement vielleicht?
	tmp := [][]interface{}{}
	// todo vielleicht anders mit q?
	for i := 0; i < qlen; i++ {
		tmp = append(tmp, []interface{}{conf.searchPerson, conf.searchSite, <-receive})
	}
	copyCount, queryErr := dbconn.CopyFrom(
		context.Background(),
		pgx.Identifier{"items"},
		[]string{"politican", "site", "data"},
		pgx.CopyFromRows(tmp),
	)
	if queryErr != nil {
		log.Printf("Query Error: %v\n", queryErr)
	}
	log.Println("CopyCount:", copyCount)
	return queryErr
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
