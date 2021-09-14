package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

type txStatusHandler struct {
	svc gopayd.TxStatusService
}

// NewTxStatusHandler will setup and return a new echo txStatus handler.
func NewTxStatusHandler(svc gopayd.TxStatusService) *txStatusHandler {
	return &txStatusHandler{svc: svc}
}

// RegisterRoutes will setup all proof routes with the supplied echo group.
func (t *txStatusHandler) RegisterRoutes(g *echo.Group) {
	g.GET(RouteTxStatus, t.status)
}

func (t *txStatusHandler) status(c echo.Context) error {
	var args gopayd.TxStatusArgs
	if err := c.Bind(&args); err != nil {
		return errors.WithStack(err)
	}
	resp, err := t.svc.Status(c.Request().Context(), args)
	if err != nil {
		return errors.WithStack(err)
	}
	return c.JSON(http.StatusOK, resp)
}
