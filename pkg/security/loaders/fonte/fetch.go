package fonte

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	fonTeURLTemplate = "https://www.fondofonte.it/gestione-finanziaria/i-valori-quota-dei-comparti/comparto-%s/"
)

type yearContent struct {
	year   string
	months []time.Month
	values []string
}

func fetchData(urlName string) ([]yearContent, error) {
	c := colly.NewCollector()
	url := fmt.Sprintf(fonTeURLTemplate, urlName)
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

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error visiting/parsing page: %w", err)
	}

	reverse(years)

	return years, nil
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
