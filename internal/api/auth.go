package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var nonceRegex = regexp.MustCompile("'csrfNonce': \"([a-z0-9]+)\",")

// GetNonce - For use in subsequent requests
func (client *Client) setNonce() error {
	resp, err := client.HttpClient.Get(client.HostUrl)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	parts := nonceRegex.FindSubmatch(body)
	nonce := parts[1]
	client.Auth.Nonce = string(nonce)

	return nil
}

// SignIn - Get a new token for user
func (client *Client) SignIn() error {
	err := client.setNonce()
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Set("nonce", client.Auth.Nonce)
	form.Set("name", client.Auth.Username)
	form.Set("password", client.Auth.Password)

	_, err = client.HttpClient.PostForm(fmt.Sprintf("%s/login", client.HostUrl), form)
	if err != nil {
		return err
	}

	return nil
}

// SignOut - Revoke the token for a user
func (client *Client) SignOut() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/logout", client.HostUrl), nil)
	if err != nil {
		return err
	}

	_, err = client.DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
