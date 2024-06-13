package metrics

import (
	"github.com/ZeljkoBenovic/apgom/internal/ami"
	"github.com/ZeljkoBenovic/apgom/internal/metrics/asterisk"
	scraperAsterisk "github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
)

type Metrics struct {
	asteriskMetrics *asterisk.MetricsAsterisk
}

func NewMetrics(ami *ami.Ami) *Metrics {
	asteriskMetrics := asterisk.NewMetricsAsterisk(scraperAsterisk.NewAsteriskScraper(ami))
	return &Metrics{
		asteriskMetrics: asteriskMetrics,
	}
}

func (m *Metrics) StartAsteriskMetrics() {
	m.asteriskMetrics.ActiveCalls()
	m.asteriskMetrics.TotalProcessedCalls()

	go m.asteriskMetrics.RunAsteriskMetricsCollector()
}
