package config

import "flag"

type Config struct {
	MetricsHttpListenPort string
}

func NewConfig() Config {
	c := Config{}

	flag.StringVar(&c.MetricsHttpListenPort, "port", "3000", "metrics http listen port")
	flag.Parse()

	return c
}
