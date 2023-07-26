package financialtimes

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
	"golang.org/x/exp/slices"
)

// QuoteLoader struct for FinancialTimes.
type QuoteLoader struct {
	name   string
	isin   string
	symbol string
}

// New creates a FinancialTimes QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	isinLabelSymbol := strings.Split(isin, ".")

	if len(isinLabelSymbol) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for FinancialTimes QuoteLoader: \"%s\" - should be \"ISIN.symbol\"", isin)
	}

	return &QuoteLoader{
		name:   name,
		isin:   isinLabelSymbol[0],
		symbol: isinLabelSymbol[1],
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

// Symbol returns the QuoteLoader symbol.
func (f *QuoteLoader) Symbol() string {
	return f.symbol
}

// LoadQuotes fetches quotes from FinancialTimes.
func (f *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	result, err := fetchData(f.symbol)
	if err != nil {
		return nil, err
	}

	if len(result.Elements) == 0 {
		return nil, nil
	}

	series := result.Elements[0].ComponentSeries
	closeIdx := slices.IndexFunc(series, func(c componentSeries) bool {
		return c.Type == "Close"
	})
	if closeIdx == -1 {
		return nil, nil
	}

	component := series[closeIdx]

	dates := result.Dates
	values := component.Values

	if len(dates) != len(values) {
		log.Warn("Dates and Values must be the same length")
		return nil, nil
	}

	quotesData := []quotes.Quote{}

	for idx, dateString := range dates {
		value := values[idx]

		date, err := time.Parse("2006-01-02T15:04:05", dateString)
		if err != nil {
			log.Warnf("Error parsing date: %s", err)
			continue
		}

		quotesData = append(quotesData, quotes.Quote{
			Date:  date,
			Close: value,
		})
	}

	return quotesData, nil
}
