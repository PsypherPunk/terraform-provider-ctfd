package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NewTeam struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Website     string `json:"website"`
	Affiliation string `json:"affiliation"`
	Country     string `json:"country"`
	Hidden      bool   `json:"hidden"`
	Banned      bool   `json:"banned"`
}

type Team struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Website     string   `json:"website"`
	Affiliation string   `json:"affiliation"`
	Country     string   `json:"country"`
	Hidden      bool     `json:"hidden"`
	Banned      bool     `json:"banned"`
	CaptainId   *uint    `json:"captain_id"`
	Bracket     string   `json:"bracket"`
	Id          uint     `json:"id"`
	Secret      string   `json:"secret"`
	OauthId     string   `json:"oauth_id"`
	Members     []uint   `json:"members"`
	Created     string   `json:"created"`
	Fields      []string `json:"fields"`
}

// GetTeams - Returns list of teams
func (client *Client) GetTeams() (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/teams", client.HostUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	teams := make([]map[string]interface{}, 0)
	err = json.Unmarshal(*body, &teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

// GetTeam - Returns details of a team
func (client *Client) GetTeam(id uint) (*Team, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/teams/%d", client.HostUrl, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	team := new(Team)
	err = json.Unmarshal(*body, &team)
	if err != nil {
		return nil, err
	}

	return team, nil
}

// CreateTeam - create a new team
func (client *Client) CreateTeam(team NewTeam) (*Team, error) {
	rb, err := json.Marshal(team)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/teams", client.HostUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	newTeam := new(Team)
	err = json.Unmarshal(*body, &newTeam)
	if err != nil {
		return nil, err
	}

	return newTeam, nil
}

// UpdateTeam - updated an existing team
func (client *Client) UpdateTeam(id uint, team NewTeam) (*Team, error) {
	rb, err := json.Marshal(team)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/teams/%d", client.HostUrl, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	updatedTeam := new(Team)
	err = json.Unmarshal(*body, &updatedTeam)
	if err != nil {
		return nil, err
	}

	return updatedTeam, nil
}

// DeleteTeam - remove an existing team
func (client *Client) DeleteTeam(id uint) error {
	emptyRequest := []byte("{}")
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/teams/%d", client.HostUrl, id), bytes.NewBuffer(emptyRequest))
	if err != nil {
		return err
	}

	_, err = client.DoApiRequest(req)
	if err != nil {
		return err
	}

	return nil
}
