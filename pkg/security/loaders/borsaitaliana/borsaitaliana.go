package borsaitaliana

import (
	"fmt"
	"strings"
	"time"

	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
)

// QuoteLoader struct for BorsaItaliana.
type QuoteLoader struct {
	name   string
	isin   string
	market string
}

// New creates a BorsaItaliana QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	isinMarket := strings.Split(isin, ".")

	if len(isinMarket) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for BorsaItaliana QuoteLoader: \"%s\" - should be \"ISIN.market\"", isin)
	}

	return &QuoteLoader{
		name:   name,
		isin:   isinMarket[0],
		market: isinMarket[1],
	}, nil
}

// Name returns the QuoteLoader name.
func (b *QuoteLoader) Name() string {
	return b.name
}

// ISIN returns the QuoteLoader isin.
func (b *QuoteLoader) ISIN() string {
	return b.isin
}

// LoadQuotes fetches quotes from BorsaItaliana.
func (b *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	result, err := fetchData(b.isin, b.market)
	if err != nil {
		return nil, err
	}

	quotesData := []quotes.Quote{}
	for _, quote := range result.Data {
		quotesData = append(quotesData, quotes.Quote{
			Date:  time.Unix(int64(quote[0]/1000), 0).In(time.UTC),
			Close: quote[1],
		})
	}

	return quotesData, nil
}
