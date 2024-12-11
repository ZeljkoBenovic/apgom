package asterisk

import "github.com/ZeljkoBenovic/apgom/internal/ami"

type AsteriskScraper struct {
	ami *ami.Ami
}

func NewAsteriskScraper(ami *ami.Ami) *AsteriskScraper {
	return &AsteriskScraper{
		ami: ami,
	}
}

func (as *AsteriskScraper) GetActiveAndTotalCalls() (activeCalls, totalCalls float64) {
	return as.ami.GetActiveAndTotalCalls()
}

func (as *AsteriskScraper) GetExtensions() (availableExtensions, unavailableExtensions, totalExtensions float64) {
	return as.ami.GetExtensions()
}

func (as *AsteriskScraper) GetRegistries() (registered, unRegistered, totalRegistries float64) {
	return as.ami.GetRegistries()
}

func (as *AsteriskScraper) GetCallsByTrunkAndDirection() (map[string]ami.CallsDirectionRegistry, error) {
	return as.ami.GetInboundAndOutboundCallsPerTrunk()
}
