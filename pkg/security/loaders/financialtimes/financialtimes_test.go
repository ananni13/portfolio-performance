package financialtimes

import (
	"strconv"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	loader, err := New("Name", "ISIN.SYMBOL")
	require.Nil(t, err)
	assert.Equal(t, "Name", loader.name)
	assert.Equal(t, "Name", loader.Name())
	assert.Equal(t, "ISIN", loader.isin)
	assert.Equal(t, "ISIN", loader.ISIN())
	assert.Equal(t, "SYMBOL", loader.symbol)
}

var (
	testIsoTimestamp = "2023-01-27T00:00:00"
	testSeriesType   = "Close"
	testClosePrice   = float32(185.48)
	testResponse     = "{\"Dates\": [\"" + testIsoTimestamp + "\", \"2023-01-28T00:00:00\"], \"Elements\": [{\"ComponentSeries\": [{\"Type\": \"" + testSeriesType + "\", \"Values\": [" + strconv.FormatFloat(float64(testClosePrice), 'f', 2, 64) + ", 186.48]}]}]}"
)

func setGock() {
	gock.New("https://markets.ft.com").
		Post("/data/chartapi/series").
		Reply(200).
		BodyString(testResponse)
}

func TestFetchData(t *testing.T) {
	defer gock.Off()
	setGock()

	response, err := fetchData("SYMBOL")
	require.Nil(t, err)
	require.Len(t, response.Dates, 2)
	require.Len(t, response.Elements, 1)
	require.Len(t, response.Elements[0].ComponentSeries, 1)

	assert.Equal(t, testSeriesType, response.Elements[0].ComponentSeries[0].Type)
	require.Len(t, response.Elements[0].ComponentSeries[0].Values, 2)

	assert.Equal(t, response.Dates[0], testIsoTimestamp)
	assert.Equal(t, response.Elements[0].ComponentSeries[0].Values[0], testClosePrice)
}

func TestLoadQuotes(t *testing.T) {
	defer gock.Off()
	setGock()

	loader, err := New("Name", "ISIN.SYMBOL")
	require.Nil(t, err)

	quotes, err := loader.LoadQuotes()
	require.Nil(t, err)
	require.Len(t, quotes, 2)

	quote := quotes[0]
	assert.Equal(t, testClosePrice, quote.Close)

	testDate, err := time.Parse("2006-01-02T15:04:05", testIsoTimestamp)
	require.Nil(t, err)
	assert.Equal(t, testDate, quote.Date)
}
