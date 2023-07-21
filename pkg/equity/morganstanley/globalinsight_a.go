package morganstanley

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New(34215, "Global Insight Fund A", "LU0868753731", "A"))
}
