package quotes

import (
	"time"
)

// QuoteLoader interface.
type QuoteLoader interface {
	Name() string
	ISIN() string
	LoadQuotes() ([]Quote, error)
}

// Quote struct.
type Quote struct {
	Date  time.Time `json:"date"`
	Close float32   `json:"close"`
}
