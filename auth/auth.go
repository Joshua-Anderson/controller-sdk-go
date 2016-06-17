package auth

import (
	"encoding/json"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// Register a new user with the controller.
func Register(c *deis.Client, username, password, email string) error {
	user := api.AuthRegisterRequest{Username: username, Password: password, Email: email}
	body, err := json.Marshal(user)

	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", "/v2/auth/register/", body)
	return err
}

// Login to the controller and get a token
func Login(c *deis.Client, username, password string) (string, error) {
	user := api.AuthLoginRequest{Username: username, Password: password}
	reqBody, err := json.Marshal(user)

	if err != nil {
		return "", err
	}

	body, err := c.BasicRequest("POST", "/v2/auth/login/", reqBody)

	if err != nil {
		return "", err
	}

	token := api.AuthLoginResponse{}
	if err = json.Unmarshal([]byte(body), &token); err != nil {
		return "", err
	}

	return token.Token, nil
}

// Delete deletes a user.
func Delete(c *deis.Client, username string) error {
	var body []byte
	var err error

	if username != "" {
		req := api.AuthCancelRequest{Username: username}
		body, err = json.Marshal(req)

		if err != nil {
			return err
		}
	}

	_, err = c.BasicRequest("DELETE", "/v2/auth/cancel/", body)
	return err
}

// Regenerate user's auth tokens.
func Regenerate(c *deis.Client, username string, all bool) (string, error) {
	var reqBody []byte
	var err error

	if all == true {
		reqBody, err = json.Marshal(api.AuthRegenerateRequest{All: all})
	} else if username != "" {
		reqBody, err = json.Marshal(api.AuthRegenerateRequest{Name: username})
	}

	if err != nil {
		return "", err
	}

	body, err := c.BasicRequest("POST", "/v2/auth/tokens/", reqBody)

	if err != nil {
		return "", err
	}

	if all == true {
		return "", nil
	}

	token := api.AuthRegenerateResponse{}
	if err = json.Unmarshal([]byte(body), &token); err != nil {
		return "", err
	}

	return token.Token, nil
}

// Passwd changes a user's password.
func Passwd(c *deis.Client, username, password, newPassword string) error {
	req := api.AuthPasswdRequest{Password: password, NewPassword: newPassword}

	if username != "" {
		req.Username = username
	}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", "/v2/auth/passwd/", body)
	return err
}
