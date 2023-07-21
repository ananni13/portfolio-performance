package secondapensione

import (
	"github.com/enrichman/portfolio-perfomance/pkg/security"
)

func init() {
	security.Register(New("SecondaPensione Sviluppo ESG", "QS0000003564"))
}
