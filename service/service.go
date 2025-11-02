package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Prasang-money/searchSvc/cache"
	"github.com/Prasang-money/searchSvc/models"
)

// make baseURL a variable so tests can override it
var baseURL = "https://restcountries.com/v3.1/name/"

type ServiceInterface interface {
	SearchCountries(name string) (*models.CountryMetadata, error)
}
type Service struct {
	cache *cache.Cache
}

func NewService(cache *cache.Cache) *Service {
	return &Service{
		cache: cache,
	}
}

func (s *Service) SearchCountries(name string) (*models.CountryMetadata, error) {
	// Check if results are in cache
	if cachedResults, found := s.cache.Get(name); found {
		return &cachedResults, nil

	}

	// If not found in cache, fetch from REST API
	url := baseURL + name
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch country data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	// fmt.Println(string(body))

	var countries []models.Country
	// Unmarshal JSON into struct
	if err := json.Unmarshal(body, &countries); err != nil {
		fmt.Println("Error:", err)
	}

	for _, country := range countries {
		if country.Name.Common == name {

			countryMetaData := models.CountryMetadata{
				Name:       country.Name.Common,
				Population: country.Population,
			}

			if len(country.Capital) > 0 {
				countryMetaData.Capital = country.Capital[0]
			}
			for _, curr := range country.Currencies {
				countryMetaData.Currency = curr.Symbol
				break
			}
			// Store results in cache before returning
			s.cache.Set(name, &countryMetaData)
			//fmt.Println(countryMetaData)
			return &countryMetaData, nil
		}
	}

	// Store results in cache before returning

	return &models.CountryMetadata{}, nil
}
