package main

import (
	"errors"
	"regexp"
	"strings"
)

type subredditGraph struct {
	Nodes []subredditGraphNode `json:"nodes"`
	Edges []subredditGraphEdge `json:"edges"`
}

type subredditGraphNode struct {
	Name        string `json:"name"`
	Over18      bool   `json:"over18"`
	Private     bool   `json:"private"`
	Subscribers int64  `json:"subscribers"`
}

type subredditGraphEdge struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Count  int `json:"count"`
}

func processNodes(subreddits []subredditData) (nodes []subredditGraphNode, indices map[string]int, err error) {
	// Allocate index map
	indices = make(map[string]int)

	for i, subreddit := range subreddits {
		// Add to index map
		indices[strings.ToLower(subreddit.URL)] = i

		// Add to list of nodes
		subscribers, err := subreddit.Subscribers.Int64()
		if err != nil {
			return nodes, indices, err
		}
		nodes = append(nodes, subredditGraphNode{
			Name:        subreddit.URL,
			Over18:      subreddit.Over18,
			Subscribers: subscribers,
			Private:     false,
		})
	}

	return
}

func addPrivateSubreddits(subreddits []subredditData, nodes *[]subredditGraphNode, indices *map[string]int) {
	for _, subreddit := range subreddits {
		// Go through links
		links := rslash.FindAllString(subreddit.Description, -1)
		for _, link := range links {
			// Sometimes, we'll find links to private subreddits that didn't show up
			// while crawling the subreddits endpoint
			_, ok := (*indices)[strings.ToLower(link)+"/"]
			if !ok {
				(*indices)[strings.ToLower(link)+"/"] = len(*nodes)

				*nodes = append(*nodes, subredditGraphNode{
					Name:        link + "/",
					Over18:      false,
					Subscribers: 0,
					Private:     true,
				})
			}
		}
	}
}

func processEdges(subreddits []subredditData, indices map[string]int) (edges []subredditGraphEdge, err error) {
	for i, subreddit := range subreddits {
		// Count links to weight edges
		linkCount := make(map[string]int)
		links := rslash.FindAllString(subreddit.Description, -1)
		for _, link := range links {
			linkCount[strings.ToLower(link)]++
		}

		// Go through links
		for link, count := range linkCount {
			target, ok := indices[strings.ToLower(link)+"/"]
			if ok {
				// Add edges
				edges = append(edges, subredditGraphEdge{
					Source: i,
					Target: target,
					Count:  count,
				})
			} else {
				return nil, errors.New("Failed to find \"" + link + "/" + "\"\n")
			}
		}
	}
	return
}

var rslash = regexp.MustCompile(`\/r\/\w+`)

func generate(subreddits []subredditData) (graph subredditGraph, err error) {
	logV("Building graph...")

	var nodes []subredditGraphNode
	var edges []subredditGraphEdge
	indices := make(map[string]int)

	// Build index map and add nodes
	logV("Building nodes...")
	nodes, indices, err = processNodes(subreddits)
	addPrivateSubreddits(subreddits, &nodes, &indices)

	// Add edges
	logV("Building edges...")
	edges, err = processEdges(subreddits, indices)

	// Build graph
	graph = subredditGraph{Nodes: nodes, Edges: edges}
	logV("Done building graph.")

	return
}
