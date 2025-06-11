package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIClient реализует запросы к внешним API
type APIClient struct {
	agifyURL       string
	genderizeURL   string
	nationalizeURL string
	httpClient     *http.Client
}

// NewAPIClient создаёт клиент для работы с API
func NewAPIClient(agifyURL, genderizeURL, nationalizeURL string) *APIClient {
	return &APIClient{
		agifyURL:       agifyURL,
		genderizeURL:   genderizeURL,
		nationalizeURL: nationalizeURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // Таймаут на запрос
		},
	}
}

// GetAge возвращает предполагаемый возраст по имени
func (c *APIClient) GetAge(ctx context.Context, name string) (int, error) {
	url := fmt.Sprintf("%s?name=%s", c.agifyURL, url.QueryEscape(name))

	resp, err := c.doRequest(ctx, url)
	if err != nil {
		return 0, err
	}

	var result struct {
		Age int `json:"age"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to parse age response: %w", err)
	}

	return result.Age, nil
}

// GetGender возвращает предполагаемый пол по имени
func (c *APIClient) GetGender(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s?name=%s", c.genderizeURL, url.QueryEscape(name))

	resp, err := c.doRequest(ctx, url)
	if err != nil {
		return "", err
	}

	var result struct {
		Gender string `json:"gender"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse gender response: %w", err)
	}

	return strings.ToLower(result.Gender), nil // "male" вместо "Male"
}

// GetNationality возвращает предполагаемую национальность по имени
func (c *APIClient) GetNationality(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s?name=%s", c.nationalizeURL, url.QueryEscape(name))

	resp, err := c.doRequest(ctx, url)
	if err != nil {
		return "", err
	}

	var result struct {
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse nationality response: %w", err)
	}

	if len(result.Country) == 0 {
		return "", fmt.Errorf("no nationality data")
	}

	// 🎯 Лямбда-функция для взвешенного выбора страны
	randCountry := func(countries []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	}) string {
		var total float64
		for _, c := range countries {
			total += c.Probability
		}
		r := rand.Float64() * total

		var acc float64
		for _, c := range countries {
			acc += c.Probability
			if r < acc {
				return c.CountryID
			}
		}
		return countries[len(countries)-1].CountryID // На случай, если что-то пошло не так
	}

	return randCountry(result.Country), nil
}

// doRequest общий метод для остальных GET HTTP-запросов API-клиента
func (c *APIClient) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}
