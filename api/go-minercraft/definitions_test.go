package minercraft

import (
	"encoding/json"
	"testing"

	"github.com/libsv/go-bk/envelope"
	"github.com/stretchr/testify/assert"
)

func Test_JSONEnvelope_process(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		env envelope.JSONEnvelope
		exp *JSONEnvelope
		err error
	}{
		"JSONEnvelope with no sig should map correctly": {
			env: envelope.JSONEnvelope{
				Payload:  "{\"index\":7\"}",
				MimeType: "application/json",
			},
			exp: &JSONEnvelope{
				Miner:     nil,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:  "{\"index\":7\"}",
					MimeType: "application/json",
				},
			},
			err: nil,
		}, "JSONEnvelope with sig should map correctly and set validated to true": {
			env: envelope.JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			exp: &JSONEnvelope{
				Miner:     nil,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
					Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
					PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
					Encoding:  "UTF-8",
					MimeType:  "application/json",
				},
			},
			err: nil,
		}, "JSONEnvelope with modified payload should map correctly and set validated to false": {
			env: envelope.JSONEnvelope{
				Payload:   `{\"Test\":\"abc124\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			exp: &JSONEnvelope{
				Miner:     nil,
				Validated: false,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   `{\"Test\":\"abc124\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
					Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
					PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
					Encoding:  "UTF-8",
					MimeType:  "application/json",
				},
			},
			err: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bb, err := json.Marshal(test.env)
			assert.NoError(t, err)
			env := &JSONEnvelope{}
			assert.NoError(t, env.process(nil, bb))
			assert.Equal(t, test.exp, env)
		})
	}
}

func strToPtr(s string) *string {
	return &s
}
