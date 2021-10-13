package handcash

// CurrencyCode is an enum for supported currencies
type CurrencyCode string

// CurrencyCode enums
const (
	CurrencyARS CurrencyCode = "ARS"
	CurrencyAUD CurrencyCode = "AUD"
	CurrencyBRL CurrencyCode = "BRL"
	CurrencyBSV CurrencyCode = "BSV"
	CurrencyCAD CurrencyCode = "CAD"
	CurrencyCHF CurrencyCode = "CHF"
	CurrencyCNY CurrencyCode = "CNY"
	CurrencyCOP CurrencyCode = "COP"
	CurrencyCZK CurrencyCode = "CZK"
	CurrencyDKK CurrencyCode = "DKK"
	CurrencyEUR CurrencyCode = "EUR"
	CurrencyGBP CurrencyCode = "GBP"
	CurrencyHKD CurrencyCode = "HKD"
	CurrencyJPY CurrencyCode = "JPY"
	CurrencyMXN CurrencyCode = "MXN"
	CurrencyNOK CurrencyCode = "NOK"
	CurrencyNZD CurrencyCode = "NZD"
	CurrencyPHP CurrencyCode = "PHP"
	CurrencyRUB CurrencyCode = "RUB"
	CurrencySAT CurrencyCode = "SAT"
	CurrencySEK CurrencyCode = "SEK"
	CurrencySGD CurrencyCode = "SGD"
	CurrencyTHB CurrencyCode = "THB"
	CurrencyUSD CurrencyCode = "USD"
	CurrencyZAR CurrencyCode = "ZAR"
)

// Profile are the user fields returned by the public and private profile endpoints
type Profile struct {
	PublicProfile  PublicProfile  `json:"publicProfile"`
	PrivateProfile PrivateProfile `json:"privateProfile"`
}

// PublicProfile is the public profile
type PublicProfile struct {
	AvatarURL         string       `json:"avatarUrl"`
	DisplayName       string       `json:"displayName"`
	Handle            string       `json:"handle"`
	ID                string       `json:"id"`
	LocalCurrencyCode CurrencyCode `json:"localCurrencyCode"`
	Paymail           string       `json:"paymail"`
	BitcoinUnit       string       `json:"bitcoinUnit"`
}

// PrivateProfile is the private profile
type PrivateProfile struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

// Participant is used for payments
type Participant struct {
	Alias             string          `json:"alias"`
	DisplayName       string          `json:"displayName"`
	ProfilePictureURL string          `json:"profilePictureUrl"`
	ResponseNote      string          `json:"responseNote"`
	Type              ParticipantType `json:"type"`
}

// PaymentResponse is returned from the GetPayment function
type PaymentResponse struct {
	AppAction         AppAction      `json:"appAction"`
	Attachments       []*Attachment  `json:"attachments"`
	FiatCurrencyCode  CurrencyCode   `json:"fiatCurrencyCode"`
	FiatExchangeRate  float64        `json:"fiatExchangeRate"`
	Note              string         `json:"note"`
	Participants      []*Participant `json:"participants"`
	RawTransactionHex string         `json:"rawTransactionHex,omitempty"`
	SatoshiAmount     uint64         `json:"satoshiAmount"`
	SatoshiFees       uint64         `json:"satoshiFees"`
	Time              uint64         `json:"time"`
	TransactionID     string         `json:"transactionId"`
	Type              PaymentType    `json:"type"`
}

// PaymentRequest is used for GetPayment()
type PaymentRequest struct {
	TransactionID string `json:"transactionId"`
}

// BalanceRequest is used for GetSpendableBalance()
type BalanceRequest struct {
	CurrencyCode CurrencyCode `json:"currencyCode"`
}

// AppAction enum
type AppAction string

// AppAction enum
const (
	AppActionLike     AppAction = "like"
	AppActionPublish  AppAction = "publish"
	AppActionTip      AppAction = "tip"
	AppActionTipGroup AppAction = "tip-group"
)

// AttachmentFormat enum
type AttachmentFormat string

// AttachmentFormat enum
const (
	AttachmentFormatBase64 AttachmentFormat = "base64"
	AttachmentFormatHex    AttachmentFormat = "hex"
	AttachmentFormatJSON   AttachmentFormat = "json"
)

// Attachment is for additional data
type Attachment struct {
	Format AttachmentFormat `json:"format,omitempty"`
	Value  interface{}      `json:"value,omitempty"`
}

// Payment is used by PayParameters
type Payment struct {
	Amount       float64      `json:"amount"`
	CurrencyCode CurrencyCode `json:"currencyCode"`
	To           string       `json:"to"`
}

// PayParameters is used by Pay()
type PayParameters struct {
	AppAction   AppAction   `json:"appAction,omitempty"`
	Attachment  *Attachment `json:"attachment,omitempty"`
	Description string      `json:"description,omitempty"`
	Receivers   []*Payment  `json:"receivers,omitempty"`
}

// PaymentType enum
type PaymentType string

// PaymentType enum
const (
	PaymentSend PaymentType = "send"
)

// ParticipantType enum
type ParticipantType string

// ParticipantType enum
const (
	ParticipantUser ParticipantType = "user"
)

// oAuthHeaders are used for signed requests
type oAuthHeaders struct {
	OauthPublicKey string `json:"oauth-publickey"`
	OauthSignature string `json:"oauth-signature"`
	OauthTimestamp string `json:"oauth-timestamp"`
}

// signedRequest is used to communicate with HandCash Connect API
type signedRequest struct {
	Body    interface{}  `json:"body"`
	Headers oAuthHeaders `json:"headers"`
	JSON    bool         `json:"json"`
	Method  string       `json:"method"`
	URI     string       `json:"uri"`
}

// requestBody is for constructing the request
type requestBody struct {
	authToken string
}

// errorResponse is the error response
type errorResponse struct {
	Message string `json:"message"`
}
