package main

import (
	"log"
	"net/http"
)

// entry point
func main() {
	log.Println("Starting server...")

	// Tell the server what to do when a request is made
	http.HandleFunc("/", handleRequest)

	// Start listening on the port
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// handle a request from the socket
func handleRequest(w http.ResponseWriter, request *http.Request) {
	log.Println("Handling request...")

	url := DecodeRequest(request)

	if (ValidateURL(url)) {

		// the channel that our helpers will all read from
		chFanOut := Tokenize(url)

		// start the (10) helpers
		fanOutHelpers := ParseInput(chFanOut, 10)

		/* we have now fanned out, so now we can start fanning back in */

		// we can now merge our slice of channels
		chFanIn := FanIn(fanOutHelpers)

		x := 0
		for range chFanIn {
			x++
		}
		log.Println(x)
	}

}