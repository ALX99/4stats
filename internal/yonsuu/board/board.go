package board

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Board struct {
	name   string
	apiURL string

	running atomic.Bool
	stop    chan any
	sync.WaitGroup
}

type boardThreads []struct {
	Page    int `json:"page"`
	Threads []struct {
		No           int `json:"no"`
		LastModified int `json:"last_modified"`
		Replies      int `json:"replies"`
	} `json:"threads"`
}

func New(name string) Board {
	return Board{
		name:   name,
		apiURL: fmt.Sprintf("https://a.4cdn.org/%s/catalog.json", name),
	}
}

func (b *Board) StartWatch() error {
	b.stop = make(chan any)
	b.Add(1)
	delay := time.Minute

	if b.running.Swap(true) {
		return errors.New("board is already running")
	}

	go func() {
		ticker := time.NewTicker(delay)
		defer b.Done()
		defer ticker.Stop()

		for {
			if err := b.refresh(); err != nil {
				fmt.Println(err)
			}

			ticker.Reset(delay)
			select {
			case <-ticker.C:
				if err := b.refresh(); err != nil {
					fmt.Println(err)
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
	resp, err := http.Get(b.apiURL)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	threads := boardThreads{}
	if err = json.Unmarshal(body, &threads); err != nil {
		return err
	}

	fmt.Printf("%+v", threads)

	return nil
}

func (b *Board) StopWatch() {
	close(b.stop)
	b.Wait()
}
