package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ZeljkoBenovic/apgom/internal/config"
	"github.com/ZeljkoBenovic/apgom/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	conf    config.Config
	log     *slog.Logger
	metrics *metrics.Metrics
}

func NewApp(conf config.Config) App {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mtrsc := metrics.NewMetrics()
	return App{
		conf:    conf,
		log:     log,
		metrics: mtrsc,
	}
}

func (a App) Run() error {
	a.setupMetrics()

	a.log.Info("starting metrics listener", slog.String("port", a.conf.MetricsHttpListenPort))

	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%s", a.conf.MetricsHttpListenPort), nil)
}

func (a App) setupMetrics() {
	a.metrics.AsteriskMetrics.SipChannel()
}
