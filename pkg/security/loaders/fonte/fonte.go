package fonte

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/enrichman/portfolio-perfomance/pkg/security"
	"github.com/gocolly/colly/v2"
)

const (
	fonTeURLTemplate = "https://www.fondofonte.it/gestione-finanziaria/i-valori-quota-dei-comparti/comparto-%s/"
)

// QuoteLoader ...
type QuoteLoader struct {
	name    string
	isin    string
	urlName string
}

// New ...
func New(name, isin string) (*QuoteLoader, error) {
	isinURLName := strings.Split(isin, ".")

	if len(isinURLName) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for FonteQuoteLoader: \"%s\" - should be \"ISIN.URLName\"", isin)
	}

	return &QuoteLoader{
		name:    name,
		isin:    isinURLName[0],
		urlName: isinURLName[1],
	}, nil
}

// Name ...
func (f *QuoteLoader) Name() string {
	return f.name
}

// ISIN ...
func (f *QuoteLoader) ISIN() string {
	return f.isin
}

// URLName ...
func (f *QuoteLoader) URLName() string {
	return f.urlName
}

// LoadQuotes ...
func (f *QuoteLoader) LoadQuotes() ([]security.Quote, error) {
	c := colly.NewCollector()

	url := fmt.Sprintf(fonTeURLTemplate, f.urlName)

	type yearContent struct {
		year   string
		months []time.Month
		values []string
	}
	years := []yearContent{}

	c.OnHTML("article.content-text-page", func(e *colly.HTMLElement) {
		e.ForEach("h5.toggle-acf", func(i int, e *colly.HTMLElement) {
			years = append(years, yearContent{year: e.Text})
		})

		e.ForEach("div.toggle-content-acf", func(i int, e *colly.HTMLElement) {
			year := years[i]

			e.ForEach("div.toggle_element_row", func(i int, e *colly.HTMLElement) {
				monthString, valueString := parseRow(e.ChildTexts("span"))
				month, ok := convertMonth(monthString)
				if !ok {
					return
				}

				year.months = append(year.months, month)
				year.values = append(year.values, valueString)
			})

			reverse(year.months)
			reverse(year.values)

			years[i] = year
		})
	})

	c.Visit(url)

	reverse(years)

	quotes := []security.Quote{}

	for _, y := range years {
		for i, month := range y.months {
			dateString := fmt.Sprintf("%s %s", y.year, month)
			tt, err := time.Parse("2006 January", dateString)
			if err != nil {
				panic(err)
			}
			tt = tt.AddDate(0, 1, -1)

			y.values[i] = strings.ReplaceAll(y.values[i], ",", ".")
			closeQuote, err := strconv.ParseFloat(y.values[i], 32)
			if err != nil {
				panic(err)
			}

			quotes = append(quotes, security.Quote{
				Date:  tt,
				Close: float32(closeQuote),
			})
		}
	}

	return quotes, nil
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func convertMonth(month string) (time.Month, bool) {
	switch month {
	case "Gennaio":
		return time.January, true
	case "Febbraio":
		return time.February, true
	case "Marzo":
		return time.March, true
	case "Aprile":
		return time.April, true
	case "Maggio":
		return time.May, true
	case "Giugno":
		return time.June, true
	case "Luglio":
		return time.July, true
	case "Agosto":
		return time.August, true
	case "Settembre":
		return time.September, true
	case "Ottobre":
		return time.October, true
	case "Novembre":
		return time.November, true
	case "Dicembre":
		return time.December, true
	}
	return 0, false
}

func parseRow(values []string) (string, string) {
	return strings.TrimSpace(values[0]), strings.TrimSpace(values[1])
}
