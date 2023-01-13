package board

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Board struct {
	name string
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
	return Board{name: name}
}

func (b *Board) Watch() {
	threadsAPI := fmt.Sprintf("https://a.4cdn.org/%s/catalog.json", b.name)

	for {
		resp, err := http.Get(threadsAPI)
		if err != nil {
			log.Println(err)
			continue
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		resp.Body.Close()

		threads := boardThreads{}
		if err = json.Unmarshal(body, &threads); err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("%+v", threads)
		time.Sleep(time.Minute)
	}
}
