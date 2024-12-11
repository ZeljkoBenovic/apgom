package config

import (
	"flag"
	"log"
)

type Config struct {
	MetricsHttpListenPort string
	MetricsHttpListenHost string
	MetricsHttpListenPath string
	AsteriskAMIHost       string
	AsteriskAMIUser       string
	AsteriskAMIPass       string
	ScrapeTimeSec         int
}

func NewConfig() Config {
	c := Config{}

	flag.StringVar(&c.MetricsHttpListenPort, "port", "3000", "metrics http listen port")
	flag.StringVar(&c.MetricsHttpListenHost, "host", "0.0.0.0", "metrics http listen host/ip")
	flag.StringVar(&c.MetricsHttpListenPath, "path", "/metrics", "metrics http path")
	flag.StringVar(&c.AsteriskAMIHost, "ami-host", "localhost", "asterisk ami hostname or ip")
	flag.StringVar(&c.AsteriskAMIUser, "ami-user", "admin", "asterisk ami username")
	flag.StringVar(&c.AsteriskAMIPass, "ami-pass", "", "asterisk ami password")
	flag.IntVar(&c.ScrapeTimeSec, "scrape-time", 1, "metrics scrape time in seconds which is usually aligned with Prometheus scraper time")
	flag.Parse()

	if c.AsteriskAMIPass == "" {
		log.Fatalln("asterisk ami password not defined")
	}

	return c
}
