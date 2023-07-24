package financialtimes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security"
	"golang.org/x/exp/slices"
)

const (
	financialTimesURL = "https://markets.ft.com/data/chartapi/series"
)

// QuoteLoader ...
type QuoteLoader struct {
	name   string
	isin   string
	symbol string
}

// New ...
func New(name, isin string) (*QuoteLoader, error) {
	isinLabelSymbol := strings.Split(isin, ".")

	if len(isinLabelSymbol) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for FinancialTimesQuoteLoader: \"%s\" - should be \"ISIN.symbol\"", isin)
	}

	return &QuoteLoader{
		name:   name,
		isin:   isinLabelSymbol[0],
		symbol: isinLabelSymbol[1],
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

// Symbol ...
func (f *QuoteLoader) Symbol() string {
	return f.symbol
}

type requestPayload struct {
	Days           int              `json:"days"`
	DataPeriod     string           `json:"dataPeriod"`
	DataInterval   int              `json:"dataInterval"`
	YFormat        string           `json:"yFormat"`
	ReturnDateType string           `json:"returnDateType"`
	Elements       []requestElement `json:"elements"`
}

type requestElement struct {
	Type   string `json:"Type"`
	Symbol string `json:"Symbol"`
}

type responsePayload struct {
	Dates    []string          `json:"Dates"`
	Elements []responseElement `json:"Elements"`
}

type responseElement struct {
	ComponentSeries []componentSeries `json:"ComponentSeries"`
}

type componentSeries struct {
	Type   string    `json:"Type"`
	Values []float32 `json:"Values"`
}

// LoadQuotes ...
func (f *QuoteLoader) LoadQuotes() ([]security.Quote, error) {
	payload := requestPayload{
		Days:           365 * 30,
		DataPeriod:     "Day",
		DataInterval:   1,
		YFormat:        "0.###",
		ReturnDateType: "ISO8601",
		Elements: []requestElement{
			{
				Type:   "price",
				Symbol: f.symbol,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body")
	}

	res, err := http.Post(
		financialTimesURL,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("error during post request: %w", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var result responsePayload
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling body: %w", err)
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

	quotes := []security.Quote{}

	for idx, dateString := range dates {
		value := values[idx]

		date, err := time.Parse("2006-01-02T15:04:05", dateString)
		if err != nil {
			log.Warnf("Error parsing date: %s", err)
			continue
		}

		quotes = append(quotes, security.Quote{
			Date:  date,
			Close: value,
		})
	}

	return quotes, nil
}
