package morganstanley

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New(1209, "Global Brands Fund AHX (EUR)", "LU0552899998", "A3"))
}
