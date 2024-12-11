package asterisk

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/ZeljkoBenovic/apgom/internal/config"
	"github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsAsterisk struct {
	ctx             context.Context
	log             *slog.Logger
	conf            config.Config
	asteriskScraper *asterisk.AsteriskScraper

	hostName string
	hostIP   string

	activeCalls *prometheus.GaugeVec
	totalCalls  prometheus.Gauge

	totalPeers       prometheus.Gauge
	availablePeers   prometheus.Gauge
	unavailablePeers prometheus.Gauge

	totalRegistries        prometheus.Gauge
	registeredRegistries   prometheus.Gauge
	unregisteredRegistries prometheus.Gauge
}

var (
	commonLabels = []string{"hostname", "host_ips"}
)

func NewMetricsAsterisk(ctx context.Context, log *slog.Logger, conf config.Config, asteriskScraper *asterisk.AsteriskScraper) (*MetricsAsterisk, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get hostname: %w", err)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("could not get network interfaces: %w", err)
	}

	var addresses string
	for _, nic := range ifaces {
		addrs, err := nic.Addrs()
		if err != nil {
			return nil, fmt.Errorf("coult not get interface addresses: %w", err)
		}

		for _, a := range addrs {
			// skip localhost, ipv6 and broadcast address
			ad := strings.Split(a.String(), ".")
			if ad[0] == "127" || ad[0] == "169" || strings.Contains(ad[0], "::") {
				continue
			}

			// strip network mask
			msk := strings.Split(a.String(), "/")

			addresses += msk[0] + ","
		}
	}

	return &MetricsAsterisk{
		ctx:             ctx,
		log:             log,
		conf:            conf,
		asteriskScraper: asteriskScraper,
		hostName:        hostname,
		// trim the comma at the end of the string
		hostIP: strings.TrimRight(addresses, ","),
	}, nil
}

func (m *MetricsAsterisk) ActiveCalls() error {
	m.activeCalls = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "calls",
		Name:      "active",
		Help:      "The number of active calls",
	}, commonLabels)

	return nil
}

func (m *MetricsAsterisk) TotalProcessedCalls() error {
	tc := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "calls",
		Name:      "total",
		Help:      "The number of total processed calls",
	}, commonLabels)

	gauge, err := tc.GetMetricWithLabelValues(m.hostName, m.hostIP)
	if err != nil {
		return err
	}

	m.totalCalls = gauge

	return nil
}

func (m *MetricsAsterisk) TotalExtensions() error {
	tp := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "extensions",
		Name:      "total",
		Help:      "Total number of extensions regardless of their status",
	}, commonLabels)

	gauge, err := tp.GetMetricWithLabelValues(m.hostName, m.hostIP)
	if err != nil {
		return err
	}

	m.totalPeers = gauge

	return nil
}
func (m *MetricsAsterisk) AvailableExtensions() error {
	ap := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "extensions",
		Name:      "available",
		Help:      "Total number of peers available extensions",
	}, commonLabels)

	gauge, err := ap.GetMetricWithLabelValues(m.hostName, m.hostIP)
	if err != nil {
		return err
	}

	m.availablePeers = gauge

	return nil
}

func (m *MetricsAsterisk) UnavailableExtensions() error {
	up := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "extensions",
		Name:      "unavailable",
		Help:      "Total number of peers unavailable extensions",
	}, commonLabels)

	gauge, err := up.GetMetricWithLabelValues(m.hostName, m.hostIP)
	if err != nil {
		return err
	}

	m.unavailablePeers = gauge

	return nil
}

func (m *MetricsAsterisk) TotalRegistries() error {
	tr := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "registries",
		Name:      "total",
		Help:      "Total number of registries",
	}, commonLabels)

	gauge, err := tr.GetMetricWithLabelValues(m.hostName, m.hostIP)
	if err != nil {
		return err
	}

	m.totalRegistries = gauge

	return nil
}

func (m *MetricsAsterisk) RegisteredRegistries() error {
	rr := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "registries",
		Name:      "active",
		Help:      "Total number of registered registries",
	}, commonLabels)

	gauge, err := rr.GetMetricWithLabelValues(m.hostName, m.hostIP)
	if err != nil {
		return err
	}

	m.registeredRegistries = gauge

	return nil
}

func (m *MetricsAsterisk) UnRegisteredRegistries() error {
	ur := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "asterisk",
		Subsystem: "registries",
		Name:      "down",
		Help:      "Total number of unregistered registries",
	}, commonLabels)

	gauge, err := ur.GetMetricWithLabelValues(m.hostName, m.hostIP)
	if err != nil {
		return err
	}

	m.unregisteredRegistries = gauge

	return nil
}

func (m *MetricsAsterisk) RunAsteriskMetricsCollector() {
	for {
		select {
		case <-time.After(time.Second * time.Duration(m.conf.ScrapeTimeSec)):
			// TODO: set debug log for when scraping stops
			activeCalls, totalCalls := m.asteriskScraper.GetActiveAndTotalCalls()
			availablePeers, unavailablePeers, totalPeers := m.asteriskScraper.GetExtensions()
			registeredRegistries, unRegisteredRegistries, totalRegistries := m.asteriskScraper.GetRegistries()

			ac, err := m.activeCalls.GetMetricWithLabelValues(m.hostName, m.hostIP)
			if err != nil {
				m.log.Error("error creating active calls gauge", "err", err)
			} else {
				ac.Set(activeCalls)
			}

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
