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

func (s *Server) SetPostCount(board string, v float64) {
	gauge, err := s.posts.GetMetricWithLabelValues(board)
	if err != nil {
		log.Println(err)
	}
	gauge.Set(v)
}
