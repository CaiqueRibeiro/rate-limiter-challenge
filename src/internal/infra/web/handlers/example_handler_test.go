package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleHandler_Get(t *testing.T) {
	handler := NewExampleHandler()

	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	rr := httptest.NewRecorder()

	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "expected status 200")

	var response ExampleDto
	expected := ExampleDto{
		Message: "Ok",
		Status:  200,
	}

	err := json.NewDecoder(rr.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Equal(t, expected, response)
}
