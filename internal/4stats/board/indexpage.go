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

type indexPage struct {
	Threads []struct {
		Posts []struct {
			No          int64  `json:"no"`
			Now         string `json:"now"`
			Name        string `json:"name"`
			Com         string `json:"com"`
			Filename    string `json:"filename,omitempty"`
			Ext         string `json:"ext,omitempty"`
			W           int    `json:"w,omitempty"`
			H           int    `json:"h,omitempty"`
			TnW         int    `json:"tn_w,omitempty"`
			TnH         int    `json:"tn_h,omitempty"`
			Tim         int64  `json:"tim,omitempty"`
			Time        int    `json:"time"`
			Md5         string `json:"md5,omitempty"`
			Fsize       int    `json:"fsize,omitempty"`
			Resto       int    `json:"resto"`
			Bumplimit   int    `json:"bumplimit,omitempty"`
			Imagelimit  int    `json:"imagelimit,omitempty"`
			SemanticURL string `json:"semantic_url,omitempty"`
			Replies     int    `json:"replies,omitempty"`
			Images      int    `json:"images,omitempty"`
		} `json:"posts"`
	} `json:"threads"`
}

func getIndexPage(ctx context.Context, client *http.Client, board string, pageNum uint8, ifModifiedSince time.Time) (page indexPage, contentModified bool, err error) {
	url, err := url.JoinPath("https://a.4cdn.org", board, fmt.Sprintf("%d.json", pageNum))
	if err != nil {
		return indexPage{}, false, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return indexPage{}, false, err
	}
	if !ifModifiedSince.IsZero() {
		r.Header.Set("If-Modified-Since", ifModifiedSince.Format(time.RFC1123))
	}

	resp, err := client.Do(r)
	if err != nil {
		return indexPage{}, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return indexPage{}, false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return indexPage{}, false, fmt.Errorf("got %d statuscode", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return indexPage{}, false, err
	}

	iPage := indexPage{}
	if err = json.Unmarshal(body, &iPage); err != nil {
		return indexPage{}, false, err
	}

	return iPage, true, nil
}
