package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/alx99/yonsuu/internal/yonsuu/board"
	"github.com/alx99/yonsuu/internal/yonsuu/metrics"
)

func main() {
	m := metrics.New()
	if err := m.Start(); err != nil {
		log.Fatalln(err)
	}

	if err := m.InitializeMetrics(); err != nil {
		log.Fatalln(err)
	}

	boards, err := board.GetBoards()
	if err != nil {
		log.Fatalln(err)
	}

	bs := make([]board.Board, 0, len(boards.Boards))
	for _, b := range boards.Boards {
		if b.Board != "vg" { // todo something not working with /vg/
			bs = append(bs, board.New(b.Board, m))
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	delay := 10 * time.Second
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		ticker.Reset(delay)
		select {
		case <-ticker.C:
			for i := 0; i < len(bs); i++ {
				if err := bs[i].Update(); err != nil {
					log.Fatalln(err)
				}
				time.Sleep(2 * time.Second)
			}

		case <-ctx.Done():
			return
		}
	}
}
