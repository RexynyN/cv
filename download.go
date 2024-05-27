// package main

// import (
// 	"archive/zip"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// 	"time"
// )

// func downloadFile(url string, wg *sync.WaitGroup) {
// 	retry := 1
// 	for retry <= 3 {
// 		fmt.Println("Try ", retry)
// 		response, err := http.Get(url)
// 		if err != nil || response.StatusCode != 200 {
// 			fmt.Println("Error on request: ", err, response.StatusCode)
// 			time.Sleep(10)
// 			retry++
// 			continue
// 		}
// 		defer response.Body.Close()

// 		filename := getFileName(url)
// 		out, err := os.Create(filename)
// 		if err != nil {
// 			fmt.Println("Error creating a file: ", err)
// 			time.Sleep(10)
// 			retry++
// 			continue
// 		}
// 		defer out.Close()

// 		_, err = io.Copy(out, response.Body)
// 		if err != nil {
// 			fmt.Println("Error writing file: ", err)
// 			time.Sleep(10)
// 			retry++
// 			continue
// 		}
// 		break
// 	}

// 	wg.Done()
// }

// func getFileName(urlstr string) string {
// 	u, err := url.Parse(urlstr)
// 	if err != nil {
// 		log.Fatal("Error due to parsing url: ", err)
// 	}

// 	ex, _ := url.QueryUnescape(u.EscapedPath())
// 	return filepath.Base(ex)
// }

// func main() {
// 	urls := []string{
// 		"https://pbs.twimg.com/media/GLAhmB_WcAAlVvL.jpg",
// 		"https://pbs.twimg.com/media/GK_4xkOXsAA2wQW.jpg",
// 		"https://pbs.twimg.com/media/GK_9MM_Wo3DxG.jpg", // This one doesn't exist, must error
// 		"https://pbs.twimg.com/media/GK_9MM_WoAA3DxG.jpg",
// 		"https://pbs.twimg.com/media/GLx-gdRXEAAtHiD.jpg",
// 		"https://pbs.twimg.com/media/GLuQlhgXYAA7EbZ.jpg",
// 	}
// 	var wg sync.WaitGroup

// 	for _, url := range urls {
// 		wg.Add(1)
// 		go downloadFile(url, &wg)
// 	}

// 	wg.Wait()

// 	archive, err := os.Create("archive.cbz")
// 	if err != nil {
// 		panic(err)
// 	}
// 	zipWriter := zip.NewWriter(archive)
// 	defer archive.Close()

// 	files, err := os.ReadDir(".")
// 	if err != nil {
// 		panic(err)
// 	}

// 	for idx, file := range files {
// 		fmt.Println(file.Name())
// 		if !strings.HasSuffix(file.Name(), ".jpg") {
// 			continue
// 		}

// 		stream, err := os.Open(file.Name())
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer stream.Close()

// 		foil := fmt.Sprintf("foito%d.jpg", idx)
// 		fmt.Println(foil)
// 		writer, err := zipWriter.Create(foil)
// 		if err != nil {
// 			panic(err)
// 		}
// 		if _, err := io.Copy(writer, stream); err != nil {
// 			panic(err)
// 		}
// 	}
// 	zipWriter.Close()
// }

// // package main

// // import (
// // 	"fmt"
// // 	"io"
// // 	"log"
// // 	"net/http"
// // 	"os"
// // 	"strings"
// // 	"sync"
// // )

// // func goDownload(url string, ch chan<- string, wg *sync.WaitGroup) {
// // 	defer wg.Done()
// // 	// Get the image
// // 	response, err := http.Get(url)
// // 	fmt.Println(url, response.StatusCode)
// // 	if err != nil || response.StatusCode != 200 {
// // 		ch <- url
// // 		return
// // 	}
// // 	defer response.Body.Close()

// // 	tokens := strings.Split(url, "/")
// // 	fileName := tokens[len(tokens)-1]

// // 	// Create the file
// // 	file, err := os.Create(fileName)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer file.Close()

// // 	// Use io.Copy to just dump the response body to the file. This supports huge files
// // 	_, err = io.Copy(file, response.Body)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	fmt.Println("Downloaded", fileName, response.Body)
// // 	ch <- ""
// // }

// // func main() {
// // 	urls := []string{
// // 		"https://pbs.twimg.com/media/GLAhmB_WcAAlVvL.jpg",
// // 		"https://pbs.twimg.com/media/GK_4xkOXsAA2wQW.jpg",
// // 		"https://pbs.twimg.com/media/GK_9MM_Wo3DxG.jpg", // This one doesn't exist, must error
// // 		"https://pbs.twimg.com/media/GK_9MM_WoAA3DxG.jpg",
// // 	}

// // 	DownloadBatch(urls)
// // }

// // func GetLeftovers(ch chan string, left chan []string, wg *sync.WaitGroup) {
// // 	defer wg.Done()
// // 	leftOvers := make([]string, 0)
// // 	for val := range ch {
// // 		if val != "" {
// // 			leftOvers = append(leftOvers, val)
// // 		}
// // 	}
// // 	fmt.Println(leftOvers)
// // 	left <- leftOvers
// // }

// // func DownloadBatch(urls []string) {
// // 	urlsChan := make(chan string)
// // 	leftoversChan := make(chan []string)
// // 	var wg sync.WaitGroup

// // 	for _, url := range urls {
// // 		wg.Add(1)
// // 		go goDownload(url, urlsChan, &wg)
// // 	}
// // 	wg.Wait()

// // 	wg.Add(1)
// // 	go GetLeftovers(urlsChan, leftoversChan, &wg)
// // 	urls = <-leftoversChan
// // 	wg.Wait()

// // }
