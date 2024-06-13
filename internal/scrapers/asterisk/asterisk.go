package asterisk

type AsteriskScraper struct {
}

func NewAsteriskScraper() *AsteriskScraper {
	return &AsteriskScraper{}
}

func (as *AsteriskScraper) GetSIPChannels() float64 {
	// TODO: implement actual scraping logic
	return 2
}
