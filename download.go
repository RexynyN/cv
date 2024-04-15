package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func goDownload(url string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Get the image
	response, err := http.Get(url)
	fmt.Println(url, response.StatusCode)
	if err != nil || response.StatusCode != 200 {
		ch <- url
		return
	}
	defer response.Body.Close()

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded", fileName, response.Body)
	ch <- ""
}

func main() {
	urls := []string{
		"https://pbs.twimg.com/media/GLAhmB_WcAAlVvL.jpg",
		"https://pbs.twimg.com/media/GK_4xkOXsAA2wQW.jpg",
		"https://pbs.twimg.com/media/GK_9MM_Wo3DxG.jpg",
		"https://pbs.twimg.com/media/GK_9MM_WoAA3DxG.jpg",
	}

	DownloadBatch(urls)
}

func GetLeftovers(ch chan string, left chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	leftOvers := make([]string, 0)
	for val := range ch {
		if val != "" {
			leftOvers = append(leftOvers, val)
		}
	}
	fmt.Println(leftOvers)
	left <- leftOvers
}

func DownloadBatch(urls []string) {
	urlsChan := make(chan string)
	leftoversChan := make(chan []string)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go goDownload(url, urlsChan, &wg)
	}
	wg.Wait()

	wg.Add(1)
	go GetLeftovers(urlsChan, leftoversChan, &wg)
	urls = <-leftoversChan
	wg.Wait()

}
