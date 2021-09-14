package gopayd

// MerchantData to be displayed to the user.
type MerchantData struct {
	// AvatarURL displays a canonical url to a merchants avatar.
	AvatarURL string `json:"avatar"`
	// MerchantName is a human readable string identifying the merchant.
	MerchantName string `json:"name"`
	// Email can be sued to contact the merchant about this transaction.
	Email string `json:"email"`
	// Address is the merchants store / head office address.
	Address string `json:"address"`
	// PaymentReference can be sent to link this request with a specific payment id.
	PaymentReference string `json:"paymentReference"`
	// ExtendedData can be supplied if the merchant wishes to send some arbitrary data back to the wallet.
	ExtendedData map[string]interface{} `json:"extendedData,omitempty"`
}
