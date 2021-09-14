// Package b is a library for working with B transactions (Bitcoin OP_RETURN protocol) in Go
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BitcoinSchema Organization (https://bitcoinschema.org)
package b

import (
	"encoding/base64"
	"fmt"
)

// Prefix is the Bitcom prefix used by B
const Prefix = "19HxigV4QyBv3tHpQVcUEQyq1pzZVdoAut"

// Data is the content portion of the B data
type Data struct {
	Bytes []byte `json:"data,omitempty"`
	UTF8  string `json:"utf8,omitempty"`
}

// B is B protocol
type B struct {
	Data      Data
	MediaType string `json:"media_type"`
	Encoding  string `json:"encoding"`
	Filename  string `json:"filename,omitempty"`
}

// EncodingType is an enum for the possible types of data encoding
type EncodingType string

// Various encoding types
const (
	EncodingBinary  EncodingType = "binary"
	EncodingGzip    EncodingType = "gzip"
	EncodingUtf8    EncodingType = "utf8"
	EncodingUtf8Alt EncodingType = "utf-8"
)

// DataURI returns a b64 encoded image that can be set directly. Ex: <img src="b64data" />
func (b *B) DataURI() string {
	return fmt.Sprintf("data:%s;base64,%s", b.Encoding, base64.StdEncoding.EncodeToString(b.Data.Bytes))
}

// BitFsURL is a helper to create a bitfs url to fetch the content over HTTP
func BitFsURL(txID string, outIndex, scriptChunk int) string {
	return fmt.Sprintf("https://x.bitfs.network/%s.out.%d.%d", txID, outIndex, scriptChunk)
}
