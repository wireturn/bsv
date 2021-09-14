package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

type balance struct {
	svc gopayd.BalanceService
}

// NewBalance will setup and return a balance handler.
func NewBalance(svc gopayd.BalanceService) *balance {
	return &balance{svc: svc}
}

// RegisterRoutes will hook up the routes to the echo group.
func (b *balance) RegisterRoutes(g *echo.Group) {
	g.GET(RouteBalance, b.balance)
}

func (b *balance) balance(e echo.Context) error {
	resp, err := b.svc.Balance(e.Request().Context())
	if err != nil {
		return errors.WithStack(err)
	}
	return e.JSON(http.StatusOK, resp)
}
