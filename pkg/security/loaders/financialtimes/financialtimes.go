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
	"github.com/enrichman/portfolio-perfomance/pkg/security"
	"golang.org/x/exp/slices"
)

const (
	FinancialTimesUrl = "https://markets.ft.com/data/chartapi/series"
)

type FinancialTimesQuoteLoader struct {
	name   string
	isin   string
	symbol string
}

func New(name, isin string) (*FinancialTimesQuoteLoader, error) {
	isinLabelSymbol := strings.Split(isin, ".")

	if len(isinLabelSymbol) != 2 {
		return nil, fmt.Errorf("Wrong ISIN format for FinancialTimesQuoteLoader: \"%s\" - should be \"ISIN.symbol\"", isin)
	}

	return &FinancialTimesQuoteLoader{
		name:   name,
		isin:   isinLabelSymbol[0],
		symbol: isinLabelSymbol[1],
	}, nil
}

func (f *FinancialTimesQuoteLoader) Name() string {
	return f.name
}

func (f *FinancialTimesQuoteLoader) ISIN() string {
	return f.isin
}

func (f *FinancialTimesQuoteLoader) Symbol() string {
	return f.symbol
}

type RequestPayload struct {
	Days           int              `json:"days"`
	DataPeriod     string           `json:"dataPeriod"`
	DataInterval   int              `json:"dataInterval"`
	YFormat        string           `json:"yFormat"`
	ReturnDateType string           `json:"returnDateType"`
	Elements       []RequestElement `json:"elements"`
}

type RequestElement struct {
	Type   string `json:"Type"`
	Symbol string `json:"Symbol"`
}

type ResponsePayload struct {
	Dates    []string          `json:"Dates"`
	Elements []ResponseElement `json:"Elements"`
}

type ResponseElement struct {
	ComponentSeries []ComponentSeries `json:"ComponentSeries"`
}

type ComponentSeries struct {
	Type   string    `json:"Type"`
	Values []float32 `json:"Values"`
}

func (f *FinancialTimesQuoteLoader) LoadQuotes() ([]security.Quote, error) {
	payload := RequestPayload{
		Days:           365 * 30,
		DataPeriod:     "Day",
		DataInterval:   1,
		YFormat:        "0.###",
		ReturnDateType: "ISO8601",
		Elements: []RequestElement{
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
		FinancialTimesUrl,
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

	var result ResponsePayload
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling body: %w", err)
	}

	componentSeries := result.Elements[0].ComponentSeries
	componentIdx := slices.IndexFunc(componentSeries, func(c ComponentSeries) bool {
		return c.Type == "Close"
	})
	if componentIdx == -1 {
		return nil, nil
	}

	component := componentSeries[componentIdx]

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
