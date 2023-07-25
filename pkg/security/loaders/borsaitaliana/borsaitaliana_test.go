package borsaitaliana

import (
	"strconv"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	loader, err := New("Name", "ISIN.MARKET")
	require.Nil(t, err)
	assert.Equal(t, "Name", loader.name)
	assert.Equal(t, "Name", loader.Name())
	assert.Equal(t, "ISIN", loader.isin)
	assert.Equal(t, "ISIN", loader.ISIN())
	assert.Equal(t, "MARKET", loader.market)
}

var (
	testTimestamp  = float32(1686614400000)
	testClosePrice = float32(100.11)
	testResponse   = "{\"d\": [[" + strconv.FormatInt(int64(testTimestamp), 10) + ", " + strconv.FormatFloat(float64(testClosePrice), 'f', 2, 64) + ", 100.1, 100.38, 99.96, 100.1],[1686700800000, 100.2, 100.09, 100.3, 100.04, 100.2]]}"
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

	response, err := fetchData("ISIN", "MARKET")
	require.Nil(t, err)
	require.Len(t, response.Data, 2)

	timestamp := response.Data[0][0]
	assert.Equal(t, testTimestamp, timestamp)

	closePrice := response.Data[0][1]
	assert.Equal(t, testClosePrice, closePrice)
}

func TestLoadQuotes(t *testing.T) {
	defer gock.Off()
	setGock()

	loader, err := New("Name", "ISIN.MARKET")
	require.Nil(t, err)

	quotes, err := loader.LoadQuotes()
	require.Nil(t, err)
	require.Len(t, quotes, 2)

	quote := quotes[0]
	assert.Equal(t, testClosePrice, quote.Close)

	testDate := time.Unix(int64(testTimestamp/1000), 0).In(time.UTC)
	assert.Equal(t, testDate, quote.Date)
}
