package api

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/google/go-querystring/query"
)

const (
	ApiBaseUri = "https://archive.ragtag.moe"
)

type SearchQuery struct {
	Query     string    `url:"q"`
	VideoId   string    `url:"v"`
	ChannelId string    `url:"channel_id"`
	Sort      SortBy    `url:"sort"`
	SortOrder SortOrder `url:"sort_order"`
	From      *int      `url:"from"`
	Size      *int      `url:"size"`
}

type SearchResult struct {
	Took     int          `json:"took"`
	TimedOut bool         `json:"timed_out"`
	Shards   ResultShards `json:"_shards"`
	Hits     ResultHits   `json:"hits"`
}

type ResultShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skiped     int `json:"skipped"`
	Failed     int `json:"failed"`
}

type ResultHits struct {
	Total    HitsTotal   `json:"total"`
	MaxScore float64     `json:"max_score"`
	Hits     []VideoData `json:"hits"`
}

type HitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type VideoData struct {
	Index  string      `json:"_index"`
	Id     string      `json:"_id"`
	Score  float64     `json:"_score"`
	Source VideoSource `json:"_source"`
}

type VideoSource struct {
	ChannelName       string          `json:"channel_name"`
	DriveBase         string          `json:"drive_base"`
	LikeCount         int             `json:"like_count"`
	Timestamps        VideoTimestamps `json:"timestamps"`
	Fps               int             `json:"fps"`
	Description       string          `json:"description"`
	Title             string          `json:"title"`
	Duration          int             `json:"duration"`
	ArchivedTimestamp string          `json:"archived_timestamp"`
	Width             int             `json:"width"`
	FormatId          string          `json:"format_id"`
	Files             []VideoFiles    `json:"files"`
	ChannelId         string          `json:"channel_id"`
	ViewCount         int             `json:"view_count"`
	DislikeCount      int             `json:"dislike_count"`
	VideoId           string          `json:"video_id"`
	UploadDate        string          `json:"upload_date"`
	Height            int             `json:"height"`
}

type VideoTimestamps struct {
	ActualStartTime    string `json:"actualStartTime"`
	PublishedAt        string `json:"publishedAt"`
	ScheduledStartTime string `json:"scheduledStartTime"`
	ActualEndTime      string `json:"actualEndTime"`
}

type VideoFiles struct {
	Size int    `json:"size"`
	Name string `json:"name"`
}

type Channel struct {
	ChannelId   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	VideosCount int    `json:"video_count"`
}

type SortBy struct {
	slug string
}

func (s SortBy) String() string {
	return s.slug
}

var (
	None              = SortBy{""}
	ArchivedTimestamp = SortBy{"archived_timestamp"}
	UploadDate        = SortBy{"upload_date"}
	Duration          = SortBy{"duration"}
	ViewCount         = SortBy{"view_count"}
	LikeCount         = SortBy{"like_count"}
	DislikeCount      = SortBy{"dislike_count"}
)

type SortOrder struct {
	slug string
}

func (s SortOrder) String() string {
	return s.slug
}

var (
	Ascending  = SortOrder{"asc"}
	Descending = SortOrder{"desc"}
)

func ApiSearch(q SearchQuery) (*SearchResult, error) {
	val, err := query.Values(q)
	if err != nil {
		return nil, err
	}

	var json *SearchResult

	if err := requests.
		URL(ApiBaseUri).
		Path("/api/v1/search").
		Param("q", val.Get("q")).
		Param("v", val.Get("v")).
		Param("channel_id", val.Get("channel_id")).
		Param("sort", val.Get("sort")).
		Param("sort_order", val.Get("sort_order")).
		Param("from", val.Get("from")).
		Param("size", val.Get("size")).
		ToJSON(&json).
		Fetch(context.Background()); err != nil {
		return nil, err
	}

	return json, nil
}

func ApiChannel(channel string) (*SearchResult, error) {
	return ApiSearch(SearchQuery{
		ChannelId: channel,
	})
}

func ApiChannels() (*[]Channel, error) {
	var json *[]Channel

	if err := requests.URL(ApiBaseUri).
		Path("/api/v1/channels").
		ToJSON(&json).
		Fetch(context.Background()); err != nil {
		return nil, err
	}

	return json, nil
}
