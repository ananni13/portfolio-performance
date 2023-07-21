package morganstanley

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New(1209, "Global Brands Fund AX", "LU0239683559", "AX"))
}
