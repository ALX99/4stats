package board

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/alx99/4stats/internal/4stats/catalog"
	"github.com/alx99/4stats/internal/4stats/metrics"
)

type Board struct {
	name string

	catalog *catalog.Catalog

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
	if b.catalog == nil {
		catalog, err := catalog.GetCatalog(ctx, b.client, b.name)
		if err != nil {
			return err
		}
		b.catalog = &catalog
	} else if err := b.catalog.Update(ctx, b.client); err != nil {
		return err
	}

	// todo /vg/ not working due to last_replies not
	// being populated
	if b.name != "vg" {
		b.m.SetPPM(b.name, b.catalog.GetPPM())
	}
	b.m.SetPostCount(b.name, float64(b.catalog.GetPostCount()))
	b.m.SetImageCount(b.name, float64(b.catalog.GetImageCount()))

	return nil
}

func (b *Board) getThreadList() (threadList, error) {
	url, err := url.JoinPath("http://a.4cdn.org", b.name, "threads.json")
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

func (b *Board) Name() string {
	return b.name
}
