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

	b := board.New("g", m)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := b.Update(); err != nil {
		log.Fatalln(err)
	}

	delay := 10 * time.Second
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		ticker.Reset(delay)
		select {
		case <-ticker.C:
			if err := b.Update(); err != nil {
				log.Fatalln(err)
			}

		case <-ctx.Done():
			return
		}
	}
}
