package corepension

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
)

// QuoteLoader struct for CorePension.
type QuoteLoader struct {
	name string
	isin string
}

// New creates a CorePension QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	return &QuoteLoader{
		name: name,
		isin: isin,
	}, nil
}

// Name returns the QuoteLoader name.
func (s *QuoteLoader) Name() string {
	return s.name
}

// ISIN returns the QuoteLoader isin.
func (s *QuoteLoader) ISIN() string {
	return s.isin
}

// LoadQuotes fetches quotes from CorePension.
func (s *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	data, err := fetchData(s.isin)
	if err != nil {
		return nil, err
	}

	quotesData := []quotes.Quote{}

	for _, quote := range data[0].NavHistory {
		date, err := time.Parse("2006-01-02", quote.Date)
		if err != nil {
			log.Errorf("failed to parse date: %v", err)
			continue
		}

		quotesData = append(quotesData, quotes.Quote{
			Date:  date,
			Close: float32(quote.Value),
		})
	}

	return quotesData, nil
}
