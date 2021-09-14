package whatsonchain

// newMockClient returns a client for mocking
func newMockClient(httpClient httpInterface) *Client {
	client := NewClient(NetworkTest, nil, nil)
	client.httpClient = httpClient
	return client
}
