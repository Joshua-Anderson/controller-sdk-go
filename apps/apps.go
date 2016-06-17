package apps

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

const workflowURLPrefix = "deis."

// ErrNoLogs is returned when logs are missing from an app.
var ErrNoLogs = errors.New(
	`There are currently no log messages. Please check the following things:
1) Logger and fluentd pods are running.
2) The application is writing logs to the logger component by checking that an entry in the ring buffer was created: kubectl logs <logger pod> --namespace=deis
3) Making sure that the container logs were mounted properly into the fluentd pod: kubectl exec <fluentd pod> --namespace=deis ls /var/log/containers`)

// List lists apps on a Deis controller.
func List(c *deis.Client, results int) ([]api.App, int, error) {
	body, count, err := c.LimitedRequest("/v2/apps/", results)

	if err != nil {
		return []api.App{}, -1, err
	}

	var apps []api.App
	if err = json.Unmarshal([]byte(body), &apps); err != nil {
		return []api.App{}, -1, err
	}

	for name, app := range apps {
		// Add in app URL based on controller hostname, port included
		app.URL = fmt.Sprintf("%s.%s", app.ID, strings.TrimPrefix(c.ControllerURL.Host, workflowURLPrefix))
		apps[name] = app
	}

	return apps, count, nil
}

// New creates a new app.
func New(c *deis.Client, id string) (api.App, error) {
	body := []byte{}

	var err error
	if id != "" {
		req := api.AppCreateRequest{ID: id}
		body, err = json.Marshal(req)

		if err != nil {
			return api.App{}, err
		}
	}

	resBody, err := c.BasicRequest("POST", "/v2/apps/", body)

	if err != nil {
		return api.App{}, err
	}

	app := api.App{}
	if err = json.Unmarshal([]byte(resBody), &app); err != nil {
		return api.App{}, err
	}

	// Add in app URL based on controller hostname, port included
	app.URL = fmt.Sprintf("%s.%s", app.ID, strings.TrimPrefix(c.ControllerURL.Host, workflowURLPrefix))

	return app, nil
}

// Get app details from a Deis controller.
func Get(c *deis.Client, appID string) (api.App, error) {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	body, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return api.App{}, err
	}

	app := api.App{}

	if err = json.Unmarshal([]byte(body), &app); err != nil {
		return api.App{}, err
	}

	// Add in app URL based on controller hostname, port included
	app.URL = fmt.Sprintf("%s.%s", app.ID, strings.TrimPrefix(c.ControllerURL.Host, workflowURLPrefix))

	return app, nil
}

// Logs retrieves logs from an app.
func Logs(c *deis.Client, appID string, lines int) (string, error) {
	u := fmt.Sprintf("/v2/apps/%s/logs", appID)

	if lines > 0 {
		u += "?log_lines=" + strconv.Itoa(lines)
	}

	body, err := c.BasicRequest("GET", u, nil)

	if err != nil || len(body) < 1 {

		return "", ErrNoLogs
	}

	// We need to trim a few characters off the front and end of the string
	return body[2 : len(body)-1], nil
}

// Run one time command in an app.
func Run(c *deis.Client, appID string, command string) (api.AppRunResponse, error) {
	req := api.AppRunRequest{Command: command}
	body, err := json.Marshal(req)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/run", appID)

	resBody, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	res := api.AppRunResponse{}

	if err = json.Unmarshal([]byte(resBody), &res); err != nil {
		return api.AppRunResponse{}, err
	}

	return res, nil
}

// Delete an app.
func Delete(c *deis.Client, appID string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}

// Transfer an app to another user.
func Transfer(c *deis.Client, appID string, username string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	req := api.AppUpdateRequest{Owner: username}
	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", u, body)
	return err
}
