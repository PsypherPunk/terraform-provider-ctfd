package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetChallenges - Returns list of challenges (no auth required)
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

/*
// GetChallengeIngredients - Returns list of challenge ingredients (no auth required)
func (client *Client) GetChallengeIngredients(challengeID string) ([]Ingredient, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/challenges/%s/ingredients", client.HostUrl, challengeID), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	ingredients := []Ingredient{}
	err = json.Unmarshal(body, &ingredients)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

// CreateChallenge - Create new challenge
func (client *Client) CreateChallenge(challenge Challenge) (*Challenge, error) {
	rb, err := json.Marshal(challenge)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/challenges", client.HostUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	newChallenge := Challenge{}
	err = json.Unmarshal(body, &newChallenge)
	if err != nil {
		return nil, err
	}

	return &newChallenge, nil
}

// CreateChallengeIngredient - Create new challenge ingredient
func (client *Client) CreateChallengeIngredient(challenge Challenge, ingredient Ingredient) (*Ingredient, error) {
	reqBody := struct {
		ChallengeID     int    `json:"challenge_id"`
		IngredientID int    `json:"ingredient_id"`
		Quantity     int    `json:"quantity"`
		Unit         string `json:"unit"`
	}{
		ChallengeID:     challenge.ID,
		IngredientID: ingredient.ID,
		Quantity:     ingredient.Quantity,
		Unit:         ingredient.Unit,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/challenges/%d/ingredients", client.HostUrl, challenge.ID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	newIngredient := Ingredient{}
	err = json.Unmarshal(body, &newIngredient)
	if err != nil {
		return nil, err
	}

	return &newIngredient, nil
}

*/
