package board

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type catalog []struct {
	Page    int `json:"page"`
	Threads []struct {
		No            int64  `json:"no"`
		Sticky        int    `json:"sticky,omitempty"`
		Closed        int    `json:"closed,omitempty"`
		Now           string `json:"now"`
		Name          string `json:"name"`
		Sub           string `json:"sub,omitempty"`
		Com           string `json:"com"`
		Filename      string `json:"filename"`
		Ext           string `json:"ext"`
		W             int    `json:"w"`
		H             int    `json:"h"`
		TnW           int    `json:"tn_w"`
		TnH           int    `json:"tn_h"`
		Tim           int64  `json:"tim"`
		Time          int    `json:"time"`
		Md5           string `json:"md5"`
		Fsize         int    `json:"fsize"`
		Resto         int    `json:"resto"`
		Capcode       string `json:"capcode,omitempty"`
		SemanticURL   string `json:"semantic_url"`
		Replies       int    `json:"replies"`
		Images        int    `json:"images"`
		OmittedPosts  int    `json:"omitted_posts,omitempty"`
		OmittedImages int    `json:"omitted_images,omitempty"`
		LastReplies   []struct {
			No       int64  `json:"no"`
			Now      string `json:"now"`
			Name     string `json:"name"`
			Com      string `json:"com"`
			Filename string `json:"filename"`
			Ext      string `json:"ext"`
			W        int    `json:"w"`
			H        int    `json:"h"`
			TnW      int    `json:"tn_w"`
			TnH      int    `json:"tn_h"`
			Tim      int64  `json:"tim"`
			Time     int    `json:"time"`
			Md5      string `json:"md5"`
			Fsize    int    `json:"fsize"`
			Resto    int    `json:"resto"`
			Capcode  string `json:"capcode"`
		} `json:"last_replies"`
		LastModified int `json:"last_modified"`
		Bumplimit    int `json:"bumplimit,omitempty"`
		Imagelimit   int `json:"imagelimit,omitempty"`
	} `json:"threads"`
}

func getCatalog(ctx context.Context, client *http.Client, board string, ifModifiedSince time.Time) (c catalog, contentModified bool, err error) {
	url, err := url.JoinPath("https://a.4cdn.org", board, "catalog.json")
	if err != nil {
		return catalog{}, false, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return catalog{}, false, err
	}
	if !ifModifiedSince.IsZero() {
		r.Header.Set("If-Modified-Since", ifModifiedSince.Format(time.RFC1123))
	}

	resp, err := client.Do(r)
	if err != nil {
		return catalog{}, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return catalog{}, false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return catalog{}, false, fmt.Errorf("got %d statuscode", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return catalog{}, false, err
	}

	if err = json.Unmarshal(body, &c); err != nil {
		return catalog{}, false, err
	}

	return c, true, nil
}
