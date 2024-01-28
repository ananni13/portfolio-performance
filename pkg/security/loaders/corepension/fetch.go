package corepension

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"
)

const (
	corepensionURLTemplate = "https://www.corepension.it/ezjscore/call/ezjscamundibuzz::sfForwardFront::paramsList=service=ProxyProductSheetV3Front&routeId=_en-GB_874_%s_tab_3"
)

type parsedData struct {
	Date       time.Time
	CloseQuote float64
}

func fetchData(isin string) ([]parsedData, error) {
	c := colly.NewCollector()
	url := fmt.Sprintf(corepensionURLTemplate, isin)

	data := []parsedData{}

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

			data = append(data, parsedData{
				Date:       date,
				CloseQuote: closeQuote,
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error visiting/parsing page: %w", err)
	}

	return data, nil
}

func parseRowText(values []string) (string, string) {
	return values[0], values[1]
}
