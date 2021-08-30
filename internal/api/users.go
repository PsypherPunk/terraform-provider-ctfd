package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// NewUser - fields required when creating a new user
type NewUser struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Website     string `json:"website"`
	Affiliation string `json:"affiliation"`
	Country     string `json:"country"`
	Hidden      bool   `json:"hidden"`
	Banned      bool   `json:"banned"`
	Type        string `json:"type"`
	Verified    bool   `json:"verified"`
}

// UpdateUser - fields required when updating a user (password cannot be updated.)
type UpdateUser struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Affiliation string `json:"affiliation"`
	Country     string `json:"country"`
	Hidden      bool   `json:"hidden"`
	Banned      bool   `json:"banned"`
	Type        string `json:"type"`
	Verified    bool   `json:"verified"`
}

// User - fields as returned from the CTFd API
type User struct {
	Id          uint     `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Website     string   `json:"website"`
	Affiliation string   `json:"affiliation"`
	Country     string   `json:"country"`
	Bracket     string   `json:"bracket"`
	Secret      string   `json:"secret"`
	OauthId     string   `json:"oauth_id"`
	Fields      []string `json:"fields"`
	Type        string   `json:"type"`
	TeamId      uint     `json:"team_id"`
	Verified    bool     `json:"verified"`
	Hidden      bool     `json:"hidden"`
	Banned      bool     `json:"banned"`
}

// GetUsers - Returns list of users
func (client *Client) GetUsers() (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users", client.HostUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	users := make([]map[string]interface{}, 0)
	err = json.Unmarshal(*body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser - Returns details of a user
func (client *Client) GetUser(id uint) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/%d", client.HostUrl, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = json.Unmarshal(*body, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser - create a new user
func (client *Client) CreateUser(user NewUser) (*User, error) {
	rb, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/users", client.HostUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	newUser := new(User)
	err = json.Unmarshal(*body, &newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// UpdateUser - updated an existing user
func (client *Client) UpdateUser(id uint, user NewUser) (*User, error) {
	rb, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/users/%d", client.HostUrl, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	updatedUser := new(User)
	err = json.Unmarshal(*body, &updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser - remove an existing user
func (client *Client) DeleteUser(id uint) error {
	emptyRequest := []byte("{}")
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/users/%d", client.HostUrl, id), bytes.NewBuffer(emptyRequest))
	if err != nil {
		return err
	}

	_, err = client.DoApiRequest(req)
	if err != nil {
		return err
	}

	return nil
}
