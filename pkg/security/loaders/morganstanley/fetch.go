package morganstanley

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	morganStanleyURLTemplate = "https://www.morganstanley.com/pub/content/dam/im/json/imwebdata/im/data/product/OF/%s/chart/historicalNav.json"
)

type responsePayload struct {
	En en `json:"en"`
}

type en struct {
	ShareClasses []shareClass `json:"shareClasses"`
}

type shareClass struct {
	ID         string     `json:"id"`
	Ccy        string     `json:"ccy"`
	Currencies []currency `json:"currencies"`
}

type currency struct {
	ID     string `json:"id"`
	Series series `json:"series"`
}

type series struct {
	Name     string   `json:"name"`
	Category []string `json:"category"`
	Data     []string `json:"data"`
}

func fetchData(fundID string) (responsePayload, error) {
	url := fmt.Sprintf(morganStanleyURLTemplate, fundID)

	res, err := http.Get(url)
	if err != nil {
		return responsePayload{}, fmt.Errorf("error getting quotes: %w", err)
	}
	if res.StatusCode >= 400 {
		return responsePayload{}, fmt.Errorf("error from request: status_code %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return responsePayload{}, fmt.Errorf("error reading body: %w", err)
	}

	var historicalNav responsePayload
	err = json.Unmarshal(b, &historicalNav)
	if err != nil {
		return responsePayload{}, fmt.Errorf("error unmarshaling body: %w", err)
	}

	return historicalNav, nil
}
