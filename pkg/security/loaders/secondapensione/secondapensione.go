package secondapensione

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-perfomance/pkg/security"
	"github.com/gocolly/colly/v2"
)

const (
	SecondaPensioneUrlTemplate = "https://www.secondapensione.it/ezjscore/call/ezjscamundibuzz::sfForwardFront::paramsList=service=ProxyProductSheetV3Front&routeId=_en-GB_879_%s_tab_3"
)

type SecondaPensione struct {
	name string
	isin string
}

func New(name, isin string) (*SecondaPensione, error) {
	return &SecondaPensione{
		name: name,
		isin: isin,
	}, nil
}

func (e *SecondaPensione) Name() string {
	return e.name
}

func (e *SecondaPensione) ISIN() string {
	return e.isin
}

func (s *SecondaPensione) LoadQuotes() ([]security.Quote, error) {
	c := colly.NewCollector()

	url := fmt.Sprintf(SecondaPensioneUrlTemplate, s.isin)

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

	c.Visit(url)

	return quotes, nil
}

func parseRowText(values []string) (string, string) {
	return values[0], values[1]
}
