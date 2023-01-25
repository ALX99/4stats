package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/alx99/yonsuu/internal/yonsuu/board"
	"github.com/alx99/yonsuu/internal/yonsuu/metrics"
)

func main() {
	m := metrics.New()
	if err := m.Start(); err != nil {
		log.Fatalln(err)
	}

	b := board.New("po")
	if err := b.StartWatch(); err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	cancel()
}
