package btp

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New("BTP Italia Mz28", "IT0005532723"))
}
