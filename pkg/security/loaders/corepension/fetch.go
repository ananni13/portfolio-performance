package corepension

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	corePensionURLTemplate = "https://www.corepension.it/product-services/fdr/share/v1/full/%s"
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
	Timestamp string  `json:"date"`
	Value     float64 `json:"value"`
}

func fetchData(isin string) (responsePayload, error) {
	payload := requestPayload{
		Fields: []string{"navHistory"},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	res, err := http.Post(
		fmt.Sprintf(corePensionURLTemplate, isin),
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

	return result, nil
}
