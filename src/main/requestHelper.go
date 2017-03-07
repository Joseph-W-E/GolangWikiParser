package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"
	"log"
	"io/ioutil"
)

// struct for decoding the json request
type Wiki struct {
	URL     string `json:"url"`
}

// read the REST call and extract the url
func DecodeRequest(request *http.Request) string {
	log.Println("Decoding request...")

	decoder := json.NewDecoder(request.Body)

	var wiki Wiki

	err := decoder.Decode(&wiki)
	if err != nil {
		panic(err)
	}

	defer request.Body.Close()

	return wiki.URL
}

// make sure the url is a valid wikipedia link
func ValidateURL(url string) bool {
	var validUrl = regexp.MustCompile(`(https|http):\/\/en.wikipedia.org\/wiki\/[A-z_1-9]+`)

	return validUrl.MatchString(url)
}

// launch a goroutine to download the webpage and feed its data into a channel (read only)
func Tokenize(url string) <- chan string {
	channel := make(chan string)

	go func() {
		webpage := downloadToString(url)

		for _, slice := range strings.Split(webpage, " ") {
			channel <- slice
		}

		close(channel)
	}()

	return channel
}

// helper to download a webpage as a string
func downloadToString(url string) string {
	log.Println("Downloading webpage...")

	// download the web page
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := netClient.Get(url)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		panic(err)
	}
	defer response.Body.Close()

	// read the web page into a series of bytes
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading downloaded content", url, "-", err)
		panic(err)
	}

	// return the bytes as a string
	return string(bytes)
}