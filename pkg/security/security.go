package security

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security/loaders"
	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
	"golang.org/x/exp/maps"
)

// LoadSecuritiesFromCSV loads all securities from the securities.csv file and returns a slice of corresponding QuoteLoader
func LoadSecuritiesFromCSV(csvBytes []byte) ([]quotes.QuoteLoader, error) {
	// read csv values using csv.Reader
	csvReader := csv.NewReader(bytes.NewReader(csvBytes))
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = 3
	_, err := csvReader.Read() // skip header line
	if err != nil {
		return nil, fmt.Errorf("reading csv: %w", err)
	}

	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading csv: %w", err)
	}

	securities := make(map[string]quotes.QuoteLoader)

	for _, line := range data {
		isin := line[0]
		name := line[1]
		loader := line[2]

		quoteLoader, err := loaders.New(loader, name, isin)
		if err != nil {
			log.Errorf("Error creating quoteLoader [%s] for ISIN %s (%s): %s", loader, isin, name, err)
			continue
		}

		if _, found := securities[quoteLoader.ISIN()]; found {
			log.Warnf("security '%s' already registered", quoteLoader.ISIN())
			continue
		}

		securities[quoteLoader.ISIN()] = quoteLoader
		log.Infof("security '%s' registered", quoteLoader.ISIN())
	}

	return maps.Values(securities), nil
}

// UpdateQuotes fetches and updates quotes from the QuoteLoader
func UpdateQuotes(loader quotes.QuoteLoader) {
	start := time.Now().In(time.UTC)

	log.Infof("[%s] loading quotes for '%s'", loader.ISIN(), loader.Name())

	newQuotes, err := loader.LoadQuotes()
	if err != nil {
		log.Errorf("[%s] error loading quotes: %s", loader.ISIN(), err)
		return
	}
	if len(newQuotes) == 0 {
		log.Warnf("[%s] no quotes found", loader.ISIN())
		return
	}

	log.Debugf("[%s] new quotes loaded from %s to %s",
		loader.ISIN(),
		newQuotes[0].Date,
		newQuotes[len(newQuotes)-1].Date,
	)

	filename := fmt.Sprintf("out/json/%s.json", loader.ISIN())
	log.Debugf("[%s] loading OLD quotes from '%s'", loader.ISIN(), filename)

	oldQuotes, err := loadQuotesFromFile(filename)
	if err != nil {
		log.Errorf("[%s] error loading quotes: %s", loader.ISIN(), err.Error())
		return
	}

	if len(oldQuotes) == 0 {
		log.Warnf("[%s] no OLD quotes found", loader.ISIN())
	} else {
		log.Debugf("[%s] found OLD quotes from %s to %s",
			loader.ISIN(),
			oldQuotes[0].Date,
			oldQuotes[len(oldQuotes)-1].Date,
		)
	}

	mergedQuotes := merge(oldQuotes, newQuotes)
	log.Debugf("[%s] merged quotes from %s to %s",
		loader.ISIN(),
		mergedQuotes[0].Date,
		mergedQuotes[len(mergedQuotes)-1].Date,
	)

	err = writeQuotesToFile(filename, mergedQuotes)
	if err != nil {
		log.Errorf("[%s] error writing quotes: %s", loader.ISIN(), err.Error())
		return
	}

	addedQuotes := len(mergedQuotes) - len(oldQuotes)
	if addedQuotes == 0 {
		log.Infof("[%s] no new quotes added", loader.ISIN())
	} else {
		log.Infof(
			"[%s] new quotes added [%d] - old [%d] - new [%d]",
			loader.ISIN(), addedQuotes, len(oldQuotes), len(newQuotes),
		)
	}

	log.Infof("[%s] quotes loaded in %s", loader.ISIN(), time.Since(start))
}

func merge(quotes1 []quotes.Quote, quotes2 []quotes.Quote) []quotes.Quote {
	quotesMap := map[time.Time]quotes.Quote{}

	for _, q := range quotes1 {
		q.Date = q.Date.UTC()
		quotesMap[q.Date] = q
	}

	for _, q := range quotes2 {
		q.Date = q.Date.UTC()

		if oldQuote, found := quotesMap[q.Date]; found {
			if oldQuote.Close != q.Close {
				log.Warnf("quote for date '%v' already exists with different value [old: %v - new: %v]",
					q.Date, oldQuote.Close, q.Close,
				)
			}
		}
		quotesMap[q.Date] = q
	}

	mergedQuotes := []quotes.Quote{}
	for _, v := range quotesMap {
		mergedQuotes = append(mergedQuotes, v)
	}

	sort.Slice(mergedQuotes, func(i, j int) bool {
		return mergedQuotes[i].Date.Before(mergedQuotes[j].Date)
	})

	return mergedQuotes
}

func loadQuotesFromFile(filename string) ([]quotes.Quote, error) {
	oldQuotesByte, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading file [%s]: %s", filename, err.Error())
	}

	var oldQuotes []quotes.Quote
	err = json.Unmarshal(oldQuotesByte, &oldQuotes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling file [%s]: %s", filename, err.Error())
	}

	return oldQuotes, nil
}

func writeQuotesToFile(filename string, quotes []quotes.Quote) error {
	jsonOutput, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling file [%s]: %s", filename, err.Error())
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file [%s]: %s", filename, err.Error())
	}
	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating file [%s]: %s", filename, err.Error())
	}

	if _, err = file.Write(jsonOutput); err != nil {
		return fmt.Errorf("error writing to file [%s]: %s", filename, err.Error())
	}

	return nil
}
