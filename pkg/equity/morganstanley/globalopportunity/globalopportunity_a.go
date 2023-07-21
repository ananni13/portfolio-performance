package globalopportunity

import (
	"github.com/enrichman/portfolio-perfomance/pkg/equity/morganstanley"
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(morganstanley.New(34192, "Global Opportunity Fund A", "LU0552385295", "A"))
}
