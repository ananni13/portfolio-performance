package main

import (
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security"
)

func main() {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	securities, err := security.LoadSecuritiesFromCSV("securities.csv")
	if err != nil {
		log.Errorf("loading securities from CSV: %s", err)
		os.Exit(1)
	}

	log.Infof("loaded %d securities", len(securities))

	var wg sync.WaitGroup

	for _, loader := range securities {
		wg.Add(1)

		loader := loader

		go func() {
			defer wg.Done()
			security.UpdateQuotes(loader)
		}()
	}

	wg.Wait()
}
