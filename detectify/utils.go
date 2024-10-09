package detectify

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func signRequest(detectifyAPISecret, detectifyAPIToken, method, relativeURL, timestamp, body string) string {
	key, _ := base64.StdEncoding.DecodeString(detectifyAPISecret)
	value := fmt.Sprintf("%s;%s;%s;%s;%s", method, relativeURL, detectifyAPIToken, timestamp, body)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(value))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}

func connect(ctx context.Context, d *plugin.QueryData, endpoint string, params map[string]string) (string, error) {
	var baseUrl, token, secret, token_v3 string

	// Prefer config options given in Steampipe
	detectifyConfig := GetConfig(d.Connection)

	baseUrl = os.Getenv("DETECTIFY_URL")
	if detectifyConfig.BaseUrl != nil {
		baseUrl = *detectifyConfig.BaseUrl
	}

	token = os.Getenv("DETECTIFY_API_TOKEN")
	if detectifyConfig.Token != nil {
		token = *detectifyConfig.Token
	}

	token_v3 = os.Getenv("DETECTIFY_API_TOKEN_V3")
	if detectifyConfig.Token_v3 != nil {
		token_v3 = *detectifyConfig.Token_v3
	}

	secret = os.Getenv("DETECTIFY_API_SECRET")
	if detectifyConfig.Secret != nil {
		secret = *detectifyConfig.Secret
	}

	if baseUrl == "" {
		return "", errors.New("'baseUrl' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	if token == "" {
		return "", errors.New("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	if secret == "" {
		return "", errors.New("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	if token_v3 == "" {
		return "", errors.New("'token_v3' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	signature := signRequest(secret, token, "GET", endpoint, timestamp, "")

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", baseUrl+endpoint, nil)
	if err != nil {
		plugin.Logger(ctx).Error("Failed to create request: %v", err)
		return fmt.Sprintf("Failed to create request: %v", err), err
	}

	// Set the necessary request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Detectify-Key", token)
	req.Header.Set("X-Detectify-Signature", signature)
	req.Header.Set("X-Detectify-Timestamp", timestamp)

	// Set query parameters for pagination
	queryParams := req.URL.Query()
	for key, value := range params {
		queryParams.Add(key, value)
	}
	req.URL.RawQuery = queryParams.Encode()

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		plugin.Logger(ctx).Error("Failed to execute request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != 200 {
		plugin.Logger(ctx).Error("utils.connect", "status_code_error", resp.Status)
		return "", fmt.Errorf("API returned HTTP status %s", resp.Status)
	}

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.Logger(ctx).Error("Failed to read response body: %v", err)
		return "", err
	}

	return string(body), nil
}

// paginatedResponse fetches paginated responses from the given endpoint.
func paginatedResponse(ctx context.Context, d *plugin.QueryData, endpoint string, optionalParams ...map[string]string) ([]string, error) {
	var paginatedResponse []interface{}
	plugin.Logger(ctx).Info("Getting Detectify findings...")

	pageSize := 100
	markerID := ""

	// Set default params
	params := map[string]string{
		"marker":    markerID,
		"page_size": fmt.Sprintf("%d", pageSize),
	}

	// Override default params with optionalParams if provided
	if len(optionalParams) > 0 {
		for key, value := range optionalParams[0] {
			params[key] = value
		}
	}

	// Iteration for Pagination
	for {
		// Update marker in params for each iteration
		params["marker"] = markerID

		findingsStr, err := connect(ctx, d, endpoint, params)
		if err != nil {
			plugin.Logger(ctx).Error("utils.paginatedResponse", "connection_error", err)
			return nil, err
		}

		var findings map[string]interface{}
		if err := json.Unmarshal([]byte(findingsStr), &findings); err != nil {
			plugin.Logger(ctx).Error("Failed to parse response body: %v", err)
			return nil, err
		}

		paginatedResponse = append(paginatedResponse, findingsStr)

		if !findings["has_more"].(bool) {
			break
		}
		markerID = findings["next_marker"].(string)
	}

	// Convert paginatedResponse to []string
	var result []string
	for _, v := range paginatedResponse {
		result = append(result, v.(string))
	}

	return result, nil
}

func connectV3(ctx context.Context, d *plugin.QueryData, endpoint string, params map[string]string) (string, error) {
	var baseUrl, token_v3 string

	// Prefer config options given in Steampipe
	detectifyConfig := GetConfig(d.Connection)

	baseUrl = os.Getenv("DETECTIFY_URL")
	if detectifyConfig.BaseUrl != nil {
		baseUrl = *detectifyConfig.BaseUrl
	}

	token_v3 = os.Getenv("DETECTIFY_API_TOKEN_V3")
	if detectifyConfig.Token_v3 != nil {
		token_v3 = *detectifyConfig.Token_v3
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", baseUrl+endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the necessary request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_v3)

	// Set query parameters
	queryParams := req.URL.Query()
	for key, value := range params {
		queryParams.Add(key, value)
	}
	req.URL.RawQuery = queryParams.Encode()

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		plugin.Logger(ctx).Error("Failed to create request: %v", err)
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != 200 {
		plugin.Logger(ctx).Error("utils.connectV3 -> API returned HTTP status %s", resp.Status)
		return "", fmt.Errorf("utils.connectV3 -> API returned HTTP status %s", resp.Status)
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.Logger(ctx).Error("utils.connectV3 -> Failed to read response body: %v", err)
		return "", fmt.Errorf("utils.connectV3 -> failed to read response body: %v", err)
	}

	return string(body), nil
}

func paginatedResponseV3(ctx context.Context, d *plugin.QueryData, endpoint string) ([]string, error) {
	var paginatedResponse []interface{}
	plugin.Logger(ctx).Info("Getting Detectify findings...")

	cursor := ""

	// Iteration for Pagination
	for {
		params := map[string]string{
			"cursor": cursor,
		}
		findingsStr, err := connectV3(ctx, d, endpoint, params)
		if err != nil {
			plugin.Logger(ctx).Error("utils.paginatedResponseV3", "connection_error", err)
			return nil, err
		}

		var findings map[string]interface{}
		if err := json.Unmarshal([]byte(findingsStr), &findings); err != nil {
			plugin.Logger(ctx).Error("Failed to parse response body: %v", err)
			return nil, err
		}

		paginatedResponse = append(paginatedResponse, findingsStr)

		pagination, ok := findings["pagination"].(map[string]interface{})
		if !ok {
			break
		}

		nextURL, ok := pagination["next"].(string)
		if !ok || nextURL == "" {
			break
		}

		// Extract cursor from next URL
		u, err := url.Parse(nextURL)
		if err != nil {
			plugin.Logger(ctx).Error("utils.paginatedResponseV3", "invalid_next_url", err)
			return nil, err
		}

		cursor = u.Query().Get("cursor")
		if cursor == "" {
			break
		}
	}

	// Convert paginatedResponse to []string
	var result []string
	for _, v := range paginatedResponse {
		result = append(result, v.(string))
	}

	return result, nil
}
