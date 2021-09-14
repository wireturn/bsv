package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestClient_GetCapabilities will test the method GetCapabilities()
func TestClient_GetCapabilities(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockCapabilities(http.StatusOK)

		capabilities, err := client.GetCapabilities(testDomain, DefaultPort)
		assert.NoError(t, err)
		assert.NotNil(t, capabilities)
		assert.Equal(t, DefaultBsvAliasVersion, capabilities.BsvAlias)
		assert.Equal(t, http.StatusOK, capabilities.StatusCode)
		assert.Equal(t, true, capabilities.Has(BRFCPki, ""))
	})

	t.Run("status not modified", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockCapabilities(http.StatusNotModified)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
		assert.NoError(t, err)
		assert.NotNil(t, capabilities)
		assert.Equal(t, DefaultBsvAliasVersion, capabilities.BsvAlias)
		assert.Equal(t, http.StatusNotModified, capabilities.StatusCode)
		assert.Equal(t, true, capabilities.Has(BRFCPki, ""))
	})

	t.Run("bad request", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
		assert.Error(t, err)
		assert.NotNil(t, capabilities)
		assert.Equal(t, http.StatusBadRequest, capabilities.StatusCode)
		assert.Equal(t, 0, len(capabilities.Capabilities))
	})

	t.Run("missing target", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities("", DefaultPort)
		assert.Error(t, err)
		assert.Nil(t, capabilities)
	})

	t.Run("missing port", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, 0)
		assert.Error(t, err)
		assert.Nil(t, capabilities)
	})

	t.Run("http error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
		assert.Error(t, err)
		assert.Nil(t, capabilities)
	})

	t.Run("bad error in request", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
		assert.Error(t, err)
		assert.NotNil(t, capabilities)
		assert.Equal(t, http.StatusBadRequest, capabilities.StatusCode)
		assert.Equal(t, 0, len(capabilities.Capabilities))
	})

	t.Run("invalid quotes - good response", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{“`+DefaultServiceName+`“: “`+DefaultBsvAliasVersion+`“,“capabilities“: {“6745385c3fc0“: false,
“pki“: “`+testServerURL+`id/{alias}@{domain.tld}“,“paymentDestination“: “`+testServerURL+`address/{alias}@{domain.tld}“}}`,
			),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
		assert.NoError(t, err)
		assert.NotNil(t, capabilities)
		assert.Equal(t, DefaultBsvAliasVersion, capabilities.BsvAlias)
		assert.Equal(t, http.StatusOK, capabilities.StatusCode)
		assert.Equal(t, true, capabilities.Has(BRFCPki, ""))
	})

	t.Run("invalid alias", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusNotModified,
				`{"`+DefaultServiceName+`": "","capabilities": {"6745385c3fc0": false,"pki": "`+testServerURL+`id/{alias}@{domain.tld}",
"paymentDestination": "`+testServerURL+`address/{alias}@{domain.tld}"}}`,
			),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
		assert.Error(t, err)
		assert.NotNil(t, capabilities)
		assert.NotEqual(t, DefaultBsvAliasVersion, capabilities.BsvAlias)
		assert.Equal(t, http.StatusNotModified, capabilities.StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusNotModified,
				`{"`+DefaultServiceName+`": ,capabilities: {6745385c3fc0: ,pki: `+testServerURL+`id/{alias}@{domain.tld}",
"paymentDestination": "`+testServerURL+`address/{alias}@{domain.tld}"}}`,
			),
		)

		var capabilities *Capabilities
		capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
		assert.Error(t, err)
		assert.NotNil(t, capabilities)
		assert.Equal(t, http.StatusNotModified, capabilities.StatusCode)
		assert.Equal(t, 0, len(capabilities.Capabilities))
	})
}

// mockCapabilities is used for mocking the response
func mockCapabilities(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
		httpmock.NewStringResponder(
			statusCode,
			`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`","capabilities": 
{"6745385c3fc0": false,"pki": "`+testServerURL+`id/{alias}@{domain.tld}",
"paymentDestination": "`+testServerURL+`address/{alias}@{domain.tld}"}}`,
		),
	)
}

