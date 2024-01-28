package corepension

import (
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

	for _, quote := range data {
		quotesData = append(quotesData, quotes.Quote{
			Date:  quote.Date,
			Close: float32(quote.CloseQuote),
		})
	}

	return quotesData, nil
}
