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
	Token    string `json:"token"`
}

// Pagination -
type Pagination struct {
	Page    uint  `json:"page"`
	Next    *uint `json:"next"`
	Prev    *uint `json:"prev"`
	Pages   uint  `json:"pages"`
	PerPage uint  `json:"per_page"`
	Total   uint  `json:"total"`
}

// Meta -
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// ApiResponse -
type ApiResponse struct {
	Meta    *Meta            `json:"meta"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

// NewClient -
func NewClient(host, username, password *string, userAgent *string) (*Client, error) {
	c := Client{
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		HostUrl: strings.TrimRight(*host, "/"),
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
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", client.Auth.Token))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", res.StatusCode)
	}

	result := new(ApiResponse)
	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("success: %v", result.Success)
	}

	return result.Data, err
}

func (client *Client) DoRequest(req *http.Request) (io.ReadCloser, error) {
	err := client.setNonce("/settings")
	if err != nil {
		return nil, err
	}
	req.Header.Set("CSRF-Token", client.Auth.Nonce)

	res, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, res.Body)
	}

	return res.Body, err
}
