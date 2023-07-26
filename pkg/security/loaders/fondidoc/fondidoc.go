package fondidoc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
)

const (
	fondiDocURLTemplate = "https://www.fondidoc.it/Chart/ChartData?ids=%s&cur=EUR"
)

// QuoteLoader struct for FondiDoc.
type QuoteLoader struct {
	name   string
	isin   string
	fundID string
}

// New creates a FondiDoc QuoteLoader.
func New(name, isin string) (*QuoteLoader, error) {
	isinID := strings.Split(isin, ".")

	if len(isinID) != 2 {
		return nil, fmt.Errorf("wrong ISIN format for FondiDoc QuoteLoader: \"%s\" - should be \"ISIN.fundID\"", isin)
	}

	return &QuoteLoader{
		name:   name,
		isin:   isinID[0],
		fundID: isinID[1],
	}, nil
}

// Name returns the QuoteLoader name.
func (f *QuoteLoader) Name() string {
	return f.name
}

// ISIN returns the QuoteLoader isin.
func (f *QuoteLoader) ISIN() string {
	return f.isin
}

// FundID returns the QuoteLoader fundID.
func (f *QuoteLoader) FundID() string {
	return f.fundID
}

// LoadQuotes fetches quotes from FondiDoc.
func (f *QuoteLoader) LoadQuotes() ([]quotes.Quote, error) {
	fondiDoc, err := fetchData(f.fundID)
	if err != nil {
		return nil, err
	}

	fundData := fondiDoc[f.fundID].(map[string]interface{})

	quotesData := []quotes.Quote{}
	for _, quote := range fundData["data"].([]interface{}) {
		quotesData = append(quotesData, quotes.Quote{
			Date:  time.Unix(int64(quote.([]interface{})[0].(float64))*100, 0),
			Close: float32(quote.([]interface{})[1].(float64)),
		})
	}

	return quotesData, nil
}

func fetchData(fundID string) (map[string]interface{}, error) {
	url := fmt.Sprintf(fondiDocURLTemplate, fundID)

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

	return fondiDoc, nil
}
