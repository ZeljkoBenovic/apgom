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

func (as *AsteriskScraper) GetActiveCalls() float64 {
	return as.ami.GetActiveCalls()
}

func (as *AsteriskScraper) GetTotalProcessedCalls() float64 {
	return as.ami.GetTotalProcessedCalls()
}
