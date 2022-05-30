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

	var file *os.File
	var urls []string
	var record []string
	var currentCount int

	var url string
	flag.StringVar(&url, "url", "", "URL to crawl, eg. https://www.your-website.com/sitemap.xml")

	var output string
	flag.StringVar(&output, "output", "", "Path to results csv file.")

	var count int
	flag.IntVar(&count, "count", 0, "Maximum count of urls to crawl.")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Show verbose output.")

	flag.Parse()

	if url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if output != "" {
		file = createFile(output)
		defer closeFile(file)
		writeFile([]string{"id", "datetime", "state", "url"}, file)
	}

	verboseOutput("Starting ...", verbose)

	urls = collectUrls(url)
	currentCount = 1

	for _, url := range urls {
		if count <= 0 {
			record = callUrl(currentCount, url, verbose)
			if output != "" {
				writeFile(record, file)
			}
		} else {
			if currentCount <= count {
				record = callUrl(currentCount, url, verbose)
				if output != "" {
					writeFile(record, file)
				}
			}
		}
		currentCount++
	}

	verboseOutput("Done ...", verbose)
}

/* network functions */

func collectUrls(source string) []string {
	var urls []string
	resultIndex := getIndex(source)
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
		resultEndpoints := getEndpoint(source)
		if len(resultEndpoints) > 0 {
			for _, url := range resultEndpoints {
				urls = append(urls, url)
			}
		}
	}
	return urls
}

func getIndex(url string) []string {
	var result []string
	err := sitemap.ParseIndexFromSite(url, func(entry sitemap.IndexEntry) error {
		result = append(result, entry.GetLocation())
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return result
}

func getEndpoint(url string) []string {
	var result []string
	err := sitemap.ParseFromSite(url, func(entry sitemap.Entry) error {
		result = append(result, entry.GetLocation())
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return result
}

func callUrl(count int, url string, verbose bool) []string {
	response, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	verboseOutput(fmt.Sprintf("%v, %v : %s", count, response.StatusCode, url), verbose)
	return []string{fmt.Sprintf("%v", count), time.Now().Format(time.UnixDate), fmt.Sprintf("%v", response.StatusCode), url}
}

/* verbose functions */

func verboseOutput(message string, verbose bool) {
	if verbose == true {
		fmt.Println(fmt.Sprintf("# %s : %s", time.Now().Format(time.UnixDate), message))
	}
}

/* file functions */

func createFile(p string) *os.File {
	f, err := os.Create(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return f
}

func writeFile(record []string, f *os.File) {
	w := csv.NewWriter(f)
	err := w.Write(record)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	w.Flush()
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
