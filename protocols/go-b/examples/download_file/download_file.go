package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	exampleFileURL := "https://x.bitfs.network/6ce94f75b88a6c24815d480437f4f06ae895afdab8039ddec10748660c29f910.out.0.3"
	downloadDir := "downloads"

	// Build fileName from fullPath
	fileURL, err := url.Parse(exampleFileURL)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Set the file parts
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create the downloads dir
	if _, err = os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.Mkdir(downloadDir, 0777)
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
	}

	// Create blank file
	var file *os.File
	if file, err = os.Create(downloadDir + "/" + fileName); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Create the request to get the file contents
	var req *http.Request
	if req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, exampleFileURL, nil); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Make the actual request
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Close the body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Copy
	var size int64
	if size, err = io.Copy(file, resp.Body); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	defer func() {
		_ = file.Close()
	}()

	log.Println("url: ", exampleFileURL)
	log.Printf("downloaded a file %s with size %d\n\n", downloadDir+"/"+fileName, size)
}
