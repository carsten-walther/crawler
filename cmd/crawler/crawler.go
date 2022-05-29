package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	sitemap "github.com/oxffaa/gopher-parse-sitemap"
)

func main() {

	urlParameter := flag.String("url", "", "URL to crawl, eg. https://www.your-website.com/sitemap.xml")
	csvParameter := flag.Bool("csv", false, "If set, write csv file.")
	maxParameter := flag.Int("max", 0, "If set, max URL request will be done.")
	silentParameter := flag.Bool("silent", false, "If set, there will be no output to console.")

	flag.Parse()

	if *urlParameter == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *silentParameter != true {
		fmt.Println("starting ...")
	}

	urls := make([]string, 0, 0)
	resultIndex := getIndex(*urlParameter)

	if len(resultIndex) > 0 {
		for _, index := range resultIndex {
			resultEndpoints := getEndpoint(index)
			if len(resultEndpoints) > 0 {
				for _, url := range resultEndpoints {
					urls = append(urls, url)
				}
			}
		}
	} else {
		resultEndpoints := getEndpoint(*urlParameter)
		if len(resultEndpoints) > 0 {
			for _, url := range resultEndpoints {
				urls = append(urls, url)
			}
		}
	}

	records := [][]string{
		{"timestamp", "state", "url"},
		{time.Now().Format(time.UnixDate), "200", *urlParameter},
	}

	count := 0

	for _, url := range urls {
		if *maxParameter <= 0 {
			records = append(records, callUrl(url, *silentParameter))
		} else {
			if count < *maxParameter {
				records = append(records, callUrl(url, *silentParameter))
			}
		}
		count++
	}

	if *csvParameter == true {
		writeCsv(records)
	}

	if *silentParameter != true {
		fmt.Println("done ...")
	}
}

func getIndex(url string) []string {

	result := make([]string, 0, 0)

	err := sitemap.ParseIndexFromSite(url, func(entry sitemap.IndexEntry) error {
		result = append(result, entry.GetLocation())
		return nil
	})

	if err != nil {
		panic(err)
	}

	return result
}

func getEndpoint(url string) []string {

	result := make([]string, 0, 0)

	err := sitemap.ParseFromSite(url, func(entry sitemap.Entry) error {
		result = append(result, entry.GetLocation())
		return nil
	})

	if err != nil {
		panic(err)
	}

	return result
}

func callUrl(url string, silent bool) []string {

	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	if silent == false {
		fmt.Println(fmt.Sprintf("%s, %v, %s", time.Now().Format(time.UnixDate), response.StatusCode, url))
	}

	return []string{time.Now().Format(time.UnixDate), fmt.Sprintf("%v", response.StatusCode), url}
}

func writeCsv(records [][]string) {

	f, err := os.Create("results.csv")
	defer f.Close()

	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		panic(err)
	}
}
