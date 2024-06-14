package asterisk

import (
	"context"
	"time"

	"github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsAsterisk struct {
	ctx context.Context

	activeCalls      prometheus.Gauge
	totalCalls       prometheus.Gauge
	totalPeers       prometheus.Gauge
	availablePeers   prometheus.Gauge
	unavailablePeers prometheus.Gauge
	asteriskScraper  *asterisk.AsteriskScraper
}

func NewMetricsAsterisk(ctx context.Context, asteriskScraper *asterisk.AsteriskScraper) *MetricsAsterisk {
	return &MetricsAsterisk{
		ctx:             ctx,
		asteriskScraper: asteriskScraper,
	}
}

func (m *MetricsAsterisk) ActiveCalls() {
	m.activeCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Name:      "active_calls",
		Help:      "The number of active calls",
	})
}

func (m *MetricsAsterisk) TotalProcessedCalls() {
	m.totalCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Name:      "total_calls",
		Help:      "The number of total processed calls",
	})
}

func (m *MetricsAsterisk) GetTotalPeers() {
	m.totalPeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Name:      "total_extensions",
		Help:      "Total number of extensions regardless of their status",
	})
}
func (m *MetricsAsterisk) GetAvailablePeers() {
	m.availablePeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Name:      "available_extensions",
		Help:      "Total number of peers available extensions",
	})
}
func (m *MetricsAsterisk) GetUnavailablePeers() {
	m.unavailablePeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Name:      "unavailable_extensions",
		Help:      "Total number of peers unavailable extensions",
	})
}

func (m *MetricsAsterisk) RunAsteriskMetricsCollector() {
	for {
		select {
		case <-time.After(time.Second * 1):
			activeCalls, totalCalls := m.asteriskScraper.GetActiveAndTotalCalls()
			availablePeers, unavailablePeers, totalPeers := m.asteriskScraper.GetExtensions()

			m.activeCalls.Set(activeCalls)
			m.totalCalls.Set(totalCalls)
			m.totalPeers.Set(totalPeers)
			m.availablePeers.Set(availablePeers)
			m.unavailablePeers.Set(unavailablePeers)
		case <-m.ctx.Done():
			return
		}
	}
}
