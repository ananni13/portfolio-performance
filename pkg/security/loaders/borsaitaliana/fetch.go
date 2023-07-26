package borsaitaliana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	borsaItalianaURL = "https://charts.borsaitaliana.it/charts/services/ChartWService.asmx/GetPricesWithVolume"
)

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

type responsePayload struct {
	Data [][5]float32 `json:"d"`
}

func fetchData(isin, market string) (responsePayload, error) {
	payload := requestPayload{
		SampleTime:           "1d",
		TimeFrame:            "10y",
		RequestedDataSetType: "ohlc",
		ChartPriceType:       "price",
		Key:                  fmt.Sprintf("%s.%s", isin, market),
		KeyType:              "Topic",
		KeyType2:             "Topic",
		Language:             "en-US",
	}

	payloadBytes, err := json.Marshal(struct {
		Request requestPayload `json:"request"`
	}{Request: payload})
	if err != nil {
		return responsePayload{}, fmt.Errorf("error marshaling request body: %w", err)
	}

	res, err := http.Post(
		borsaItalianaURL,
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
