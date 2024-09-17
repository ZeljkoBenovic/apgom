package app

import (
	"context"
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
	ctx     context.Context
	conf    config.Config
	log     *slog.Logger
	metrics *metrics.Metrics
}

func NewApp(conf config.Config) App {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	am, err := ami.NewAmi(conf, log)
	if err != nil {
		log.Error("could not create new ami instance", "err", err)
		os.Exit(1)
	}

	go handleOsSignals(log)

	mtrsc, err := metrics.NewMetrics(ctx, log, am)
	if err != nil {
		log.Error("could not setup new metrics", "err", err)
		os.Exit(1)
	}

	return App{
		ctx:     ctx,
		conf:    conf,
		log:     log,
		metrics: mtrsc,
	}
}

func (a App) Run() error {
	if err := a.metrics.StartAsteriskMetrics(); err != nil {
		return fmt.Errorf("could not start asterisk metrics: %v", err)
	}

	a.log.Info("starting metrics listener",
		slog.String(
			"on",
			fmt.Sprintf(
				"%s:%s%s",
				a.conf.MetricsHttpListenHost,
				a.conf.MetricsHttpListenPort,
				a.conf.MetricsHttpListenPath,
			),
		),
	)

	http.Handle(a.conf.MetricsHttpListenPath, promhttp.Handler())

	return http.ListenAndServe(
		fmt.Sprintf("%s:%s", a.conf.MetricsHttpListenHost, a.conf.MetricsHttpListenPort),
		nil,
	)
}

func handleOsSignals(log *slog.Logger) {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	<-sig
	log.Info("shutting down metrics server")
	os.Exit(0)
}
