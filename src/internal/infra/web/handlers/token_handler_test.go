package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestTokenHandler_Create_Success(t *testing.T) {
	token := "dummy_token"
	db, clientMock := redismock.NewClientMock()
	handler := NewTokenHandler(db)

	clientMock.ExpectSet("token_max_req:"+token, 10, time.Duration(0)).SetVal("OK")

	body := `{"token":"dummy_token","max_requests":10}`
	req := httptest.NewRequest(http.MethodPost, "/create-token", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, `{"message":"Token registered"}`, rr.Body.String())
	assert.NoError(t, clientMock.ExpectationsWereMet())
}

func TestTokenHandlerCreateBadRequestInvalidBody(t *testing.T) {
	db, _ := redismock.NewClientMock()
	handler := NewTokenHandler(db)

	body := `{"token":"","max_requests":0}`
	req := httptest.NewRequest(http.MethodPost, "/create-token", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"message":"Invalid body"}`, rr.Body.String())
}

func TestTokenHandlerCreateBadRequestBodyDecodeError(t *testing.T) {
	db, _ := redismock.NewClientMock()
	handler := NewTokenHandler(db)

	body := `{"token":`
	req := httptest.NewRequest(http.MethodPost, "/create-token", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"message":"Unable to read the body"}`, rr.Body.String())
}
