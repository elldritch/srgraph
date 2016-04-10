package main

import "encoding/json"

type subredditListing struct {
	Data  subredditListingData
	Kind  string
	Error string
}

type subredditListingData struct {
	Modhash  string
	Children []subreddit
	After    string
	Before   string
}

type subreddit struct {
	Kind string
	Data subredditData
}

type subredditData struct {
	Quarantine              bool
	HideAds                 bool
	UserSrThemeEnabled      bool
	WikiEnabled             bool
	CollapseDeletedComments bool
	Over18                  bool
	BannerImg               string
	SubmitTextHTML          string
	UserIsBanned            string
	ID                      string
	SubmitText              string
	DisplayName             string
	HeaderImg               string
	DescriptionHTML         string
	Title                   string
	PublicDescriptionHTML   string
	IconSize                []json.Number
	SuggestedCommentSort    string
	IconImg                 string
	HeaderTitle             string
	Description             string
	UserIsMuted             string
	SubmitLinkLabel         string
	AccountsActive          string
	PublicTraffic           string
	HeaderSize              []json.Number
	Subscribers             json.Number
	SubmitTextLabel         string
	Lang                    string
	UserIsModerator         string
	KeyColor                string
	Name                    string
	Created                 json.Number
	URL                     string
	CreatedUtc              json.Number
	BannerSize              []json.Number
	UserIsContributor       string
	PublicDescription       string
	CommentScoreHideMins    json.Number
	SubredditType           string
	SubmissionType          string
	UserIsSubscriber        string
}
