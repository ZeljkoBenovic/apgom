package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/ZeljkoBenovic/apgom/internal/ami"
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
	am, err := ami.NewAmi(conf, log)
	if err != nil {
		log.Error("could not create new ami instance", "err", err)
		os.Exit(1)
	}

	go handleOsSignals(log)

	mtrsc := metrics.NewMetrics(am)
	return App{
		conf:    conf,
		log:     log,
		metrics: mtrsc,
	}
}

func (a App) Run() error {
	a.metrics.StartAsteriskMetrics()

	a.log.Info("starting metrics listener", slog.String("port", a.conf.MetricsHttpListenPort))

	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%s", a.conf.MetricsHttpListenPort), nil)
}

func handleOsSignals(log *slog.Logger) {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	<-sig
	log.Info("shutting down metrics server")
	os.Exit(0)
}
