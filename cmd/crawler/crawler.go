package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	color "github.com/TwiN/go-color"
	sitemap "github.com/oxffaa/gopher-parse-sitemap"
	terminal "golang.org/x/term"
)

var url string
var output string
var count int
var verbose bool

func main() {

	var file *os.File

	width, _, err := terminal.GetSize(0)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	flag.StringVar(&url, "url", "", "URL to crawl, eg. https://www.your-website.com/sitemap.xml")
	flag.StringVar(&output, "output", "", "Path to results csv file.")
	flag.IntVar(&count, "count", 0, "Maximum count of urls to crawl.")
	flag.BoolVar(&verbose, "verbose", false, "Show verbose output.")

	flag.Parse()

	fmt.Printf("\nSimple CLI crawler\n")
	fmt.Printf("\n%s\n", PrintLine("=", width*2/3))

	if url == "" {
		fmt.Printf("\nUsage:\n")
		flag.PrintDefaults()
		fmt.Printf("\n")
		os.Exit(0)
	}

	fmt.Printf("\n")

	if output != "" {
		file = createFile(output)
		defer closeFile(file)
		writeFile([]string{"id", "datetime", "state", "url"}, file)
	}

	verboseOutput("> Fetching URLs ...")

	urls := collectUrls(url)

	verboseOutput(fmt.Sprintf("> %v URLs found ...\n", len(urls)))

	currentCount := 0

	for _, url := range urls {
		currentCount++
		if count <= 0 {
			record := callUrl(currentCount, url, len(urls))
			if output != "" {
				writeFile(record, file)
			}
		} else {
			if currentCount <= count {
				record := callUrl(currentCount, url, len(urls))
				if output != "" {
					writeFile(record, file)
				}
			}
		}
	}

	verboseOutput("\n> Done ...")
}

/* network functions */

func collectUrls(source string) []string {
	var urls []string
	resultIndex := getIndex(source)
	if len(resultIndex) > 0 {
		for _, index := range resultIndex {
			urls = append(urls, getEndpoint(index)...)
		}
	} else {
		urls = append(urls, getEndpoint(source)...)
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

func callUrl(count int, url string, index int) []string {
	response, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	message := fmt.Sprintf("%s, %s, %v, %s", lpad(fmt.Sprint(count), " ", recursionCountDigits(index)), time.Now().Format(time.UnixDate), response.StatusCode, url)

	switch response.StatusCode {
	case 200:
		verboseOutput(color.Colorize(color.Green, message))
	case 404:
		verboseOutput(color.Colorize(color.Red, message))
	case 500:
		verboseOutput(color.Colorize(color.Red, message))
	default:
		verboseOutput(message)
	}

	return []string{fmt.Sprintf("%v", count), time.Now().Format(time.UnixDate), fmt.Sprintf("%v", response.StatusCode), url}
}

/* verbose functions */

func verboseOutput(message string) {
	if verbose {
		fmt.Println(message)
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

/* tool functions */

func recursionCountDigits(number int) int {
	if number < 10 {
		return 1
	} else {
		return 1 + recursionCountDigits(number/10)
	}
}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func PrintLine(char string, chars int) string {
	s := ""
	for i := 0; i < chars; i++ {
		s = char + s
	}
	return s
}
