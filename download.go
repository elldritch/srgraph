package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff"
)

const errRateLimit = "Rate limited"

func setOver18Cookie(client *http.Client) (err error) {
	cookies := client.Jar
	if cookies == nil {
		cookies, err = cookiejar.New(nil)
		if err != nil {
			return
		}
	}

	url, err := url.Parse("https://www.reddit.com")
	if err != nil {
		return
	}

	cookies.SetCookies(url, []*http.Cookie{&http.Cookie{Name: "over18", Value: "1"}})
	client.Jar = cookies
	return nil
}

func getSubredditListing(client *http.Client, listing string) (data []byte, nextListing string, err error) {
	url := "https://www.reddit.com/reddits.json"
	if len(listing) > 0 {
		url += "?after=" + listing
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Subreddit Description Crawler")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if res.Body != nil {
			closeErr := res.Body.Close()
			if closeErr != nil {
				err = closeErr
			}
		}
	}()
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	srl := subredditListing{}
	err = json.Unmarshal(data, &srl)
	if err != nil {
		return
	}
	if srl.Error != "" {
		err = errors.New(errRateLimit)
		return
	}

	nextListing = srl.Data.After
	return
}

func saveSubredditListing(client *http.Client, directory string, listing string) (nextListing string, err error) {
	resp, nextListing, err := getSubredditListing(client, listing)
	if err != nil {
		return
	}

	if listing == "" {
		err = ioutil.WriteFile(filepath.Join(directory, "first"), resp, 0664)
	} else {
		err = ioutil.WriteFile(filepath.Join(directory, listing), resp, 0664)
	}

	return
}

func saveListingWithBackoff(client *http.Client, directory string, listing string) (nextListing string, err error) {
	logV("Getting listing \"" + listing + "\"...")
	backoffAlgorithm := backoff.NewExponentialBackOff()
	backoffAlgorithm.InitialInterval = time.Second
	backoffAlgorithm.MaxElapsedTime = 0
	_ = backoff.Retry(func() error {
		nextListing, err = saveSubredditListing(client, directory, listing)
		if err != nil && err.Error() == errRateLimit {
			return err
		}
		return nil
	}, backoffAlgorithm)
	return
}

func saveAllSubreddits(directory string, startPage string) (err error) {
	client := http.Client{}
	err = os.MkdirAll(directory, 0775)
	if err != nil {
		return
	}

	// Set cookie so we can crawl over18 subreddits as well
	err = setOver18Cookie(&client)
	if err != nil {
		return
	}

	logV("Done configuring client, downloading...")

	// First page
	nextListing := startPage
	lastListing := nextListing
	nextListing, err = saveListingWithBackoff(&client, directory, nextListing)
	if err != nil {
		return
	}
	defer func() {
		log.Printf("Last page downloaded: \"%s\"\n", lastListing)
	}()

	// All other pages
	for nextListing != "" {
		lastListing = nextListing
		nextListing, err = saveListingWithBackoff(&client, directory, nextListing)
	}
	return
}
