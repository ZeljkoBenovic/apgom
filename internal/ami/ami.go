package ami

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/ZeljkoBenovic/apgom/internal/config"
	"github.com/ivahaev/amigo"
)

type Ami struct {
	conf config.Config
	log  *slog.Logger
	cl   *amigo.Amigo
}

func NewAmi(conf config.Config, log *slog.Logger) (*Ami, error) {
	cl := amigo.New(&amigo.Settings{
		Username: conf.AsteriskAMIUser,
		Password: conf.AsteriskAMIPass,
		Host:     conf.AsteriskAMIHost,
	})
	connected := make(chan bool)

	cl.Connect()

	cl.On("connect", func(s string) {
		log.Info("connected to ami interface", "msg", s)
		connected <- true
	})

	cl.On("error", func(s string) {
		log.Error("could not connect to ami interface", "msg", s)
		connected <- false
	})

	isConnected := <-connected
	if !isConnected {
		return nil, fmt.Errorf("ami connection failed")
	}

	return &Ami{
		conf: conf,
		log:  log,
		cl:   cl,
	}, nil
}

func (a *Ami) GetActiveCalls() float64 {
	active, _, err := a.getActiveAndTotalCalls()
	if err != nil {
		a.log.Error("could not get active calls", "err", err)
		return -1
	}

	return active
}

func (a *Ami) GetTotalProcessedCalls() float64 {
	_, processed, err := a.getActiveAndTotalCalls()
	if err != nil {
		a.log.Error("could not get active calls", "err", err)
		return -1
	}

	return processed
}

func (a *Ami) getActiveAndTotalCalls() (float64, float64, error) {
	resp, err := a.cl.Action(map[string]string{
		"Action":  "Command",
		"Command": "core show channels",
	})
	if err != nil {
		a.log.Error("could not run ami command", "err", err)
	}

	cmdRespArr := strings.Split(resp["CommandResponse"], "\n")
	activeCallsArr := strings.Split(cmdRespArr[len(cmdRespArr)-2], " ")
	totalCallsArr := strings.Split(cmdRespArr[len(cmdRespArr)-1], " ")

	activeCalls, err := strconv.ParseFloat(activeCallsArr[0], 64)
	if err != nil {
		return -1, -1, err
	}

	totalCalls, err := strconv.ParseFloat(totalCallsArr[0], 64)
	if err != nil {
		return -1, -1, err
	}

	return activeCalls, totalCalls, nil
}
