package metrics

import (
	"github.com/ZeljkoBenovic/apgom/internal/metrics/asterisk"
	scraperAsterisk "github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
)

type Metrics struct {
	AsteriskMetrics *asterisk.MetricsAsterisk
}

func NewMetrics() *Metrics {
	asteriskMetrics := asterisk.NewMetricsAsterisk(scraperAsterisk.NewAsteriskScraper())
	return &Metrics{
		AsteriskMetrics: asteriskMetrics,
	}
}
