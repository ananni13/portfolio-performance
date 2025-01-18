package corepension

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

const (
	corePensionURLTemplate = "https://www.corepension.it/product-services/fdr/share/v2/full/%s"
)

type requestPayload struct {
	Fields []string `json:"fields"`
}

type responsePayload []responseItem

type responseItem struct {
	ID         string       `json:"_id"`
	NavHistory []parsedData `json:"navHistory"`
}

type parsedData struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

func fetchData(isin string) (responsePayload, error) {
	payload := requestPayload{
		Fields: []string{"navHistory"},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(corePensionURLTemplate, isin), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", "www.corepension.it")
	req.Header.Set("Origin", "https://www.corepension.it")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(req)
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
		log.Info(string(bodyBytes))
		log.Infof("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Infof("syntax error at byte offset %d", e.Offset)
		}
		return nil, fmt.Errorf("error unmarshaling body: %w", err)
	}

	return result, nil
}
