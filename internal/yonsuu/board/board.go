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
	name          string
	apiThreadList string

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

func New(name string) Board {
	return Board{
		name:          name,
		apiThreadList: fmt.Sprintf("https://a.4cdn.org/%s/threads.json", name),
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

		if err := b.refresh(); err != nil {
			fmt.Println(err)
		}

		for {
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
	tList, err := b.getThreadList()
	if err != nil {
		return err
	}

	fmt.Printf("%+v", tList)
	return nil
}

func (b *Board) getThreadList() (threadList, error) {
	resp, err := http.Get(b.apiThreadList)
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

func (b *Board) StopWatch() {
	close(b.stop)
	b.Wait()
}
