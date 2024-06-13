package asterisk

import (
	"github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsAsterisk struct {
	sipChannels     prometheus.Gauge
	asteriskScraper *asterisk.AsteriskScraper
}

func NewMetricsAsterisk(asteriskScraper *asterisk.AsteriskScraper) *MetricsAsterisk {
	return &MetricsAsterisk{
		asteriskScraper: asteriskScraper,
	}
}

func (m *MetricsAsterisk) SipChannel() {
	m.sipChannels = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "asterisk",
		Name:        "sip_channels",
		Help:        "The number of Asterisk SIP channels",
		ConstLabels: nil,
	})

	m.getSipChannels()
}

func (m *MetricsAsterisk) getSipChannels() {
	// TODO: implement asterisk scraping module
	m.sipChannels.Set(m.asteriskScraper.GetSIPChannels())
}
