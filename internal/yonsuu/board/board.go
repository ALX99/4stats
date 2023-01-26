package board

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alx99/yonsuu/internal/yonsuu/metrics"
)

type Board struct {
	name string

	prevFirstPage       indexPage
	prevFirstPageScrape time.Time

	m       metrics.Metrics
	running atomic.Bool
	stop    chan any
	sync.WaitGroup
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
			No          int    `json:"no"`
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
	}
}

func (b *Board) StartWatch() error {
	b.stop = make(chan any)
	b.Add(1)
	delay := 10 * time.Second

	if b.running.Swap(true) {
		return errors.New("board is already running")
	}

	go func() {
		ticker := time.NewTicker(delay)
		defer b.Done()
		defer ticker.Stop()

		if err := b.refresh(); err != nil {
			log.Fatalln(err)
		}

		for {
			ticker.Reset(delay)
			select {
			case <-ticker.C:
				if err := b.refresh(); err != nil {
					log.Fatalln(err)
				}

			case <-b.stop:
				return
			}
		}
	}()

	return nil
}

// refresh calls the API and refreshes the statistics
func (b *Board) refresh() error {
	ppm, err := b.calculatePPM()
	if err != nil {
		return err
	}

	b.m.SetPPM(b.name, ppm)
	log.Println("ppm:", ppm)

	return nil
}

func (b *Board) calculatePPM() (float64, error) {
	// save last valuse
	prevScrape := b.prevFirstPageScrape
	prevFirstIndexPage := b.prevFirstPage

	firstPage, err := b.getIndexPage(1)
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
		return threadList{}, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return threadList{}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return threadList{}, nil
	}
	resp.Body.Close()

	tList := threadList{}
	if err = json.Unmarshal(body, &tList); err != nil {
		return threadList{}, nil
	}

	return tList, nil
}

func (b *Board) getIndexPage(page uint8) (indexPage, error) {
	url, err := url.JoinPath("https://a.4cdn.org", b.name, fmt.Sprintf("%d.json", page))
	if err != nil {
		return indexPage{}, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return indexPage{}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return indexPage{}, nil
	}
	resp.Body.Close()

	iPage := indexPage{}
	if err = json.Unmarshal(body, &iPage); err != nil {
		return indexPage{}, nil
	}

	return iPage, nil
}

func (b *Board) StopWatch() {
	close(b.stop)
	b.Wait()
}
