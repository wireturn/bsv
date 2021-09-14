package envelope

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONEnvelope_IsValid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		env      *JSONEnvelope
		expValid bool
		err      error
	}{
		"valid envelope should return true": {
			env: &JSONEnvelope{
				Payload:   `{\"name\":\"simon\",\"colour\":\"blue\"}`,
				Signature: strToPtr("30450221008209b19ffe2182d859ce36fdeff5ded4b3f70ad77e0e8715238a539db97c1282022043b1a5b260271b7c833ca7c37d1490f21b7bd029dbb8970570c7fdc3df5c93ab"),
				PublicKey: strToPtr("02b01c0c23ff7ff35f774e6d3b3491a123afb6c98965054e024d2320f7dbd25d8a"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: true,
			err:      nil,
		}, "JSON envelope should be valid when checked with correct sig and pk": {
			env: &JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: true,
			err:      nil,
		}, "JSON envelope invalid signature should fail": {
			env: &JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aaaaaa32813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: false,
			err:      nil,
		}, "JSON envelope invalid public key should fail": {
			env: &JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d9999"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: false,
			err:      nil,
		}, "JSON envelope signature not hex should return error": {
			env: &JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1aZZf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d9999"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: false,
			err:      errors.New("failed to decode json envelope signature encoding/hex: invalid byte: U+005A 'Z'"),
		}, "JSON envelope public key not hex should return error": {
			env: &JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d99ZZ"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: false,
			err:      errors.New("failed to decode json envelope publicKey encoding/hex: invalid byte: U+005A 'Z'"),
		}, "JSON envelope signature invalid signature prefix should fail": {
			env: &JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("9945022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("0394890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: false,
			err:      errors.New("failed to parse json envelope signature malformed signature: no header magic"),
		}, "JSON envelope public key invalid prefix should fail": {
			env: &JSONEnvelope{
				Payload:   `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Signature: strToPtr("3045022100b2b3000353b1acaf6e0190a44fc26b0b43830e5aa8d1232813c928d003697c010220294796e63da19d238b29f9cb17e2f31f728ef77a41bfd0f5e355f99f347ff4bf"),
				PublicKey: strToPtr("9994890eeb9888e68cb953d56c598ab0aaa6789e20522cc8b937353694799d7ab1"),
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			expValid: false,
			err:      errors.New("failed to parse json envelope publicKey invalid magic in compressed pubkey string: 153"),
		}, "JSON envelope no sigs should not validate and return valid": {
			env: &JSONEnvelope{
				Payload:  `{\"Test\":\"abc123\",\"Name\":\"4567890\",\"Thing\":\"%$oddchars££$-\"}`,
				Encoding: "UTF-8",
				MimeType: "application/json",
			},
			expValid: true,
			err:      nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			isValid, err := test.env.IsValid()
			if test.err != nil {
				assert.Error(t, err)
				assert.False(t, isValid)
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expValid, isValid)
		})
	}
}

func TestNewJSONEnvelope_signAndCheckValid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		payload interface{}
		err     error
	}{
		"successful run should return no errors": {
			payload: struct {
				Test  string
				Name  string
				Thing string
			}{
				Test:  "abc123",
				Name:  "4567890",
				Thing: "%$oddchars££$-",
			},
			err: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// create env and check valid
			env, err := NewJSONEnvelope(test.payload)
			assert.NoError(t, err)
			val, valErr := env.IsValid()
			assert.NoError(t, valErr)
			assert.True(t, val)
			// convert to json then decode and retry validation
			bb, err := json.Marshal(env)
			assert.NoError(t, err)
			var unmarshalledEnv *JSONEnvelope
			assert.NoError(t, json.Unmarshal(bb, &unmarshalledEnv))
			// is valid?
			val, valErr = unmarshalledEnv.IsValid()
			assert.NoError(t, valErr)
			assert.True(t, val)
		})
	}
}

func strToPtr(s string) *string {
	return &s
}
