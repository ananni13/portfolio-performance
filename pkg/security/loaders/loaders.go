package loaders

import (
	"fmt"

	"github.com/enrichman/portfolio-performance/pkg/security/loaders/borsaitaliana"
	"github.com/enrichman/portfolio-performance/pkg/security/loaders/corepension"
	"github.com/enrichman/portfolio-performance/pkg/security/loaders/financialtimes"
	"github.com/enrichman/portfolio-performance/pkg/security/loaders/fondidoc"
	"github.com/enrichman/portfolio-performance/pkg/security/loaders/fonte"
	"github.com/enrichman/portfolio-performance/pkg/security/loaders/morganstanley"
	"github.com/enrichman/portfolio-performance/pkg/security/loaders/secondapensione"
	"github.com/enrichman/portfolio-performance/pkg/security/quotes"
)

func New(loader, name, isin string) (quotes.QuoteLoader, error) {
	switch loader {
	case "borsaitaliana":
		return borsaitaliana.New(name, isin)
	case "financialtimes":
		return financialtimes.New(name, isin)
	case "fonte":
		return fonte.New(name, isin)
	case "secondapensione":
		return secondapensione.New(name, isin)
	case "corepension":
		return corepension.New(name, isin)
	case "fondidoc":
		return fondidoc.New(name, isin)
	case "morganstanley":
		return morganstanley.New(name, isin)
	}
	return nil, fmt.Errorf("quoteLoader [%s] not found", loader)
}
