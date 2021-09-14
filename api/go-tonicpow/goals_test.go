package tonicpow

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestGoal will return a dummy example for tests
func newTestGoal() *Goal {
	return &Goal{
		CampaignID:     testCampaignID,
		Description:    "This is an example goal",
		ID:             testGoalID,
		MaxPerPromoter: 1,
		Name:           testGoalName,
		PayoutRate:     0.01,
		PayoutType:     "flat",
		Title:          "Example Goal",
	}
}

// TestClient_CreateGoal will test the method CreateGoal()
func TestClient_CreateGoal(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("create a goal (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, goal)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateGoal(goal)
		assert.NoError(t, err)
		assert.NotNil(t, goal)
		assert.NotNil(t, response)
		assert.Equal(t, testGoalID, goal.ID)
	})

	t.Run("missing campaign id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()
		goal.CampaignID = 0

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, goal)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateGoal(goal)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("missing goal name", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()
		goal.Name = ""

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, goal)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateGoal(goal)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusBadRequest, goal)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateGoal(goal)
		assert.Error(t, err)
		assert.NotNil(t, response)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodPost, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateGoal(goal)
		assert.Error(t, err)
		assert.Equal(t, apiError.Message, err.Error())
		assert.NotNil(t, response)
	})
}

// ExampleClient_CreateGoal example using CreateGoal()
//
// See more examples in /examples/
func ExampleClient_CreateGoal() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseGoal := newTestGoal()
	_ = mockResponseData(
		http.MethodPost,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal),
		http.StatusCreated,
		responseGoal,
	)

	// Create goal (using mocking response)
	if _, err = client.CreateGoal(responseGoal); err != nil {
		fmt.Printf("error creating goal: " + err.Error())
		return
	}
	fmt.Printf("created goal: %s", responseGoal.Name)
	// Output:created goal: example_goal
}

// BenchmarkClient_CreateGoal benchmarks the method CreateGoal()
func BenchmarkClient_CreateGoal(b *testing.B) {
	client, _ := newTestClient()
	goal := newTestGoal()
	_ = mockResponseData(
		http.MethodPost,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal),
		http.StatusCreated,
		goal,
	)
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateGoal(goal)
	}
}

// TestClient_GetGoal will test the method GetGoal()
func TestClient_GetGoal(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("get a goal (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf(
			"%s/%s/details/%d", EnvironmentDevelopment.apiURL,
			modelGoal, goal.ID,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, goal)
		assert.NoError(t, err)

		var newGoal *Goal
		var response *StandardResponse
		newGoal, response, err = client.GetGoal(goal.ID)
		assert.NoError(t, err)
		assert.NotNil(t, newGoal)
		assert.NotNil(t, response)
		assert.Equal(t, testGoalID, goal.ID)
	})

	t.Run("missing goal id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()
		goal.ID = 0

		endpoint := fmt.Sprintf(
			"%s/%s/details/%d", EnvironmentDevelopment.apiURL,
			modelGoal, goal.ID,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, goal)
		assert.NoError(t, err)

		var newGoal *Goal
		var response *StandardResponse
		newGoal, response, err = client.GetGoal(goal.ID)
		assert.Error(t, err)
		assert.Nil(t, newGoal)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf(
			"%s/%s/details/%d", EnvironmentDevelopment.apiURL,
			modelGoal, goal.ID,
		)
		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, goal)
		assert.NoError(t, err)

		var newGoal *Goal
		var response *StandardResponse
		newGoal, response, err = client.GetGoal(goal.ID)
		assert.Error(t, err)
		assert.Nil(t, newGoal)
		assert.NotNil(t, response)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf(
			"%s/%s/details/%d", EnvironmentDevelopment.apiURL,
			modelGoal, goal.ID,
		)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var newGoal *Goal
		var response *StandardResponse
		newGoal, response, err = client.GetGoal(goal.ID)
		assert.Error(t, err)
		assert.Nil(t, newGoal)
		assert.NotNil(t, response)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// ExampleClient_GetGoal example using GetGoal()
//
// See more examples in /examples/
func ExampleClient_GetGoal() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseGoal := newTestGoal()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/details/%d", EnvironmentDevelopment.apiURL,
			modelGoal, responseGoal.ID,
		),
		http.StatusOK,
		responseGoal,
	)

	// Get goal (using mocking response)
	if responseGoal, _, err = client.GetGoal(responseGoal.ID); err != nil {
		fmt.Printf("error getting goal: " + err.Error())
		return
	}
	fmt.Printf("goal: %s", responseGoal.Name)
	// Output:goal: example_goal
}

// BenchmarkClient_GetGoal benchmarks the method GetGoal()
func BenchmarkClient_GetGoal(b *testing.B) {
	client, _ := newTestClient()
	goal := newTestGoal()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/details/%d", EnvironmentDevelopment.apiURL,
			modelGoal, goal.ID,
		),
		http.StatusOK,
		goal,
	)
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetGoal(goal.ID)
	}
}

