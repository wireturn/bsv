package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

type invoice struct {
	svc gopayd.InvoiceService
}

// NewInvoice will setup and return a new invoice handler.
func NewInvoice(svc gopayd.InvoiceService) *invoice {
	return &invoice{svc: svc}
}

// RegisterRoutes will hook up the routes to the echo group.
func (i *invoice) RegisterRoutes(g *echo.Group) {
	g.GET(RouteInvoices, i.invoices)
	g.GET(RouteInvoice, i.invoice)
	g.POST(RouteInvoices, i.create)
	g.DELETE(RouteInvoice, i.delete)
}

// invoices returns all invoices currently stored.
func (i *invoice) invoices(e echo.Context) error {
	ii, err := i.svc.Invoices(e.Request().Context())
	if err != nil {
		return errors.WithStack(err)
	}
	return e.JSON(http.StatusOK, ii)
}

func (i *invoice) invoice(e echo.Context) error {
	var args gopayd.InvoiceArgs
	if err := e.Bind(&args); err != nil {
		return errors.Wrap(err, "failed to parse invoice args")
	}
	inv, err := i.svc.Invoice(e.Request().Context(), args)
	if err != nil {
		return errors.WithStack(err)
	}
	return e.JSON(http.StatusOK, inv)
}

func (i *invoice) create(e echo.Context) error {
	var req gopayd.InvoiceCreate
	if err := e.Bind(&req); err != nil {
		return errors.Wrap(err, "failed to parse invoice create req")
	}
	inv, err := i.svc.Create(e.Request().Context(), req)
	if err != nil {
		return errors.WithStack(err)
	}
	return e.JSON(http.StatusCreated, inv)
}

func (i *invoice) delete(e echo.Context) error {
	var args gopayd.InvoiceArgs
	if err := e.Bind(&args); err != nil {
		return errors.Wrap(err, "failed to parse invoice args")
	}
	if err := i.svc.Delete(e.Request().Context(), args); err != nil {
		return errors.WithStack(err)
	}
	return e.NoContent(http.StatusNoContent)
}
