package metrics

type Metrics interface {
	SetPPM(board string, v float64)
}
