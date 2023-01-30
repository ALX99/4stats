package board

import (
	"context"
	"encoding/json"
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
	page, modified, err := getIndexPage(ctx, &b.client, b.name, 1, b.prevFirstPageScrape)
	if err != nil {
		return indexPage{}, err
	}
	if !modified {
		return b.prevFirstPage, nil
	}
	return page, nil
}

func (b *Board) Name() string {
	return b.name
}