// TestClient_UpdateCampaign will test the method UpdateCampaign()
func TestClient_UpdateGoal(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("update a goal (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()
		goal.Title = "Updated Title"

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		err = mockResponseData(http.MethodPut, endpoint, http.StatusOK, goal)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateGoal(goal)
		assert.NoError(t, err)
		assert.NotNil(t, goal)
		assert.NotNil(t, response)
		assert.Equal(t, "Updated Title", goal.Title)
	})

	t.Run("missing id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()
		goal.ID = 0

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		err = mockResponseData(http.MethodPut, endpoint, http.StatusOK, goal)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateGoal(goal)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		err = mockResponseData(http.MethodPut, endpoint, http.StatusBadRequest, goal)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateGoal(goal)
		assert.Error(t, err)
		assert.NotNil(t, response)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodPut, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateGoal(goal)
		assert.Error(t, err)
		assert.Equal(t, apiError.Message, err.Error())
		assert.NotNil(t, response)
	})
}

// ExampleClient_UpdateGoal example using UpdateGoal()
//
// See more examples in /examples/
func ExampleClient_UpdateGoal() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseGoal := newTestGoal()
	responseGoal.Title = "Updated Title"
	_ = mockResponseData(
		http.MethodPut,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal),
		http.StatusOK,
		responseGoal,
	)

	// Update goal (using mocking response)
	_, err = client.UpdateGoal(responseGoal)
	if err != nil {
		fmt.Printf("error updating goal: " + err.Error())
		return
	}
	fmt.Printf("goal: %s", responseGoal.Title)
	// Output:goal: Updated Title
}

// BenchmarkClient_UpdateGoal benchmarks the method UpdateGoal()
func BenchmarkClient_UpdateGoal(b *testing.B) {
	client, _ := newTestClient()
	goal := newTestGoal()
	_ = mockResponseData(
		http.MethodPut,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelGoal),
		http.StatusOK,
		goal,
	)
	for i := 0; i < b.N; i++ {
		_, _ = client.UpdateGoal(goal)
	}
}

// TestClient_DeleteGoal will test the method DeleteGoal()
func TestClient_DeleteGoal(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("delete a goal (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s?%s=%d", EnvironmentDevelopment.apiURL, modelGoal, fieldID, goal.ID)
		err = mockResponseData(http.MethodDelete, endpoint, http.StatusOK, nil)
		assert.NoError(t, err)

		var deleted bool
		var response *StandardResponse
		deleted, response, err = client.DeleteGoal(goal.ID)
		assert.NoError(t, err)
		assert.Equal(t, true, deleted)
		assert.Equal(t, testGoalID, goal.ID)
		assert.NotNil(t, response)
	})

	t.Run("missing goal id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()
		goal.ID = 0

		endpoint := fmt.Sprintf("%s/%s?%s=%d", EnvironmentDevelopment.apiURL, modelGoal, fieldID, goal.ID)

		err = mockResponseData(http.MethodDelete, endpoint, http.StatusOK, nil)
		assert.NoError(t, err)

		var deleted bool
		var response *StandardResponse
		deleted, response, err = client.DeleteGoal(goal.ID)
		assert.Error(t, err)
		assert.Equal(t, false, deleted)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s?%s=%d", EnvironmentDevelopment.apiURL, modelGoal, fieldID, goal.ID)
		err = mockResponseData(http.MethodDelete, endpoint, http.StatusBadRequest, nil)
		assert.NoError(t, err)

		var deleted bool
		var response *StandardResponse
		deleted, response, err = client.DeleteGoal(goal.ID)
		assert.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, false, deleted)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		goal := newTestGoal()

		endpoint := fmt.Sprintf("%s/%s?%s=%d", EnvironmentDevelopment.apiURL, modelGoal, fieldID, goal.ID)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodDelete, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var deleted bool
		var response *StandardResponse
		deleted, response, err = client.DeleteGoal(goal.ID)
		assert.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, false, deleted)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// ExampleClient_DeleteGoal example using DeleteGoal()
//
// See more examples in /examples/
func ExampleClient_DeleteGoal() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseGoal := newTestGoal()
	_ = mockResponseData(
		http.MethodDelete,
		fmt.Sprintf("%s/%s?%s=%d", EnvironmentDevelopment.apiURL, modelGoal, fieldID, responseGoal.ID),
		http.StatusOK,
		nil,
	)

	// Delete goal (using mocking response)
	var deleted bool
	if deleted, _, err = client.DeleteGoal(responseGoal.ID); err != nil {
		fmt.Printf("error deleting goal: " + err.Error())
		return
	}
	fmt.Printf("goal deleted: %t", deleted)
	// Output:goal deleted: true
}

// BenchmarkClient_DeleteGoal benchmarks the method DeleteGoal()
func BenchmarkClient_DeleteGoal(b *testing.B) {
	client, _ := newTestClient()
	goal := newTestGoal()
	_ = mockResponseData(
		http.MethodDelete,
		fmt.Sprintf("%s/%s?%s=%d", EnvironmentDevelopment.apiURL, modelGoal, fieldID, goal.ID),
		http.StatusOK,
		nil,
	)
	for i := 0; i < b.N; i++ {
		_, _, _ = client.DeleteGoal(goal.ID)
	}
}
