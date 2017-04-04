package main

import (
	"log"
	"sync"
	"regexp"
	"time"
)

/* ~~~ Fan Out ~~~ */
func ParseInput(chTokens <- chan string, numParsers int) [] <- chan string {
	log.Println("Fanning out parsers...")

	helpers := make([] <- chan string, numParsers)

	// here, we send off a bunch of functions that all read from the same channel (to break up the work)
	// aka, fanning out

	for i := range helpers {
		helpers[i] = parseHelper(chTokens)
	}

	return helpers
}

func parseHelper(input <- chan string) <- chan string {
	channel := make(chan string)

	go func() {

		for s := range input {

			// simulate taking awhile
			time.Sleep(100)

			expression := regexp.MustCompile("[^A-z]+")
			charactersOnly := expression.ReplaceAllString(s, "")

			if (charactersOnly != "") {
				channel <- charactersOnly
			}
		}

		close(channel)

	}()

	return channel
}

/* ~~~ Fan In ~~~ */
func FanIn(incomingChannels [] <- chan string) <- chan string {
	log.Println("Fanning in parsers...")

	var waitGroup sync.WaitGroup
	chMerged := make(chan string)

	// this is our function to pull data from a helper into the merged channel
	funnel := func(chHelper <- chan string) {
		for val := range chHelper {
			chMerged <- val
		}
		waitGroup.Done()
	}

	// once .Done() is called len(incomingChannels) times, .Wait() will no longer block
	waitGroup.Add(len(incomingChannels))

	// start funneling all the helpers' channels into chMerged
	for _, channel := range incomingChannels {
		go funnel(channel)
	}

	// wait for the helpers to finish
	go func() {
		waitGroup.Wait()
		close(chMerged)
	}()

	return chMerged
}

/* ~~~ Final Operation ~~~ */
func WordCount(chData <- chan string, bound int) {
	counts := make(map[string]int)

	for data := range chData {
		counts[data]++
	}

	for s, i := range counts {
		if i < bound {
			delete(counts, s)
		}
	}

	log.Println(counts)
}