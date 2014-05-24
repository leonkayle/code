package search

import (
	"log"
	"sync"
)

// A map of registered matchers for searching.
var matchers map[string]Matcher = make(map[string]Matcher)

// Register is called to register a matcher for use
// by the program.
func Register(feedType string, matcher Matcher) {
	// Assign the matcher to the specified key.
	log.Println("Register", feedType)
	matchers[feedType] = matcher
}

// Run performs the search logic.
func Run(searchTerm string) {
	// RetrieveFeed returns the list of feeds to search through.
	feeds, err := RetrieveFeed()
	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to receive match results to display.
	results := make(chan *Result)

	// Setup a wait group so we can process all the feeds.
	var waitGroup sync.WaitGroup

	// Set the number of goroutines we need to wait for while
	// they process the individual feeds.
	waitGroup.Add(len(feeds))

	// Launch a goroutine for each feed to find the results.
	for _, feed := range feeds {
		// Retrieve a matcher for the search.
		matcher, ok := matchers[feed.Type]
		if !ok {
			matcher = matchers["default"]
		}

		// Make a copy of the value to give each goroutine
		// their own copy.
		find := feed

		// Launch the goroutine to perform the search.
		go func() {
			Match(matcher, find, searchTerm, results)
			waitGroup.Done()
		}()
	}

	// Launch a groutine to monitor when all the work is done.
	go func() {
		// Wait for everything to be processed.
		waitGroup.Wait()

		// Close the channel to signal to the Display
		// function that we can exit the program.
		close(results)
	}()

	// Start displaying results as they are avaiable and
	// return after the final result is displayed.
	Display(results)
}
