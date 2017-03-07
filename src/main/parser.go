package main

import "log"

/* ~~~ Fan Out ~~~ */
func ParseInput(chFanOut <- chan string, numParsers int) [] <- chan string {
	log.Println("Fanning out parsers...")

	helpers := make([] <- chan string, numParsers)

	// here, we send off a bunch of functions that all read from the same channel (to break up the work)
	// aka, fanning out

	for i := range helpers {
		helpers[i] = parseHelper(chFanOut)
	}

	return helpers

}

func parseHelper(input <- chan string) <- chan string {
	channel := make(chan string)

	go func() {

		for s := range input {
			if s == "the" {
				channel <- "potato"
			}
		}

		close(channel)

	}()

	return channel
}

/* ~~~ Fan In ~~~ */
func FanIn(incomingChannels [] <- chan string) <- chan string {
	log.Println("Fanning in parsers...")

	channel := make(chan string)
	count := 0

	for _, helperChannel := range incomingChannels {

		// here, we launch a bunch of functions to feed all the separate channels into ONE channel
		// aka, fanning in

		go func() {

			for s := range helperChannel {
				channel <- s
			}

			count++
			if (count == len(incomingChannels)) {
				close(channel)
			}

		}()

	}

	return channel
}