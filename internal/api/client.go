package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

// Client -
type Client struct {
	HostUrl    string
	HttpClient *http.Client
	Auth       AuthStruct
	UserAgent  string
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
}

// ApiResponse -
type ApiResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

// NewClient -
func NewClient(host, username, password *string, userAgent *string) (*Client, error) {
	c := Client{
		HttpClient: &http.Client{Timeout: 5 * time.Second},
		HostUrl:    strings.TrimRight(*host, "/"),
		Auth: AuthStruct{
			Username: *username,
			Password: *password,
		},
		UserAgent: *userAgent,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c.HttpClient.Jar = jar

	return &c, nil
}

func (client *Client) DoApiRequest(req *http.Request) (*json.RawMessage, error) {
	err := client.SignIn()
	if err != nil {
		return nil, err
	}

	err = client.setNonce()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Csrf-Token", client.Auth.Nonce)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	result := new(ApiResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("success: %s", result.Success)
	}

	return result.Data, err
}

func (client *Client) DoRequest(req *http.Request) (io.ReadCloser, error) {
	err := client.setNonce()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Csrf-Token", client.Auth.Nonce)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, resp.Body)
	}

	return resp.Body, err
}
