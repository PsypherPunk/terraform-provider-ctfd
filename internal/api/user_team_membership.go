package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserTeamMembership - fields as returned from the CTFd API
type UserTeamMembership struct {
	UserId uint `json:"user_id"`
	TeamId uint `json:"team_id"`
}

// CreateUserTeamMembership - create a new userTeamMembership
func (client *Client) CreateUserTeamMembership(teamId uint, userId uint) (*UserTeamMembership, error) {
	userTeamMembership := map[string]interface{}{
		"user_id": userId,
	}
	rb, err := json.Marshal(userTeamMembership)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/teams/%d/members", client.HostUrl, teamId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	_, err = client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	newUserTeamMembership := UserTeamMembership{
		TeamId: teamId,
		UserId: userId,
	}

	return &newUserTeamMembership, nil
}

// DeleteUserTeamMembership - remove an existing userTeamMembership
func (client *Client) DeleteUserTeamMembership(teamId uint, userId uint) error {
	userTeamMembership := map[string]interface{}{
		"user_id": userId,
	}
	rb, err := json.Marshal(userTeamMembership)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/teams/%d/members", client.HostUrl, teamId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = client.DoApiRequest(req)
	if err != nil {
		return err
	}

	return nil
}
