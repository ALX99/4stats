package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	server *http.Server

	// metrics
	ppm *prometheus.GaugeVec
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
	s.ppm = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fourchan",
			Subsystem: "posts",
			Name:      "per_minute",
			Help:      "Number of posts per minute",
		},
		[]string{"board"},
	)

	return prometheus.Register(s.ppm)
}

func (s *Server) Start() error {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	return nil
}

func (s *Server) SetPPM(board string, v float64) {
	gauge, err := s.ppm.GetMetricWithLabelValues(board)
	if err != nil {
		log.Println(err)
	}
	gauge.Set(v)
}
