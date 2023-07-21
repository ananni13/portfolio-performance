package globalopportunity

import (
	"github.com/enrichman/portfolio-perfomance/pkg/equity/morganstanley"
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(morganstanley.New(34192, "Global Opportunity Fund AH (EUR)", "LU0552385618", "Ae"))
}
