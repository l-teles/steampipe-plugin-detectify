package detectify

import (
	"context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
	"encoding/json"
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
	var baseUrl, token, secret, tokenv3 string

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

	tokenv3 = os.Getenv("DETECTIFY_API_TOKEN_V3")
	if detectifyConfig.Tokenv3 != nil {
		tokenv3 = *detectifyConfig.Tokenv3
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
	
	if tokenv3 == "" {
		return "", errors.New("'tokenv3' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
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

func paginatedResponse(ctx context.Context, d *plugin.QueryData, endpoint string) ([]string, error) {
	var paginatedResponse []interface{}
	plugin.Logger(ctx).Info("Getting Detectify findings...")

    pageSize := 100
    markerID := ""

	// Iteration for Pagination
	for {
		params := map[string]string{
			"marker":    markerID,
			"page_size": fmt.Sprintf("%d", pageSize),
		}
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