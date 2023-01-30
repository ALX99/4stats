package board

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/alx99/yonsuu/internal/yonsuu/metrics"
)

type Board struct {
	name string

	prevFirstPage       indexPage
	prevFirstPageScrape time.Time

	m      metrics.Metrics
	client http.Client
}

type threadList []struct {
	Page    int `json:"page"`
	Threads []struct {
		No           int `json:"no"`
		LastModified int `json:"last_modified"`
		Replies      int `json:"replies"`
	} `json:"threads"`
}

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

func New(name string, m metrics.Metrics) Board {
	return Board{
		name: name,
		m:    m,
		client: http.Client{
			Transport: http.DefaultTransport,
			Timeout:   10 * time.Second,
		},
	}
}

// Update this board's metrics
func (b *Board) Update(ctx context.Context) error {
	ppm, err := b.calculatePPM(ctx)
	if err != nil {
		return err
	}

	b.m.SetPPM(b.name, ppm)

	return nil
}

func (b *Board) calculatePPM(ctx context.Context) (float64, error) {
	// save last valuse
	prevScrape := b.prevFirstPageScrape
	prevFirstIndexPage := b.prevFirstPage

	firstPage, err := b.getFirstIndexPage(ctx)
	if err != nil {
		return 0, err
	}

	// save new ones
	b.prevFirstPageScrape = time.Now()
	b.prevFirstPage = firstPage

	if prevScrape.IsZero() {
		return 0, nil // first scrape nothing to calculate
	}

	newPostsCount := getHighestPostNo(firstPage) - getHighestPostNo(prevFirstIndexPage)
	return (float64(newPostsCount) / time.Since(prevScrape).Seconds()) * 60, nil
}

func (b *Board) getThreadList() (threadList, error) {
	url, err := url.JoinPath("https://a.4cdn.org", b.name, "threads.json")
	if err != nil {
		return threadList{}, err
	}

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return threadList{}, err
	}

	resp, err := b.client.Do(r)
	if err != nil {
		return threadList{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return threadList{}, err
	}
	resp.Body.Close()

	tList := threadList{}
	if err = json.Unmarshal(body, &tList); err != nil {
		return threadList{}, err
	}

	return tList, nil
}

func (b *Board) getFirstIndexPage(ctx context.Context) (indexPage, error) {
	page := 1
	url, err := url.JoinPath("https://a.4cdn.org", b.name, fmt.Sprintf("%d.json", page))
	if err != nil {
		return indexPage{}, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return indexPage{}, err
	}
	if !b.prevFirstPageScrape.IsZero() {
		r.Header.Set("If-Modified-Since", b.prevFirstPageScrape.Format(time.RFC1123))
	}

	resp, err := b.client.Do(r)
	if err != nil {
		return indexPage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return b.prevFirstPage, nil
	}

	if resp.StatusCode != http.StatusOK {
		return indexPage{}, fmt.Errorf("got %d statuscode", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return indexPage{}, err
	}

	iPage := indexPage{}
	if err = json.Unmarshal(body, &iPage); err != nil {
		return indexPage{}, err
	}

	return iPage, nil
}

func (b *Board) Name() string {
	return b.name
}
