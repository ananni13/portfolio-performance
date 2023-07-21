package morganstanley

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New(1209, "Global Brands Fund AH (EUR)", "LU0335216932", "Ae"))
}
