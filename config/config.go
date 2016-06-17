package config

import (
	"encoding/json"
	"fmt"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's config.
func List(c *deis.Client, app string) (api.Config, error) {
	u := fmt.Sprintf("/v2/apps/%s/config/", app)

	body, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return api.Config{}, err
	}

	config := api.Config{}
	if err = json.Unmarshal([]byte(body), &config); err != nil {
		return api.Config{}, err
	}

	return config, nil
}

// Set sets an app's config variables.
func Set(c *deis.Client, app string, config api.Config) (api.Config, error) {
	body, err := json.Marshal(config)

	if err != nil {
		return api.Config{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/config/", app)

	resBody, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return api.Config{}, err
	}

	newConfig := api.Config{}
	if err = json.Unmarshal([]byte(resBody), &newConfig); err != nil {
		return api.Config{}, err
	}

	return newConfig, nil
}
