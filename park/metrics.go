package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// ParkMetrics contains all metrics specific to the main park simulator
var metrics = struct {
	Money        prometheus.Gauge
	Time         prometheus.Gauge
	EntranceFee  prometheus.Gauge
	OpensAt      prometheus.Gauge
	ClosesAt     prometheus.Gauge
	IsParkClosed prometheus.Gauge
	Guests       prometheus.Gauge
}{
	Money: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "park_money",
		Help: "Current money amount in the park",
	}),

	Time: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "park_time",
		Help: "Current time in the park (Unix timestamp)",
	}),

	EntranceFee: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "park_entrance_fee",
		Help: "Current entrance fee",
	}),

	OpensAt: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "park_opens_at",
		Help: "Hour at which the park opens",
	}),

	ClosesAt: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "park_closes_at",
		Help: "Hour at which the park closes",
	}),

	IsParkClosed: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "park_is_closed",
		Help: "Whether the park is closed (1) or open (0)",
	}),

	Guests: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "park_guests",
		Help: "Current number of guests in the park",
	}),
}

// RegisterParkMetrics registers all park-specific metrics
func RegisterParkMetrics() {
	prometheus.MustRegister(metrics.Money)
	prometheus.MustRegister(metrics.Time)
	prometheus.MustRegister(metrics.EntranceFee)
	prometheus.MustRegister(metrics.OpensAt)
	prometheus.MustRegister(metrics.ClosesAt)
	prometheus.MustRegister(metrics.IsParkClosed)
	prometheus.MustRegister(metrics.Guests)
}
