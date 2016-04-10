package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func isParseNeeded(directory string) bool {
	_, err := os.Stat(filepath.Join(directory, parsedFile))
	return err != nil
}

func parseListing(data []byte, subreddits *[]subredditData) (err error) {
	var listing subredditListing
	err = json.Unmarshal(data, &listing)
	if err != nil {
		return
	}

	for _, child := range listing.Data.Children {
		*subreddits = append(*subreddits, child.Data)
	}
	return
}

func marshalToFile(path string, subreddits *[]subredditData) (err error) {
	jsonSubreddits, err := json.Marshal(*subreddits)
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Dir(path), 0775)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, jsonSubreddits, 0664)
	if err != nil {
		return
	}
	return
}

// We use a different directory because most filesystems have O(n) lookup with
// large numbers of files in them, and the point of serialising results is so we
// can quickly tell whether or not we can skip parsing
var parsedFile = filepath.Join("..", "srgraph-cache.json")

func parse(directory string) (subreddits []subredditData, err error) {
	logV("Loading listing files...")
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}
	logV("Done loading listing files.")

	// Loop through data files
	logV("Parsing data files...")
	for i, file := range files {
		contents, err := ioutil.ReadFile(filepath.Join(directory, file.Name()))
		if err != nil {
			return subreddits, err
		}

		err = parseListing(contents, &subreddits)
		if err != nil {
			return subreddits, err
		}

		if i != 0 && i%1000 == 0 && *verbose {
			log.Printf("Done with file %d\n", i)
		}
	}
	logV("Done parsing.")

	// Serialise the parsed result, because otherwise this step is extremely slow
	// on old filesystems
	err = marshalToFile(filepath.Join(directory, parsedFile), &subreddits)
	return
}
