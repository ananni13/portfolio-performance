package cometa

import (
	"fmt"
	"strings"

	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
)

// QuoteLoader struct for Cometa.
type QuoteLoader struct {
	name    string
	isin    string
	urlName string
}

// New creates a Cometa QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	isinURLName := strings.Split(isin, ".")

	if len(isinURLName) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for Cometa QuoteLoader: \"%s\" - should be \"ISIN.URLName\"", isin)
	}

	return &QuoteLoader{
		name:    name,
		isin:    isinURLName[0],
		urlName: isinURLName[1],
	}, nil
}

// Name returns the QuoteLoader name.
func (f *QuoteLoader) Name() string {
	return f.name
}

// ISIN returns the QuoteLoader isin.
func (f *QuoteLoader) ISIN() string {
	return f.isin
}

// URLName returns the QuoteLoader urlName.
func (f *QuoteLoader) URLName() string {
	return f.urlName
}

// LoadQuotes fetches quotes from Cometa.
func (s *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	data, err := fetchData(s.urlName)
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
