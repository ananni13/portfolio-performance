package fondidoc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/enrichman/portfolio-performance/pkg/security"
)

const (
	fondiDocURLTemplate = "https://www.fondidoc.it/Chart/ChartData?ids=%s&cur=EUR"
)

// QuoteLoader ...
type QuoteLoader struct {
	name   string
	isin   string
	fundID string
}

// New ...
func New(name, isin string) (*QuoteLoader, error) {
	isinID := strings.Split(isin, ".")

	if len(isinID) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for FondiDocQuoteLoader: \"%s\" - should be \"ISIN.fundID\"", isin)
	}

	return &QuoteLoader{
		name:   name,
		isin:   isinID[0],
		fundID: isinID[1],
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

// FundID ...
func (f *QuoteLoader) FundID() string {
	return f.fundID
}

// LoadQuotes ...
func (f *QuoteLoader) LoadQuotes() ([]security.Quote, error) {
	url := fmt.Sprintf(fondiDocURLTemplate, f.fundID)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting quotes: %w", err)
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error from request: status_code %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var fondiDoc map[string]interface{}
	err = json.Unmarshal(b, &fondiDoc)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling body: %w", err)
	}

	fundData := fondiDoc[f.fundID].(map[string]interface{})

	quotes := []security.Quote{}
	for _, quote := range fundData["data"].([]interface{}) {
		quotes = append(quotes, security.Quote{
			Date:  time.Unix(int64(quote.([]interface{})[0].(float64))*100, 0),
			Close: float32(quote.([]interface{})[1].(float64)),
		})
	}

	return quotes, nil
}
