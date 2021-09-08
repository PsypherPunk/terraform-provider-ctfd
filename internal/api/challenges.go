package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetChallenges - Returns list of challenges
func (client *Client) GetChallenges() (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/challenges", client.HostUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.DoApiRequest(req)
	if err != nil {
		return nil, err
	}

	challenges := make([]map[string]interface{}, 0)
	err = json.Unmarshal(*body, &challenges)
	if err != nil {
		return nil, err
	}

	return challenges, nil
}
