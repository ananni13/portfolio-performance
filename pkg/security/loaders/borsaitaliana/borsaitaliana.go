package borsaitaliana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/enrichman/portfolio-performance/pkg/security"
)

const (
	borsaItalianaURL = "https://charts.borsaitaliana.it/charts/services/ChartWService.asmx/GetPricesWithVolume"
)

// QuoteLoader ...
type QuoteLoader struct {
	name   string
	isin   string
	market string
}

// New ...
func New(name, isin string) (*QuoteLoader, error) {
	isinMarket := strings.Split(isin, ".")

	if len(isinMarket) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for BorsaItalianaQuoteLoader: \"%s\" - should be \"ISIN.market\"", isin)
	}

	return &QuoteLoader{
		name:   name,
		isin:   isinMarket[0],
		market: isinMarket[1],
	}, nil
}

// Name ...
func (b *QuoteLoader) Name() string {
	return b.name
}

// ISIN ...
func (b *QuoteLoader) ISIN() string {
	return b.isin
}

type data struct {
	Data [][5]float32 `json:"d"`
}

type requestPayload struct {
	SampleTime           string
	TimeFrame            string
	RequestedDataSetType string
	ChartPriceType       string
	Key                  string
	OffSet               int
	FromDate             string `json:",omitempty"`
	ToDate               string `json:",omitempty"`
	UseDelay             bool
	KeyType              string
	KeyType2             string
	Language             string
}

// LoadQuotes ...
func (b *QuoteLoader) LoadQuotes() ([]security.Quote, error) {
	payload := requestPayload{
		SampleTime:           "1d",
		TimeFrame:            "10y",
		RequestedDataSetType: "ohlc",
		ChartPriceType:       "price",
		Key:                  fmt.Sprintf("%s.%s", b.isin, b.market),
		KeyType:              "Topic",
		KeyType2:             "Topic",
		Language:             "en-US",
	}

	payloadBytes, err := json.Marshal(struct {
		Request requestPayload `json:"request"`
	}{Request: payload})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body")
	}

	res, err := http.Post(
		borsaItalianaURL,
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

	var result data
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling body: %w", err)
	}

	quotes := []security.Quote{}
	for _, quote := range result.Data {
		quotes = append(quotes, security.Quote{
			Date:  time.Unix(int64(quote[0]/1000), 0).In(time.UTC),
			Close: quote[1],
		})
	}

	return quotes, nil
}
