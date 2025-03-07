package user_service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Get fetches Instagram user profile info, including the user ID.
func Get(username string, cookies map[string]string) (map[string]interface{}, error) {
	// Construct the request URL with the username
	url := fmt.Sprintf("https://www.instagram.com/api/v1/users/web_profile_info/?username=%s", username)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers to request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for key, value := range cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	// Parse the JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
