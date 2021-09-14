package http

// Routes used in the http handlers.
const (
	RoutePaymentRequest = "api/v1/payment/:paymentID"
	RoutePayment        = "api/v1/payment/:paymentID"

	RouteInvoice  = "api/v1/invoices/:paymentID"
	RouteInvoices = "api/v1/invoices"
	RouteBalance  = "api/v1/balance"

	RouteProofs   = "api/v1/proofs/:txid"
	RouteTxStatus = "api/v1/txstatus/:txid"
)
