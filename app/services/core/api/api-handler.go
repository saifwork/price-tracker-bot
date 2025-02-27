// api/api.go
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ApiClient interface {
	GetAPIRequest(ctx context.Context, apiURL string, responseData any) error
	PostAPIRequest(ctx context.Context, apiURL string, requestData interface{}, responseData interface{}) error
	PatchAPIRequest(ctx context.Context, apiURL string, requestData interface{}, responseData interface{}) error
}

type Client struct{}

func (c *Client) GetAPIRequest(ctx context.Context, apiURL string, responseData any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	if err := DecodeAPIResponse(resp.Body, responseData); err != nil {
		return err
	}

	return nil
}

func (c *Client) PostAPIRequest(ctx context.Context, apiURL string, requestData interface{}, responseData interface{}) error {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("failed to parse API URL: %v", err)
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to encode request payload: %v", err)
	}

	resp, err := http.Post(parsedURL.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	if err := DecodeAPIResponse(resp.Body, responseData); err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteAPIRequest(ctx context.Context, apiURL string, requestData interface{}, responseData interface{}) error {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("failed to parse API URL: %v", err)
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to encode request payload: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", parsedURL.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	if err := DecodeAPIResponse(resp.Body, responseData); err != nil {
		return err
	}

	return nil
}


func (c *Client) PatchAPIRequest(ctx context.Context, apiURL string, requestData interface{}, responseData interface{}) error {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("failed to parse API URL: %v", err)
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to encode request payload: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", parsedURL.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create PATCH request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make PATCH API request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println("Response Body:", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, responseData); err != nil {
		return fmt.Errorf("failed to decode response JSON: %v, response: %s", err, string(bodyBytes))
	}

	return nil
}

func DecodeAPIResponse(body io.Reader, response interface{}) error {
	if err := json.NewDecoder(body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode API response: %v", err)
	}
	return nil
}

func ConvertDataToStruct(data any, target interface{}) error {
	if dataMap, ok := data.(map[string]any); ok {
		dataBytes, err := json.Marshal(dataMap)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %v", err)
		}

		if err := json.Unmarshal(dataBytes, target); err != nil {
			return fmt.Errorf("failed to unmarshal data into struct: %v", err)
		}

		return nil
	}

	return fmt.Errorf("data is not in the expected format")
}
