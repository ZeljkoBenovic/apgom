package metrics

import (
	"context"
	"log/slog"

	"github.com/ZeljkoBenovic/apgom/internal/ami"
	"github.com/ZeljkoBenovic/apgom/internal/metrics/asterisk"
	scraperAsterisk "github.com/ZeljkoBenovic/apgom/internal/scrapers/asterisk"
)

type Metrics struct {
	ctx             context.Context
	log             *slog.Logger
	asteriskMetrics *asterisk.MetricsAsterisk
}

func NewMetrics(ctx context.Context, log *slog.Logger, ami *ami.Ami) (*Metrics, error) {
	asteriskMetrics, err := asterisk.NewMetricsAsterisk(ctx, log.With("module", "metrics.asterisk"), scraperAsterisk.NewAsteriskScraper(ami))
	if err != nil {
		return nil, err
	}

	return &Metrics{
		asteriskMetrics: asteriskMetrics,
		log:             log,
	}, nil
}

func (m *Metrics) StartAsteriskMetrics() error {
	var err error

	err = m.asteriskMetrics.ActiveCalls()
	err = m.asteriskMetrics.TotalProcessedCalls()
	err = m.asteriskMetrics.TotalExtensions()
	err = m.asteriskMetrics.AvailableExtensions()
	err = m.asteriskMetrics.UnavailableExtensions()
	err = m.asteriskMetrics.RegisteredRegistries()
	err = m.asteriskMetrics.UnRegisteredRegistries()
	err = m.asteriskMetrics.TotalRegistries()

	if err != nil {
		return err
	}

	go m.asteriskMetrics.RunAsteriskMetricsCollector()

	return nil
}
