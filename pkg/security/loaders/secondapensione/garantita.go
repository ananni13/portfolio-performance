package secondapensione

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New("SecondaPensione Garantita ESG", "QS0000013033"))
}
