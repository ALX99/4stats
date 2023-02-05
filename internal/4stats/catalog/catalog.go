package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Catalog struct {
	scrapedAt     time.Time
	prevScrapedAt time.Time
	url           string

	prevCatalogPages []CatalogPage
	CatalogPages     []CatalogPage
}

type CatalogPage struct {
	Page    int      `json:"page"`
	Threads []Thread `json:"threads"`
}

func GetCatalog(ctx context.Context, client http.Client, board string) (Catalog, error) {
	url, err := url.JoinPath("http://a.4cdn.org", board, "catalog.json")
	if err != nil {
		return Catalog{}, err
	}
	c := Catalog{url: url}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Catalog{}, err
	}

	resp, err := client.Do(r)
	if err != nil {
		return Catalog{}, err
	}
	defer resp.Body.Close()
	scrapedAt := time.Now()

	if resp.StatusCode != http.StatusOK {
		return Catalog{}, fmt.Errorf("got %d statuscode", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Catalog{}, err
	}

	if err = json.Unmarshal(body, &c.CatalogPages); err != nil {
		return Catalog{}, err
	}
	c.scrapedAt = scrapedAt
	c.prevCatalogPages = c.CatalogPages // initialize to the same value

	return c, nil
}

func (c *Catalog) Update(ctx context.Context, client http.Client) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return err
	}

	r.Header.Set("If-Modified-Since", c.scrapedAt.Format(time.RFC1123))

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	scrapedAt := time.Now()

	// Nothing to update
	if resp.StatusCode == http.StatusNotModified {
		c.prevScrapedAt = c.scrapedAt
		c.scrapedAt = scrapedAt
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d statuscode", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tmpPrevCatalogPages := c.prevCatalogPages
	c.prevCatalogPages = c.CatalogPages
	c.CatalogPages = make([]CatalogPage, 0) // Allocate new heap memory

	if err = json.Unmarshal(body, &c.CatalogPages); err != nil {
		// restore previous values in case of error
		c.CatalogPages = c.prevCatalogPages
		c.prevCatalogPages = tmpPrevCatalogPages
		return err
	}

	c.prevScrapedAt = c.scrapedAt
	c.scrapedAt = scrapedAt

	return nil
}

func (c *Catalog) GetPostCount() int {
	count := 0
	for _, page := range c.CatalogPages {
		for _, thread := range page.Threads {
			count += thread.Replies + 1
		}
	}
	return count
}

func (c *Catalog) GetImageCount() int {
	count := 0
	for _, page := range c.CatalogPages {
		for _, thread := range page.Threads {
			if thread.Ext != "" {
				count++ // Thread has an image
			}
			count += thread.Images
		}
	}
	return count
}

func (c *Catalog) GetPPM() float64 {
	newPostsCount := getHighestPostNo(c.CatalogPages) - getHighestPostNo(c.prevCatalogPages)
	return (float64(newPostsCount) / time.Since(c.prevScrapedAt).Seconds()) * 60
}

func getHighestPostNo(pages []CatalogPage) int64 {
	var maxNo int64
	for _, page := range pages {
		for _, thread := range page.Threads {
			if maxNo < thread.No {
				maxNo = thread.No
			}
			for _, reply := range thread.LastReplies {
				if maxNo < reply.No {
					maxNo = reply.No
				}
			}
		}
	}
	return maxNo
}
