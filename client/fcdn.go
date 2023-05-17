package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func download(url string, wg *sync.WaitGroup, dest string) {
	defer wg.Done()

	startTime := time.Now()

	// Create a new HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading %s: %s\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Create a new file to save the downloaded data
	file, err := os.Create(dest)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer file.Close()

	// Copy the response body to the file while measuring speed
	_, err = io.Copy(file, io.TeeReader(resp.Body, &SpeedMeasurer{}))
	if err != nil {
		fmt.Printf("Error saving file: %s\n", err)
		return
	}

	elapsedTime := time.Since(startTime)
	downloadSpeed := float64(resp.ContentLength) / elapsedTime.Seconds() / 1024 / 1024 // Speed in MB/s

	fmt.Printf("Downloaded %s\n", url)
	fmt.Printf("Elapsed Time: %s\n", elapsedTime)
	fmt.Printf("Download Speed: %.2f MB/s\n", downloadSpeed)
}

type SpeedMeasurer struct {
	totalBytes int64
	lastTime   time.Time
}

func (sm *SpeedMeasurer) Write(p []byte) (n int, err error) {
	n = len(p)
	sm.totalBytes += int64(n)

	if sm.lastTime.IsZero() {
		sm.lastTime = time.Now()
	} else {
		elapsedTime := time.Since(sm.lastTime)
		if elapsedTime.Seconds() >= 1 {
			downloadSpeed := float64(sm.totalBytes) / elapsedTime.Seconds() / 1024 / 1024 // Speed in MB/s
			fmt.Printf("Current Download Speed: %.2f MB/s\n", downloadSpeed)
			sm.totalBytes = 0
			sm.lastTime = time.Now()
		}
	}

	return
}

func main() {
	// Parse command-line arguments
	var fileName string
	flag.StringVar(&fileName, "a", "file.txt", "Destination file name")
	flag.Parse()

	// Get the server URLs from command-line arguments
	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Println("Please provide server URLs")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go download(url, &wg, fileName)
	}

	wg.Wait()
	fmt.Println("All downloads completed")
}
