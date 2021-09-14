package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	go_payd "github.com/libsv/payd"
	"github.com/pkg/errors"
)

// paymentRequestHandler is an http handler that supports BIP-270 requests.
type paymentRequestHandler struct {
	svc go_payd.PaymentRequestService
}

// NewPaymentRequestHandler will create and return a new PaymentRequestHandler.
func NewPaymentRequestHandler(svc go_payd.PaymentRequestService) *paymentRequestHandler {
	return &paymentRequestHandler{
		svc: svc,
	}
}

// RegisterRoutes will setup all routes with an echo group.
func (h *paymentRequestHandler) RegisterRoutes(g *echo.Group) {
	g.GET(RoutePaymentRequest, h.createPaymentRequest)
}

func (h *paymentRequestHandler) createPaymentRequest(e echo.Context) error {
	args := go_payd.PaymentRequestArgs{
		PaymentID: e.Param("paymentID"),
	}
	resp, err := h.svc.CreatePaymentRequest(e.Request().Context(), args)
	if err != nil {
		return errors.WithStack(err)
	}
	return e.JSON(http.StatusCreated, resp)
}
