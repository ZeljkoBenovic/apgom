package ami

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ZeljkoBenovic/apgom/internal/config"
	"github.com/ivahaev/amigo"
)

type Ami struct {
	conf config.Config
	log  *slog.Logger
	cl   *amigo.Amigo

	peerStatusCh chan map[string]string
}

type PeerStatus struct {
	Name      string
	IP        string
	Status    string
	LatencyMs float64
	Tech      string
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
		conf:         conf,
		log:          log,
		cl:           cl,
		peerStatusCh: make(chan map[string]string, 100),
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

func (a *Ami) GetPeerStatus() []PeerStatus {
	var resp []PeerStatus

	if a.peerStatusCh == nil {
		a.peerStatusCh = make(chan map[string]string)
	}

	for m := range a.peerStatusCh {
		if m == nil {
			break
		}

		var statusMs int
		var statusString string

		status := strings.Split(m["Status"], " ")
		statusString = status[0]

		if len(status) == 1 {
			statusMs = -1
		} else {
			statusMs, _ = strconv.Atoi(status[1][1:])
		}

		resp = append(resp, PeerStatus{
			Name:      m["ObjectName"],
			IP:        m["IPaddress"],
			Status:    statusString,
			LatencyMs: float64(statusMs),
			Tech:      m["Channeltype"],
		})
	}

	return resp
}

func (a *Ami) GetExtensions() (float64, float64, float64) {
	var (
		errCh               = make(chan error)
		doneCh              = make(chan struct{})
		peers               []map[string]string
		availableExtensions float64
		unavailabExtensions float64
		iaxPeersActionID    string
		sipPeersActionID    string
		closedChCounter     int
	)

	if err := a.cl.RegisterHandler("PeerEntry", func(m map[string]string) {
		a.peerStatusCh <- m

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
		// as we're catching two different peerlistcomplete events they can arrive at the same time
		// so we set up an artificial delay
		time.Sleep(time.Millisecond * 10)

		switch m["ActionID"] {
		case sipPeersActionID:
			closedChCounter++
		case iaxPeersActionID:
			closedChCounter++
		}

		if closedChCounter == 2 {
			a.peerStatusCh <- nil
			doneCh <- struct{}{}
		}
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

	sipPeersRes, err := a.cl.Action(map[string]string{
		"Action": "SIPPeers",
	})
	if err != nil {
		errCh <- fmt.Errorf("could run sippeers action: %w", err)
	}
	sipPeersActionID = sipPeersRes["ActionID"]

	iaxPeersRes, err := a.cl.Action(map[string]string{
		"Action": "IAXpeerlist",
	})
	if err != nil {
		errCh <- fmt.Errorf("could not send iaxpeer action: %w", err)
	}
	iaxPeersActionID = iaxPeersRes["ActionID"]

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

	// Asterisk 16.21.1
	if _, ok := resp["CommandResponse"]; !ok {
		split := strings.Split(resp["Output"], " ")
		totalInt, _ := strconv.Atoi(split[0])
		return -1, float64(totalInt), nil
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
