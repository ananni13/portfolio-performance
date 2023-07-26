package fondidoc

import (
	"fmt"
	"strings"
	"time"

	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
)

// QuoteLoader struct for FondiDoc.
type QuoteLoader struct {
	name   string
	isin   string
	fundID string
}

// New creates a FondiDoc QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	isinID := strings.Split(isin, ".")

	if len(isinID) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for FondiDoc QuoteLoader: \"%s\" - should be \"ISIN.fundID\"", isin)
	}

	return &QuoteLoader{
		name:   name,
		isin:   isinID[0],
		fundID: isinID[1],
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

// FundID returns the QuoteLoader fundID.
func (f *QuoteLoader) FundID() string {
	return f.fundID
}

// LoadQuotes fetches quotes from FondiDoc.
func (f *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	fondiDoc, err := fetchData(f.fundID)
	if err != nil {
		return nil, err
	}

	fundData := fondiDoc[f.fundID].(map[string]interface{})

	quotesData := []quotes.Quote{}
	for _, quote := range fundData["data"].([]interface{}) {
		quotesData = append(quotesData, quotes.Quote{
			Date:  time.Unix(int64(quote.([]interface{})[0].(float64))*100, 0),
			Close: float32(quote.([]interface{})[1].(float64)),
		})
	}

	return quotesData, nil
}
