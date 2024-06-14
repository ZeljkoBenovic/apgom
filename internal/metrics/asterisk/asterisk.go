package asterisk

import (
	"context"
	"time"

	"github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsAsterisk struct {
	ctx             context.Context
	asteriskScraper *asterisk.AsteriskScraper

	activeCalls prometheus.Gauge
	totalCalls  prometheus.Gauge

	totalPeers       prometheus.Gauge
	availablePeers   prometheus.Gauge
	unavailablePeers prometheus.Gauge

	totalRegistries        prometheus.Gauge
	registeredRegistries   prometheus.Gauge
	unregisteredRegistries prometheus.Gauge
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
		Subsystem: "calls",
		Name:      "active",
		Help:      "The number of active calls",
	})
}

func (m *MetricsAsterisk) TotalProcessedCalls() {
	m.totalCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "calls",
		Name:      "total",
		Help:      "The number of total processed calls",
	})
}

func (m *MetricsAsterisk) TotalExtensions() {
	m.totalPeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "extensions",
		Name:      "total",
		Help:      "Total number of extensions regardless of their status",
	})
}
func (m *MetricsAsterisk) AvailableExtensions() {
	m.availablePeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "extensions",
		Name:      "available",
		Help:      "Total number of peers available extensions",
	})
}

func (m *MetricsAsterisk) UnavailableExtensions() {
	m.unavailablePeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "extensions",
		Name:      "unavailable",
		Help:      "Total number of peers unavailable extensions",
	})
}

func (m *MetricsAsterisk) TotalRegistries() {
	m.totalRegistries = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "registries",
		Name:      "total",
		Help:      "Total number of registries",
	})
}

func (m *MetricsAsterisk) RegisteredRegistries() {
	m.registeredRegistries = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "registries",
		Name:      "active",
		Help:      "Total number of registered registries",
	})
}

func (m *MetricsAsterisk) UnRegisteredRegistries() {
	m.unregisteredRegistries = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "registries",
		Name:      "down",
		Help:      "Total number of unregistered registries",
	})
}

func (m *MetricsAsterisk) RunAsteriskMetricsCollector() {
	for {
		select {
		case <-time.After(time.Second * 1):
			activeCalls, totalCalls := m.asteriskScraper.GetActiveAndTotalCalls()
			availablePeers, unavailablePeers, totalPeers := m.asteriskScraper.GetExtensions()
			registeredRegistries, unRegisteredRegistries, totalRegistries := m.asteriskScraper.GetRegistries()

			m.activeCalls.Set(activeCalls)
			m.totalCalls.Set(totalCalls)
			m.totalPeers.Set(totalPeers)
			m.availablePeers.Set(availablePeers)
			m.unavailablePeers.Set(unavailablePeers)
			m.registeredRegistries.Set(registeredRegistries)
			m.unregisteredRegistries.Set(unRegisteredRegistries)
			m.totalRegistries.Set(totalRegistries)
		case <-m.ctx.Done():
			return
		}
	}
}
