package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Server struct {
	server *http.Server

	// metrics
	ppm   *prometheus.GaugeVec
	posts *prometheus.GaugeVec
}

func New() *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &Server{
		server: &http.Server{
			Addr:    ":8000",
			Handler: mux,
		},
	}
}

func (s *Server) InitializeMetrics() error {
	labels := []string{"board"}
	s.ppm = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fourchan",
			Subsystem: "posts",
			Name:      "per_minute",
			Help:      "Number of posts per minute",
		},
		labels,
	)

	s.posts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fourchan",
			Subsystem: "posts",
			Name:      "total",
			Help:      "Total number of posts",
		},
		labels,
	)

	if err := prometheus.Register(s.ppm); err != nil {
		return err
	}

	return prometheus.Register(s.posts)
}

func (s *Server) Start() error {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("Could not start metrics server")
		}
	}()
	return nil
}

func (s *Server) SetPPM(board string, v float64) {
	gauge, err := s.ppm.GetMetricWithLabelValues(board)
	if err != nil {
		log.Err(err).Msg("Could not increment PPM metric")
	}
	gauge.Set(v)
}

func (s *Server) SetPostCount(board string, v float64) {
	gauge, err := s.posts.GetMetricWithLabelValues(board)
	if err != nil {
		log.Err(err).Msg("Could not increment Posts metric")
	}
	gauge.Set(v)
}
