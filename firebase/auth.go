package firebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInResponse struct {
	LocalId        string `json:"localId"`
	Email          string `json:"email"`
	DisplayName    string `json:"displayName"`
	IdToken        string `json:"idToken"`
	ProfilePicture string `json:"profilePicture"`
	RefreshToken   string `json:"refreshToken"`
	ExpiresIn      string `json:"expiresIn"`
}

func (c *Client) SignIn(email, password string) (SignInResponse, error) {
	if len(c.apiRoute) == 0 {
		return SignInResponse{}, errors.New("empty api key")
	}

	reqBody, err := json.Marshal(credentials{Email: email, Password: password})
	if err != nil {
		return SignInResponse{}, err
	}

	resp, err := c.http.Post(c.apiRoute, "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		return SignInResponse{}, err
	}

	var profile SignInResponse
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return SignInResponse{}, err
	}

	return profile, err
}
