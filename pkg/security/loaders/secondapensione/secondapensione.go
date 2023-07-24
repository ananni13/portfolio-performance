package secondapensione

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security"
	"github.com/gocolly/colly/v2"
)

const (
	secondaPensioneURLTemplate = "https://www.secondapensione.it/ezjscore/call/ezjscamundibuzz::sfForwardFront::paramsList=service=ProxyProductSheetV3Front&routeId=_en-GB_879_%s_tab_3"
)

// QuoteLoader struct for SecondaPensione.
type QuoteLoader struct {
	name string
	isin string
}

// New creates a SecondaPensione QuoteLoader.
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

// LoadQuotes fetches quotes from SecondaPensione.
func (s *QuoteLoader) LoadQuotes() ([]security.Quote, error) {
	c := colly.NewCollector()

	url := fmt.Sprintf(secondaPensioneURLTemplate, s.isin)

	quotes := []security.Quote{}

	c.OnHTML("#tableVl", func(e *colly.HTMLElement) {
		e.ForEach("tbody tr", func(i int, e *colly.HTMLElement) {
			dateString, valueString := parseRowText(e.ChildTexts("td"))

			if dateString == "" || valueString == "" {
				return
			}

			date, err := time.Parse("02/01/2006", dateString)
			if err != nil {
				log.Warnf("Error parsing date: %s", err)
				return
			}

			closeQuote, err := strconv.ParseFloat(valueString, 32)
			if err != nil {
				log.Warnf("Error parsing quote: %s", err)
				return
			}

			quotes = append(quotes, security.Quote{
				Date:  date,
				Close: float32(closeQuote),
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error visiting/parsing page")
	}

	return quotes, nil
}

func parseRowText(values []string) (string, string) {
	return values[0], values[1]
}
