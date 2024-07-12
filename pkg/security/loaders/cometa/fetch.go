package cometa

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"
	"github.com/goodsign/monday"
)

const (
	cometaURLTemplate = "https://www.cometafondo.it/quota/%s"
)

type parsedData struct {
	Date       time.Time
	CloseQuote float64
}

func fetchData(isin string) ([]parsedData, error) {
	c := colly.NewCollector()
	url := fmt.Sprintf(cometaURLTemplate, isin)

	data := []parsedData{}

	c.OnHTML("div.content", func(e *colly.HTMLElement) {
		e.ForEach("table.table-values", func(i int, e *colly.HTMLElement) {
			e.ForEach("tbody tr", func(i int, e *colly.HTMLElement) {
				if i == 0 {
					return
				}

				yearString, monthString, valueString := parseRowText(e.ChildTexts("td"))

				date, err := monday.ParseInLocation("Jan 2006", monthString+" "+yearString, time.UTC, monday.LocaleItIT)
				if err != nil {
					log.Warnf("Error parsing date: %s", err)
					return
				}

				valueString = strings.Split(valueString, " ")[0]
				closeQuote, err := strconv.ParseFloat(strings.ReplaceAll(valueString, ",", "."), 32)
				if err != nil {
					log.Warnf("Error parsing quote: %s", err)
					return
				}

				data = append(data, parsedData{
					Date:       date.AddDate(0, 1, -1),
					CloseQuote: closeQuote,
				})
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error visiting/parsing page: %w", err)
	}

	return data, nil
}

func parseRowText(values []string) (string, string, string) {
	return values[0], values[1], values[2]
}
