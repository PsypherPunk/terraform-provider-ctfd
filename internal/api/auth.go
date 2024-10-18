package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var nonceRegex = regexp.MustCompile("'csrfNonce': \"([a-z0-9]+)\",")

// GetNonce - For use in subsequent requests
func (client *Client) setNonce(path string) error {
	path = strings.TrimLeft(path, "/")
	res, err := client.HttpClient.Get(fmt.Sprintf("%s/%s", client.HostUrl, path))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	parts := nonceRegex.FindSubmatch(body)
	nonce := parts[1]
	client.Auth.Nonce = string(nonce)

	return nil
}

// CheckSetup - verify that a CTFd instance has been setup
func (client *Client) CheckSetup() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/login", client.HostUrl), nil)
	if err != nil {
		return err
	}

	res, err := client.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 302 {
		location := res.Header.Get("Location")
		if strings.HasSuffix(location, "/setup") {
			return errors.New("CTFd instance not setup; not currently supported")
		}
	}

	return nil
}

// SignIn - Authenticate the Admin. user
func (client *Client) SignIn() error {
	err := client.setNonce("/login")
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Set("nonce", client.Auth.Nonce)
	form.Set("name", client.Auth.Username)
	form.Set("password", client.Auth.Password)

	res, err := client.HttpClient.PostForm(fmt.Sprintf("%s/login", client.HostUrl), form)
	if err != nil {
		return err
	}
	if res.StatusCode != 302 {
		location := res.Header.Get("Location")
		if !strings.HasSuffix(location, "/challenges") {
			return errors.New("unable to sign in; not redirected to /challenges")
		}
	}

	return nil
}

// SignOut - Request the logout page
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
