package main

import (
	"github.com/alx99/yonsuu/internal/yonsuu/board"
	"github.com/alx99/yonsuu/internal/yonsuu/metrics"
)

func main() {
	m := metrics.New()
	m.Start()
	b := board.New("po")
	b.Watch()
}
