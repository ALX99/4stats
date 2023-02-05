package board

import (
	"context"
	"net/http"
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

func (b *Board) Name() string {
	return b.name
}
