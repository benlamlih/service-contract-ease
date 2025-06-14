package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"contract_ease/internal/database/mocks"
)

func TestHelloWorldHandler(t *testing.T) {
	t.Parallel()
	mockDB := mocks.NewService(t)
	mockDB.On("Health", mock.Anything).Return(map[string]string{"status": "up"})

	s := NewServer(mockDB, 8080)
	r := gin.New()
	r.GET("/api/health", s.healthHandler)
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "{\"status\":\"up\"}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
