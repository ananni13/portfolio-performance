package globalbrands

import (
	"github.com/enrichman/portfolio-perfomance/pkg/equity/morganstanley"
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(morganstanley.New(1209, "Global Brands Fund AH (EUR)", "LU0335216932", "Ae"))
}
