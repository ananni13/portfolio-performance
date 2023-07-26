package fondidoc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	fondiDocURLTemplate = "https://www.fondidoc.it/Chart/ChartData?ids=%s&cur=EUR"
)

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
