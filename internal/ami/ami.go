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

// TODO: add support for PJSIP

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

func (a *Ami) GetRegistries() (float64, float64, float64) {
	var (
		errCh        = make(chan error)
		complete     = make(chan struct{})
		registries   []map[string]string
		registered   float64
		unRegistered float64
	)

	if err := a.cl.RegisterHandler("RegistryEntry", func(m map[string]string) {
		if m["State"] == "Registered" {
			registered++
		} else {
			unRegistered++
		}

		registries = append(registries, m)
	}); err != nil {
		errCh <- fmt.Errorf("could not set registry entry handler: %w", err)
	}

	if err := a.cl.RegisterHandler("RegistrationsComplete", func(m map[string]string) {
		complete <- struct{}{}
	}); err != nil {
		errCh <- fmt.Errorf("could not set registrations complete handler: %w", err)
	}

	defer func() {
		err := a.cl.UnregisterHandler("RegistrationsComplete", func(m map[string]string) {})
		err = a.cl.UnregisterHandler("RegistryEntry", func(m map[string]string) {})
		if err != nil {
			a.log.Error("unregister handler error", "err", err)
		}
	}()

	if _, err := a.cl.Action(map[string]string{
		"Action": "SIPshowregistry",
	}); err != nil {
		errCh <- fmt.Errorf("could not send sipshowregistry action: %w", err)
	}

	for {
		select {
		case err := <-errCh:
			a.log.Error("get registries error", "err", err)
			return -1, -1, -1
		case <-complete:
			return registered, unRegistered, float64(len(registries))
		}
	}
}

func (a *Ami) GetActiveAndTotalCalls() (float64, float64) {
	active, total, err := a.getActiveAndTotalCalls()
	if err != nil {
		a.log.Error("could not get active calls", "err", err)
		return -1, -1
	}

	return active, total
}

func (a *Ami) GetExtensions() (float64, float64, float64) {
	var (
		errCh               = make(chan error)
		doneCh              = make(chan struct{})
		peers               []map[string]string
		availableExtensions float64
		unavailabExtensions float64
	)
	if err := a.cl.RegisterHandler("PeerEntry", func(m map[string]string) {
		if m["Dynamic"] == "yes" {
			if strings.Contains(m["Status"], "OK") {
				availableExtensions++
			} else {
				unavailabExtensions++
			}

			peers = append(peers, m)
		}
	}); err != nil {
		errCh <- fmt.Errorf("could not register peerlist handler: %w", err)
	}

	if err := a.cl.RegisterHandler("PeerlistComplete", func(m map[string]string) {
		doneCh <- struct{}{}
	}); err != nil {
		errCh <- fmt.Errorf("could not register peerlist complete handler: %w", err)
	}

	defer func() {
		err := a.cl.UnregisterHandler("PeerEntry", func(m map[string]string) {})
		err = a.cl.UnregisterHandler("PeerlistComplete", func(m map[string]string) {})
		if err != nil {
			a.log.Error("unregister handler error", "err", err)
		}
	}()

	if _, err := a.cl.Action(map[string]string{
		"Action": "SIPPeers",
	}); err != nil {
		errCh <- fmt.Errorf("could run sippeers action: %w", err)
	}

	for {
		select {
		case <-doneCh:
			return availableExtensions, unavailabExtensions, float64(len(peers))
		case err := <-errCh:
			a.log.Error("get peers error", "err", err)
			return -1, -1, -1
		}
	}
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
