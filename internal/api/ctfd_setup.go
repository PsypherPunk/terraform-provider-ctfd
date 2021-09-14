package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type configData struct {
	Id    uint   `json:"id"`
	Value string `json:"value"`
	Key   string `json:"key"`
}

type EmailConfig struct {
	Username    string `json:"mail_username"`
	Password    string `json:"mail_password"`
	FromAddress string `json:"mailfrom_addr"`
	Server      string `json:"mail_server"`
	Port        int    `json:"mail_port"`
	UseAuth     bool   `json:"mail_useauth"`
	UseTls      bool   `json:"mail_tls"`
	UseSsl      bool   `json:"mail_ssl"`
}

// CtfdSetup - `ctfd_setup` resource
type CtfdSetup struct {
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	AdminEmail        string       `json:"admin_email"`
	ConfigurationPath string       `json:"configuration_path"`
	Email             *EmailConfig `json:"email"`
}

// GetCtfdSetup - Retrieve details of the CTFd setup
func (client *Client) GetCtfdSetup() (*CtfdSetup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/configs", client.HostUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	ctfdSetup := new(CtfdSetup)
	config := new([]configData)
	err = json.Unmarshal(*body, &config)
	if err != nil {
		return nil, err
	}
	for _, value := range *config {
		if value.Key == "ctf_name" {
			ctfdSetup.Name = value.Value
		}
		if value.Key == "ctf_description" {
			ctfdSetup.Description = value.Value
		}
	}

	// TODO: get email config.

	return ctfdSetup, nil
}

// doSetup - perform the initial setup for CTFd
func doSetup(client *Client, setup CtfdSetup) error {
	err := client.setNonce("/setup")
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Set("nonce", client.Auth.Nonce)
	form.Set("ctf_name", setup.Name)
	form.Set("ctf_description", setup.Description)
	form.Set("name", client.Auth.Username)
	form.Set("user_mode", "teams")
	form.Set("email", setup.AdminEmail)
	form.Set("password", client.Auth.Password)

	res, err := client.HttpClient.PostForm(fmt.Sprintf("%s/setup", client.HostUrl), form)
	if err != nil {
		return err
	}
	if res.StatusCode != 302 {
		msg, err := GetErrorFromHtml(*res)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: unable to setup", *msg)
	}

	return nil
}

// importConfiguration - upload challenges file, destroying setup
func importConfiguration(client *Client, setup CtfdSetup) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	err := client.setNonce("/admin/config")
	if err != nil {
		return err
	}
	err = writer.WriteField("nonce", client.Auth.Nonce)
	if err != nil {
		return err
	}

	part, err := writer.CreateFormFile("backup", filepath.Base(setup.ConfigurationPath))
	if err != nil {
		return err
	}

	file, err := os.Open(setup.ConfigurationPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/admin/import", client.HostUrl), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 302 {
		msg, err := GetErrorFromHtml(*res)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: unable to import config", *msg)
	}

	return nil
}

// setupEmail - configure email services
func setupEmail(client Client, emailConfig EmailConfig) error {
	rb, err := json.Marshal(emailConfig)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/configs", client.HostUrl), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = client.DoApiRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// CreateCtfdSetup - setup a new CTFd instance
func (client *Client) CreateCtfdSetup(setup CtfdSetup) error {
	// do initial setup
	err := doSetup(client, setup)
	if err != nil {
		return err
	}

	// import configuration file
	err = importConfiguration(client, setup)
	if err != nil {
		return err
	}

	// repeat initial setup
	err = doSetup(client, setup)
	if err != nil {
		return err
	}

	err = client.CheckSetup()
	if err != nil {
		return err
	}

	err = client.SignIn()
	if err != nil {
		return err
	}

	token, err := client.GetOrCreateToken()
	if err != nil {
		return err
	}
	client.Auth.Token = token.Value

	if setup.Email != nil {
		err := setupEmail(*client, *setup.Email)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteCtfdSetup - remove CTFd setup
func (client *Client) DeleteCtfdSetup() error {
	err := client.setNonce("/admin/config")
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Set("nonce", client.Auth.Nonce)
	form.Set("accounts", "y")
	form.Set("submissions", "y")
	form.Set("challenges", "y")
	form.Set("pages", "y")
	form.Set("notifications", "y")

	res, err := client.HttpClient.PostForm(fmt.Sprintf("%s/admin/reset", client.HostUrl), form)
	if err != nil {
		return err
	}
	if res.StatusCode != 302 {
		msg, err := GetErrorFromHtml(*res)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: unable to reset", *msg)
	}

	return nil
}
