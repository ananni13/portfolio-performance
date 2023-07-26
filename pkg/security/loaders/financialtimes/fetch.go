package financialtimes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	financialTimesURL = "https://markets.ft.com/data/chartapi/series"
)

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

func fetchData(symbol string) (responsePayload, error) {
	payload := requestPayload{
		Days:           365 * 30,
		DataPeriod:     "Day",
		DataInterval:   1,
		YFormat:        "0.###",
		ReturnDateType: "ISO8601",
		Elements: []requestElement{
			{
				Type:   "price",
				Symbol: symbol,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return responsePayload{}, fmt.Errorf("error marshaling request body: %w", err)
	}

	res, err := http.Post(
		financialTimesURL,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return responsePayload{}, fmt.Errorf("error during post request: %w", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return responsePayload{}, fmt.Errorf("error reading body: %w", err)
	}

	var result responsePayload
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return responsePayload{}, fmt.Errorf("error unmarshaling body: %w", err)
	}

	return result, nil
}
