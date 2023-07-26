package fonte

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
)

// QuoteLoader struct for FonTe.
type QuoteLoader struct {
	name    string
	isin    string
	urlName string
}

// New creates a FonTe QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	isinURLName := strings.Split(isin, ".")

	if len(isinURLName) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for FonTe QuoteLoader: \"%s\" - should be \"ISIN.URLName\"", isin)
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

// LoadQuotes fetches quotes from FonTe.
func (f *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	years, err := fetchData(f.urlName)
	if err != nil {
		return nil, err
	}

	quotesData := []quotes.Quote{}

	for _, y := range years {
		for i, month := range y.months {
			dateString := fmt.Sprintf("%s %s", y.year, month)
			tt, err := time.Parse("2006 January", dateString)
			if err != nil {
				log.Warnf("Error parsing date: %s", err)
				continue
			}
			tt = tt.AddDate(0, 1, -1)

			y.values[i] = strings.ReplaceAll(y.values[i], ",", ".")
			closeQuote, err := strconv.ParseFloat(y.values[i], 32)
			if err != nil {
				log.Warnf("Error parsing quote: %s", err)
				continue
			}

			quotesData = append(quotesData, quotes.Quote{
				Date:  tt,
				Close: float32(closeQuote),
			})
		}
	}

	return quotesData, nil
}
