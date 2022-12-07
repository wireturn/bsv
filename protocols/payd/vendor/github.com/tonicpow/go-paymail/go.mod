module github.com/tonicpow/go-paymail

go 1.15

require (
	github.com/bitcoinschema/go-bitcoin v0.3.15
	github.com/bitcoinsv/bsvd v0.0.0-20190609155523-4c29707f7173
	github.com/bitcoinsv/bsvutil v0.0.0-20181216182056-1d77cf353ea9
	github.com/go-resty/resty/v2 v2.5.0
	github.com/jarcoal/httpmock v1.0.8
	github.com/julienschmidt/httprouter v1.3.0
	github.com/miekg/dns v1.1.40
	github.com/mrz1836/go-api-router v0.3.9
	github.com/mrz1836/go-logger v0.2.4
	github.com/mrz1836/go-sanitize v1.1.3
	github.com/mrz1836/go-validate v0.2.0
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.2.4 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/sys v0.0.0-20210301091718-77cc2087c03b // indirect
	golang.org/x/text v0.3.5 // indirect
)

replace github.com/go-resty/resty/v2 => github.com/go-resty/resty/v2 v2.4.0
