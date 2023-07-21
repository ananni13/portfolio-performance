package morganstanley

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New(34215, "Global Insight Fund AH (EUR)", "LU0868754382", "Ae"))
}
