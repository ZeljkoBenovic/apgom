package asterisk

import (
	"time"

	"github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsAsterisk struct {
	activeCalls     prometheus.Gauge
	totalCalls      prometheus.Gauge
	asteriskScraper *asterisk.AsteriskScraper
}

func NewMetricsAsterisk(asteriskScraper *asterisk.AsteriskScraper) *MetricsAsterisk {
	return &MetricsAsterisk{
		asteriskScraper: asteriskScraper,
	}
}

func (m *MetricsAsterisk) ActiveCalls() {
	m.activeCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "asterisk",
		Name:        "active_calls",
		Help:        "The number of active calls",
		ConstLabels: nil,
	})
}

func (m *MetricsAsterisk) TotalProcessedCalls() {
	m.totalCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "asterisk",
		Name:        "total_calls",
		Help:        "The number of total processed calls",
		ConstLabels: nil,
	})
}

func (m *MetricsAsterisk) RunAsteriskMetricsCollector() {
	//TODO: add context
	for {
		select {
		case <-time.After(time.Second):
			m.activeCalls.Set(m.asteriskScraper.GetActiveCalls())
			m.totalCalls.Set(m.asteriskScraper.GetTotalProcessedCalls())
		}
	}
}
