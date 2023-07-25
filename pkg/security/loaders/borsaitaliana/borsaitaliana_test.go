package borsaitaliana

import (
	"strconv"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestNew(t *testing.T) {
	loader, err := New("Name", "ISIN.MARKET")
	if err != nil {
		t.Fatalf("Expected err to be nil, but was: %s", err)
	}

	if loader.name != "Name" {
		t.Errorf("Expected name to be 'Name', but was: %s", loader.name)
	}
	if loader.Name() != "Name" {
		t.Errorf("Expected Name() to return 'Name', but was: %s", loader.name)
	}
	if loader.isin != "ISIN" {
		t.Errorf("Expected isin to be 'ISIN', but was: %s", loader.isin)
	}
	if loader.ISIN() != "ISIN" {
		t.Errorf("Expected ISIN() to return 'ISIN', but was: %s", loader.name)
	}
	if loader.market != "MARKET" {
		t.Errorf("Expected market to be 'MARKET', but was: %s", loader.market)
	}
}

var (
	testTimestamp  = float32(1686614400000)
	testClosePrice = float32(100.11)
	testResponse   = "{\"d\": [[" + strconv.FormatInt(int64(testTimestamp), 10) + ", " + strconv.FormatFloat(float64(testClosePrice), 'f', 2, 32) + ", 100.1, 100.38, 99.96, 100.1],[1686700800000, 100.2, 100.09, 100.3, 100.04, 100.2]]}"
)

func setGock() {
	gock.New("https://charts.borsaitaliana.it").
		Post("/charts/services/ChartWService.asmx/GetPricesWithVolume").
		Reply(200).
		BodyString(testResponse)
}

func TestFetchData(t *testing.T) {
	defer gock.Off()
	setGock()

	response, err := fetchData("ISIN", "MTA")
	if err != nil {
		t.Fatalf("Expected err to be nil, but was: %s", err)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Expected response Data to have len 2, but was: %d", len(response.Data))
	}

	timestamp := response.Data[0][0]
	if timestamp != testTimestamp {
		t.Errorf("Expected timestamp to be %f, but was: %f", testTimestamp, timestamp)
	}

	closePrice := response.Data[0][1]
	if closePrice != testClosePrice {
		t.Errorf("Expected closePrice to be %f, but was: %f", testClosePrice, closePrice)
	}
}

func TestLoadQuotes(t *testing.T) {
	defer gock.Off()
	setGock()

	loader, err := New("Name", "ISIN.MARKET")
	if err != nil {
		t.Fatalf("Expected err to be nil, but was: %s", err)
	}

	quotes, err := loader.LoadQuotes()
	if err != nil {
		t.Fatalf("Expected err to be nil, but was: %s", err)
	}

	if len(quotes) != 2 {
		t.Fatalf("Expected quotes to have len 2, but was: %d", len(quotes))
	}

	quote := quotes[0]
	if quote.Close != testClosePrice {
		t.Errorf("Expected Close to be %f, but was: %f", testClosePrice, quote.Close)
	}

	testDate := time.Unix(int64(testTimestamp/1000), 0).In(time.UTC)
	if quote.Date.Compare(testDate) != 0 {
		t.Errorf("Expected Date to be %s, but was: %s", testDate, quote.Date)
	}
}
