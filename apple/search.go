package apple

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// AppItem represents an app in the Apple App Store
type AppItem struct {
	TrackID                   int        `json:"trackId"`
	TrackName                 string     `json:"trackName"`
	BundleID                  string     `json:"bundleId"`
	MinOS                     string     `json:"minimumOsVersion"`
	ArtistID                  int        `json:"artistId"`
	ArtistName                string     `json:"artistName"`
	Price                     float64    `json:"price"`
	Genres                    []string   `json:"genres"`
	PrimaryGenreID            int        `json:"primaryGenreId"`
	PrimaryGenre              string     `json:"primaryGenreName"`
	SellerName                string     `json:"sellerName"`
	Version                   string     `json:"version"`
	ReleaseNote               string     `json:"releaseNotes"`
	FileSizeBytes             string     `json:"fileSizeBytes"`
	ReleaseDate               *time.Time `json:"releaseDate"`
	CurrentVersionReleaseDate *time.Time `json:"currentVersionReleaseDate"`
}

type SearchOptions struct {
	Region string
	Query  string
	Limit  int
}

// Search searches for an app in the Apple App Store (login not required)
func (c *AppleClient) Search(opt SearchOptions) ([]AppItem, error) {
	if opt.Region == "" {
		if c.Cred == nil {
			opt.Region = "US"
		}
		opt.Region = c.Cred.Region
	}

	// API Documentation: https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI/Searching.html
	req := c.defaultRequest().
		SetQueryParams(map[string]string{
			"term":    opt.Query,
			"country": opt.Region,
			"entity":  "software",
			"media":   "software",
			"limit":   fmt.Sprintf("%d", opt.Limit),
		})

	resp, err := req.Get("https://itunes.apple.com/search")
	if err != nil {
		return nil, errors.Join(ErrHTTPError, err)
	}
	if resp.IsError() {
		return nil, errors.Join(ErrHTTPError, errors.New(resp.Status()))
	}

	var result struct {
		Results []AppItem `json:"results"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, errors.Join(ErrNoResults, err)
	}

	return result.Results, nil
}
