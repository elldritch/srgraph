# srgraph

This script crawls https://www.reddit.com/reddits.json, downloading all pages
and using them to generate a JSON graph of subreddits. This graph connects
subreddits via links to each other in the sidebar.

So far, I've been unsuccessful in visualising this graph because of its sheer
size (rendering headlessly with D3 causes Node to die after consuming too much
memory.) If you make any progress on this front, feel free to submit a PR.

## Installation
Run `go get -u github.com/ilikebits/srgraph`.

## Usage
Use `srgraph --help` for more information.

### `srgraph get [<flags>] <data_directory>`
Crawls reddit for subreddits, storing them in `<data_directory>`. Reddit
paginates its listings of subreddits. As `get` crawls, it will log the current
page being downloaded. Use flag `--start=ID` to specify a page to start at (e.g.
after stopping the script).

Crawling all subreddits takes about 8 hours on my connection (the bottleneck is
rate limiting). See the
[Releases](https://github.com/ilikebits/srgraph/releases) page for pre-loaded
data sets of subreddits.

### `srgraph make <data_directory>`
Uses downloaded subreddits in `<data_directory>` to generate a JSON graph, which
is emitted on `stdout`. NOTE: `<data_directory>` must not contain any files
other than the downloaded subreddits, or `make`'s behaviour is undefined (it
will probably crash, but it might do something weird if the extra file is valid
JSON).

This will cache the parsed list of subreddits in the current directory. In the
future, this may be flag-optional or cache in the data directory instead.

An example graph:
```javascript
{
  "nodes": [{
    "name": "/r/<subreddit>/", // string of the form "/r/<name>/"
    "over18": false, // boolean
    "private": false, // boolean
    "subscribers": 42 // non-negative integer
  }], // ...
  "edges": [{
    "source": 42, // index of source node
    "target": 314, // index of target node
    "count": 1, // how many times does source link to target?
  }]
}
```

## Development
Feel free to submit PRs! In particular, I think it would be interesting to allow
"transform" plugins (e.g. building a graph without any nodes of degree 0, or
building a graph that fulfils some other condition).

I use `glide` to manage dependencies. See
`src/github.com/ilikebits/srgraph/glide.yaml`.
