package preev

// newMockClient returns a client for mocking (using a custom HTTP interface)
func newMockClient(httpClient httpInterface) *Client {
	client := NewClient(nil, nil)
	client.httpClient = httpClient
	return client
}
