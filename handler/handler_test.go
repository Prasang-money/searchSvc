package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Prasang-money/searchSvc/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) SearchCountries(name string) (*models.CountryMetadata, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CountryMetadata), args.Error(1)
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", handler.HealthCheck())
	router.GET("/search", handler.SearchHandler())
	return router
}

func TestHealthCheck(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)
	router := setupTestRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response["status"])
}

func TestSearchHandler_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)
	router := setupTestRouter(handler)

	expectedResponse := &models.CountryMetadata{
		Name:       "United States",
		Population: 331002651,
		Capital:    "Washington, D.C.",
		Currency:   "USD",
	}

	// Set up mock expectation
	mockService.On("SearchCountries", "United States").Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?name=United States", nil)
	router.ServeHTTP(w, req)

	// Assert HTTP status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response models.CountryMetadata
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert response content
	assert.Equal(t, expectedResponse.Name, response.Name)
	assert.Equal(t, expectedResponse.Population, response.Population)
	assert.Equal(t, expectedResponse.Capital, response.Capital)
	assert.Equal(t, expectedResponse.Currency, response.Currency)

	// Verify that the mock was called as expected
	mockService.AssertExpectations(t)
}

func TestSearchHandler_Error(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)
	router := setupTestRouter(handler)

	expectedError := fmt.Errorf("service error")
	mockService.On("SearchCountries", "NonExistentCountry").Return(nil, expectedError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?name=NonExistentCountry", nil)
	router.ServeHTTP(w, req)

	// Assert HTTP status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert error message
	assert.Contains(t, w.Body.String(), expectedError.Error())

	// Verify that the mock was called as expected
	mockService.AssertExpectations(t)
}

func TestSearchHandler_EmptyQuery(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)
	router := setupTestRouter(handler)

	expectedResponse := &models.CountryMetadata{}
	mockService.On("SearchCountries", "").Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search", nil)
	router.ServeHTTP(w, req)

	// Assert HTTP status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response models.CountryMetadata
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert empty response
	assert.Equal(t, models.CountryMetadata{}, response)

	// Verify that the mock was called as expected
	mockService.AssertExpectations(t)
}

func TestNewHandler(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}