// ExampleClient_GetCapabilities example using GetCapabilities()
//
// See more examples in /examples/
func ExampleClient_GetCapabilities() {
	// Load the client
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	mockCapabilities(http.StatusOK)

	// Get the capabilities
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities(testDomain, DefaultPort)
	if err != nil {
		fmt.Printf("error getting capabilities: " + err.Error())
		return
	}
	fmt.Printf("found %d capabilities", len(capabilities.Capabilities))
	// Output:found 3 capabilities
}

// BenchmarkClient_GetCapabilities benchmarks the method GetCapabilities()
func BenchmarkClient_GetCapabilities(b *testing.B) {
	client, _ := newTestClient()
	mockCapabilities(http.StatusOK)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetCapabilities(testDomain, DefaultPort)
	}
}

// TestCapabilities_Has will test the method Has()
func TestCapabilities_Has(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		capabilities  *Capabilities
		brfcID        string
		alternateID   string
		expectedFound bool
	}{
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "alternate_id", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "alternate_id", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "wrong", false},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "6745385c3fc0", true},
	}

	for _, test := range tests {
		if output := test.capabilities.Has(test.brfcID, test.alternateID); output != test.expectedFound {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%v] expected, received: [%v]", t.Name(), test.brfcID, test.alternateID, test.expectedFound, output)
		}
	}
}

// ExampleCapabilities_Has example using Has()
//
// See more examples in /examples/
func ExampleCapabilities_Has() {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	found := capabilities.Has("6745385c3fc0", "alternate_id")
	fmt.Printf("found brfc: %v", found)
	// Output:found brfc: true
}

// BenchmarkCapabilities_Has benchmarks the method Has()
func BenchmarkCapabilities_Has(b *testing.B) {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	for i := 0; i < b.N; i++ {
		_ = capabilities.Has("6745385c3fc0", "alternate_id")
	}
}

// TestCapabilities_GetBool will test the method GetBool()
func TestCapabilities_GetBool(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		capabilities  *Capabilities
		brfcID        string
		alternateID   string
		expectedValue bool
	}{
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "alternate_id", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "alternate_id", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "wrong", false},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "6745385c3fc0", true},
	}

	for _, test := range tests {
		if output := test.capabilities.GetBool(test.brfcID, test.alternateID); output != test.expectedValue {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%v] expected, received: [%v]", t.Name(), test.brfcID, test.alternateID, test.expectedValue, output)
		}
	}
}

// ExampleCapabilities_GetBool example using GetBool()
//
// See more examples in /examples/
func ExampleCapabilities_GetBool() {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	found := capabilities.GetBool("6745385c3fc0", "alternate_id")
	fmt.Printf("found value: %v", found)
	// Output:found value: true
}

// BenchmarkCapabilities_GetBool benchmarks the method GetBool()
func BenchmarkCapabilities_GetBool(b *testing.B) {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	for i := 0; i < b.N; i++ {
		_ = capabilities.GetBool("6745385c3fc0", "alternate_id")
	}
}

// TestCapabilities_GetString will test the method GetString()
func TestCapabilities_GetString(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		capabilities  *Capabilities
		brfcID        string
		alternateID   string
		expectedValue string
	}{
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"pki",
			"0c4339ef99c2",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"0c4339ef99c2",
			"pki",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"0c4339ef99c2",
			"0c4339ef99c2",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"pki",
			"",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"wrong",
			"wrong",
			"",
		},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"wrong",
			"pki",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	for _, test := range tests {
		if output := test.capabilities.GetString(test.brfcID, test.alternateID); output != test.expectedValue {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.brfcID, test.alternateID, test.expectedValue, output)
		}
	}
}

// ExampleCapabilities_GetString example using GetString()
//
// See more examples in /examples/
func ExampleCapabilities_GetString() {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": false,
			"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	found := capabilities.GetString("pki", "0c4339ef99c2")
	fmt.Printf("found value: %v", found)
	// Output:found value: https://domain.com/bsvalias/id/{alias}@{domain.tld}
}

// BenchmarkCapabilities_GetString benchmarks the method GetString()
func BenchmarkCapabilities_GetString(b *testing.B) {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": false,
			"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}
	for i := 0; i < b.N; i++ {
		_ = capabilities.GetString("pki", "0c4339ef99c2")
	}
}
