package btp

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New("BTP Valore Gn27", "IT0005547408"))
}
