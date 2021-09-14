package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	go_payd "github.com/libsv/payd"
	"github.com/pkg/errors"
)

// paymentHandler is an http handler that supports BIP-270 requests.
type paymentHandler struct {
	svc go_payd.PaymentService
}

// NewPaymentHandler will create and return a new PaymentHandler.
func NewPaymentHandler(svc go_payd.PaymentService) *paymentHandler {
	return &paymentHandler{
		svc: svc,
	}
}

// RegisterRoutes will setup all routes with an echo group.
func (h *paymentHandler) RegisterRoutes(g *echo.Group) {
	g.POST(RoutePayment, h.createPayment)
}

func (h *paymentHandler) createPayment(e echo.Context) error {
	args := go_payd.CreatePaymentArgs{
		PaymentID: e.Param("paymentID"),
	}
	var req go_payd.CreatePayment
	if err := Bind(e.Request().Body, &req); err != nil {
		return errors.WithStack(err)
	}
	resp, err := h.svc.CreatePayment(e.Request().Context(), args, req)
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.Error > 0{
		return e.JSON(http.StatusUnprocessableEntity, resp)
	}
	return e.JSON(http.StatusCreated, resp)
}
