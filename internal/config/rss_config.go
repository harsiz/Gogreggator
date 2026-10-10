package config

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	client := http.DefaultClient
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	byteResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var r RSSFeed
	if err = xml.Unmarshal(byteResp, &r); err != nil {
		return &RSSFeed{}, err
	}

	// Copy Paste Boilerplate
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	for i := range r.Channel.Item {
		item := &r.Channel.Item[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &r, nil
}
