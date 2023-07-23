package fondidoc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

const (
	FondiDocUrlTemplate = "https://www.fondidoc.it/Chart/ChartData?ids=%s&cur=EUR"
)

type FondiDocQuoteLoader struct {
	name   string
	isin   string
	fundId string
}

func New(name, isin string) (*FondiDocQuoteLoader, error) {
	isinId := strings.Split(isin, ".")

	if len(isinId) != 2 {
		return nil, fmt.Errorf("Wrong ISIN format for FondiDocQuoteLoader: \"%s\" - should be \"ISIN.fundId\"", isin)
	}

	return &FondiDocQuoteLoader{
		name:   name,
		isin:   isinId[0],
		fundId: isinId[1],
	}, nil
}

func (f *FondiDocQuoteLoader) Name() string {
	return f.name
}

func (f *FondiDocQuoteLoader) ISIN() string {
	return f.isin
}

func (f *FondiDocQuoteLoader) FundId() string {
	return f.fundId
}

func (f *FondiDocQuoteLoader) LoadQuotes() ([]security.Quote, error) {
	url := fmt.Sprintf(FondiDocUrlTemplate, f.fundId)

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

	fundData := fondiDoc[f.fundId].(map[string]interface{})

	quotes := []security.Quote{}
	for _, quote := range fundData["data"].([]interface{}) {
		quotes = append(quotes, security.Quote{
			Date:  time.Unix(int64(quote.([]interface{})[0].(float64))*100, 0),
			Close: float32(quote.([]interface{})[1].(float64)),
		})
	}

	return quotes, nil
}
