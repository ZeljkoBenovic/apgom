package metrics

import (
	"context"

	"github.com/ZeljkoBenovic/apgom/internal/ami"
	"github.com/ZeljkoBenovic/apgom/internal/metrics/asterisk"
	scraperAsterisk "github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
)

type Metrics struct {
	ctx             context.Context
	asteriskMetrics *asterisk.MetricsAsterisk
}

func NewMetrics(ctx context.Context, ami *ami.Ami) *Metrics {
	asteriskMetrics := asterisk.NewMetricsAsterisk(ctx, scraperAsterisk.NewAsteriskScraper(ami))
	return &Metrics{
		asteriskMetrics: asteriskMetrics,
	}
}

func (m *Metrics) StartAsteriskMetrics() {
	m.asteriskMetrics.ActiveCalls()
	m.asteriskMetrics.TotalProcessedCalls()
	m.asteriskMetrics.GetTotalPeers()
	m.asteriskMetrics.GetAvailablePeers()
	m.asteriskMetrics.GetUnavailablePeers()

	go m.asteriskMetrics.RunAsteriskMetricsCollector()
}
