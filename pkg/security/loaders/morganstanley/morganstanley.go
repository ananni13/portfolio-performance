package morganstanley

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
	"golang.org/x/exp/slices"
)

// QuoteLoader struct for MorganStanley.
type QuoteLoader struct {
	name         string
	isin         string
	fundID       string
	shareClassID string
}

// New creates a MorganStanley QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	isinFundIDShareID := strings.Split(isin, ".")

	if len(isinFundIDShareID) != 3 {
		return nil, fmt.Errorf("wrong ISIN format for MorganStanley QuoteLoader: \"%s\" - should be \"ISIN.fundID.shareClassID\"", isin)
	}

	return &QuoteLoader{
		name:         name,
		isin:         isinFundIDShareID[0],
		fundID:       isinFundIDShareID[1],
		shareClassID: isinFundIDShareID[2],
	}, nil
}

// Name returns the QuoteLoader name.
func (m *QuoteLoader) Name() string {
	return m.name
}

// ISIN returns the QuoteLoader isin.
func (m *QuoteLoader) ISIN() string {
	return m.isin
}

// FundID returns the QuoteLoader fundID.
func (m *QuoteLoader) FundID() string {
	return m.fundID
}

// ShareClassID returns the QuoteLoader shareClassID.
func (m *QuoteLoader) ShareClassID() string {
	return m.shareClassID
}

// LoadQuotes fetches quotes from MorganStanley.
func (m *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	historicalNav, err := fetchData(m.fundID)
	if err != nil {
		return nil, err
	}

	shareIdx := slices.IndexFunc(historicalNav.En.ShareClasses, func(s shareClass) bool {
		return s.ID == m.shareClassID
	})
	if shareIdx == -1 {
		return nil, nil
	}

	shareClass := historicalNav.En.ShareClasses[shareIdx]

	eurIdx := slices.IndexFunc(shareClass.Currencies, func(c currency) bool {
		return c.ID == "EUR"
	})
	if eurIdx == -1 {
		return nil, nil
	}

	currency := shareClass.Currencies[eurIdx]

	if len(currency.Series.Category) != len(currency.Series.Data) {
		log.Warn("Series Category and Data must be the same length")
		return nil, nil
	}

	quotesData := []quotes.Quote{}

	for idx, dateString := range currency.Series.Category {
		valueString := currency.Series.Data[idx]

		if dateString == "" || valueString == "" {
			continue
		}

		date, err := time.Parse("01/02/2006", dateString)
		if err != nil {
			log.Warnf("Error parsing date: %s", err)
			continue
		}

		closeQuote, err := strconv.ParseFloat(valueString, 32)
		if err != nil {
			log.Warnf("Error parsing quote: %s", err)
			continue
		}

		quotesData = append(quotesData, quotes.Quote{
			Date:  date,
			Close: float32(closeQuote),
		})
	}

	return quotesData, nil
}
