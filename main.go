package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	getCmd        = kingpin.Command("get", "Download subreddit listings. NOTE: This will take a long time (several hours).")
	getCmdListing = getCmd.Flag("start", "The ID of the listing to start at (the `after` attribute from reddits.json pagination).").Short('s').String()
	getCmdDir     = getCmd.Arg("data_directory", "The directory to save downloaded subreddits in.").Required().String()

	makeCmd    = kingpin.Command("make", "Generate a graph of subreddits by sidebar links from a directory of subreddit listings.")
	makeCmdDir = makeCmd.Arg("data_directory", "A directory with saved subreddits").Required().String()

	verbose = kingpin.Flag("verbose", "Verbose progress reports.").Short('v').Bool()
)

func logV(msg string) {
	if *verbose {
		log.Println(msg)
	}
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadParsed(directory string) (subreddits []subredditData, err error) {
	if isParseNeeded(directory) {
		subreddits, err = parse(directory)
	} else {
		logV("Loading from cache...")
		data, err := ioutil.ReadFile(filepath.Join(directory, parsedFile))
		if err != nil {
			return subreddits, err
		}
		err = json.Unmarshal(data, &subreddits)
		if err != nil {
			return subreddits, err
		}
		logV("Done loading from cache.")
	}
	return
}

func main() {
	var err error

	switch kingpin.Parse() {
	case "get":
		log.Println("Downloading subreddit listings...")
		err = saveAllSubreddits(*getCmdDir, *getCmdListing)
		log.Println("Finished.")
		die(err)
	case "make":
		log.Printf("Generating graph from subreddit listings in %s...\n", *makeCmdDir)

		subreddits, err := loadParsed(*makeCmdDir)
		die(err)

		graph, err := generate(subreddits)
		die(err)
		output, err := json.Marshal(graph)
		fmt.Println(string(output))

		log.Println("Finished.")
		die(err)
	}
}
