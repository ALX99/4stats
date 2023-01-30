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

	"github.com/alx99/yonsuu/internal/yonsuu/board"
	"github.com/alx99/yonsuu/internal/yonsuu/metrics"
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
		if b.Board != "vg" { // todo something not working with /vg/
			bs = append(bs, board.New(b.Board, m))
		}
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
				log.Debug().Msgf("Updating board %s", bs[i].Name())
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
