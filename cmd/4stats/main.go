package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/alx99/4stats/internal/4stats/board"
	"github.com/alx99/4stats/internal/4stats/metrics"
	"github.com/rs/zerolog"
)

func main() {
	setupLogger()

	m := metrics.New()
	if err := m.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := m.InitializeMetrics(); err != nil {
		log.Fatal().Err(err).Send()
	}

	boards, err := board.GetBoards()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	bs := make([]board.Board, 0, len(boards.Boards))
	for _, b := range boards.Boards {
		bs = append(bs, board.New(b.Board, m))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	delay := time.Second
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i := 0; i < len(bs); i++ {
				if err := bs[i].Update(ctx); err != nil {
					log.Err(err).Msgf("Failed to update board %s", bs[i].Name())
				}
				time.Sleep(delay)
			}

		case <-ctx.Done():
			return
		}
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stdout,
		FormatCaller: func(i interface{}) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},
	}).
		With().
		Caller().
		Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
