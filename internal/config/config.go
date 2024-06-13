package config

import (
	"flag"
	"log"
)

type Config struct {
	MetricsHttpListenPort string
	AsteriskAMIHost       string
	AsteriskAMIUser       string
	AsteriskAMIPass       string
}

func NewConfig() Config {
	c := Config{}

	flag.StringVar(&c.MetricsHttpListenPort, "port", "3000", "metrics http listen port")
	flag.StringVar(&c.AsteriskAMIHost, "ami-host", "localhost", "asterisk ami hostname or ip")
	flag.StringVar(&c.AsteriskAMIUser, "ami-user", "admin", "asterisk ami username")
	flag.StringVar(&c.AsteriskAMIPass, "ami-pass", "", "asterisk ami password")
	flag.Parse()

	if c.AsteriskAMIPass == "" {
		log.Fatalln("asterisk ami password not defined")
	}

	return c
}
