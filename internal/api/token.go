package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Token struct {
	Id         uint   `json:"id"`
	Type       string `json:"type"`
	UserId     uint   `json:"user_id"`
	Created    string `json:"created"`
	Expiration string `json:"expiration"`
	Value      string `json:"value"`
}

// GetTokens - get token objects in bulk
func (client *Client) GetTokens() (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tokens", client.HostUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	tokens := make([]map[string]interface{}, 0)
	err = json.Unmarshal(*body, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// CreateToken - create a token object
func (client *Client) CreateToken() (newToken *Token, err error) {
	emptyRequest := []byte("{}")
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tokens", client.HostUrl), bytes.NewBuffer(emptyRequest))
	if err != nil {
		return nil, err
	}

	err = client.setNonce("/settings")
	if err != nil {
		return nil, err
	}
	req.Header.Set("CSRF-Token", client.Auth.Nonce)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	apiResponse := new(ApiResponse)
	err = json.NewDecoder(res.Body).Decode(apiResponse)
	defer res.Body.Close()
	if res.StatusCode != 200 || !apiResponse.Success {
		return nil, errors.New(fmt.Sprintf("token generation failed: %s", apiResponse.Message))
	}

	err = json.Unmarshal(*apiResponse.Data, &newToken)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
