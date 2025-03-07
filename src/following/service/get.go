package following_service

import (
	"encoding/json"
	"fmt"
	"net/http"

	log_service "github.com/Rfluid/insta-tools/src/log/service"
	"github.com/pterm/pterm"
)

// Get makes a request to Instagram's API and returns the result as a map[string]interface{}
func Get(
	userID string,
	cookies map[string]string,
	count int,
	maxID string,
) (map[string]interface{}, error) {
	// Construct the request URL with user ID
	url := fmt.Sprintf("https://www.instagram.com/api/v1/friendships/%s/following/", userID)

	// Build query parameters
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for key, value := range cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	// Add query parameters
	query := req.URL.Query()
	query.Add("count", fmt.Sprintf("%d", count))
	if maxID != "" {
		query.Add("max_id", maxID)
	}
	req.URL.RawQuery = query.Encode()

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		log_service.LogConditionally(
			pterm.DefaultLogger.Error,
			fmt.Sprintf("Error fetching following. API status code is %v", resp.StatusCode),
		)

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		return result, fmt.Errorf("bad status code (%v) in API response", resp.StatusCode)
	}

	// Parse the JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
