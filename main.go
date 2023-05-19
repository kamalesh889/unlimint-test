package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	webPages = []string{
		"https://www.yahoo.com/",
		"https://www.google.com/",
		"https://www.bing.com/?toWww=1&redig=A3A6CB4417844F3EA56B64A926F4BBFA",
		"https://www.amazon.com/",
		"https://github.com/",
		"https://about.gitlab.com/",
		"https://abcdinvalid.com/",
	}

	results struct {
		// put here content length of each page
		ContentLength map[string]int

		// accumulate here the content length of all pages
		TotalContentLength int
	}
)

/*

	One thing i am assuming that if i am not able to get the content length of any webpage ,
	then i am not adding -1 to the  TotalContentLength .

	I am only adding values to TotalContentLength for which webpage i am able to get its
	content length

*/

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch := make(chan string)
	chk := make(chan bool)

	go workerFetchWebPage(ctx, ch)
	go workerGetContentLength(ctx, ch, chk)

	<-chk

	for key, val := range results.ContentLength {
		fmt.Println("Content length of webpage:", key, "is", val)
	}

	fmt.Println("Total content length of all webpages are", results.TotalContentLength)

}

func workerFetchWebPage(ctx context.Context, page chan string) {

	for _, val := range webPages {

		select {
		case <-ctx.Done():

			fmt.Println("Time limit exceeded : Exiting from workerFetchWebPage")
			return

		default:

			page <- val
		}
	}

	close(page)

}

func workerGetContentLength(ctx context.Context, page chan string, chker chan bool) {

	pagemap := make(map[string]int)
	totallength := 0

	for webpage := range page {

		select {
		case <-ctx.Done():

			fmt.Println("Time limit exceeded : Exiting from workerGetContentLength")
			chker <- false
			return

		default:

			contentlength, err := getContentlenght(webpage)

			if err == nil {
				totallength = totallength + contentlength
			}

			pagemap[webpage] = contentlength
		}

	}

	results.ContentLength = pagemap
	results.TotalContentLength = totallength

	chker <- true

}

// This function will scrape the website and get the content length
func getContentlenght(webpage string) (int, error) {

	response, err := http.Get(webpage)
	if err != nil {
		return -1, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	return len(body), nil

}
