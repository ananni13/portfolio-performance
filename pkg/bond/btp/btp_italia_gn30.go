package btp

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New("BTP Italia Gn28", "IT0005497000"))
}
