package main

import (
	_ "embed"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/enrichman/portfolio-performance/pkg/security"
)

//go:embed securities.csv
var securities []byte

func main() {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	loaders, err := security.LoadSecuritiesFromCSV(securities)
	if err != nil {
		log.Errorf("loading securities from CSV: %s", err)
		os.Exit(1)
	}

	log.Infof("loaded %d securities", len(loaders))

	var wg sync.WaitGroup

	for _, loader := range loaders {
		wg.Add(1)

		loader := loader

		go func() {
			defer wg.Done()
			security.UpdateQuotes(loader)
		}()
	}

	wg.Wait()
}
