package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Prasang-money/searchSvc/cache"
	"github.com/Prasang-money/searchSvc/models"
)

func TestSearchCountries_CacheHit(t *testing.T) {
	c := cache.NewCache(10)
	svc := NewService(c)

	expected := &models.CountryMetadata{
		Name:       "CachedLand",
		Population: 42,
		Capital:    "CacheCity",
		Currency:   "CCH",
	}

	// Prime the cache
	c.Set("CachedLand", expected)

	res, err := svc.SearchCountries("CachedLand")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatalf("expected non-nil response")
	}
	if res.Name != expected.Name {
		t.Fatalf("expected name %s, got %s", expected.Name, res.Name)
	}
}

func TestSearchCountries_FetchFromAPI(t *testing.T) {
	// Start a test server that returns a JSON similar to the real API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a slice with one country whose Name.Common matches the requested name
		country := models.Country{
			Name:       models.Name{Common: "Testland"},
			Population: 12345,
			Capital:    []string{"T-City"},
			Currencies: map[string]models.Currencies{"TST": {Name: "TestCurrency", Symbol: "T$"}},
		}
		arr := []models.Country{country}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(arr)
	}))
	defer ts.Close()

	// Override baseURL to point to the test server
	orig := baseURL
	baseURL = ts.URL + "/"
	defer func() { baseURL = orig }()

	c := cache.NewCache(10)
	svc := NewService(c)

	res, err := svc.SearchCountries("Testland")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "Testland" {
		t.Fatalf("expected name Testland, got %s", res.Name)
	}
	if res.Capital != "T-City" {
		t.Fatalf("expected capital T-City, got %s", res.Capital)
	}
	if res.Currency != "T$" {
		t.Fatalf("expected currency T$, got %s", res.Currency)
	}

	// also ensure it was stored in cache
	cached, found := c.Get("Testland")
	if !found {
		t.Fatalf("expected cached entry for Testland")
	}
	if cached.Name != "Testland" {
		t.Fatalf("cached value mismatch: expected Testland, got %s", cached.Name)
	}
}

func TestSearchCountries_APIError(t *testing.T) {
	// Server that returns a 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("boom"))
	}))
	defer ts.Close()

	orig := baseURL
	baseURL = ts.URL + "/"
	defer func() { baseURL = orig }()

	c := cache.NewCache(10)
	svc := NewService(c)

	_, err := svc.SearchCountries("Anything")
	if err == nil {
		t.Fatalf("expected error from API non-200 status, got nil")
	}
	if !strings.Contains(err.Error(), "API returned status code") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
