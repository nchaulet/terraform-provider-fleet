package fleetapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClientBasicAuth struct {
	Username string
	Password string
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Auth       ClientBasicAuth
}

func NewClient(kibanaHost string, AuthOptions ClientBasicAuth) *Client {
	return &Client{
		BaseURL: kibanaHost,
		Auth:    AuthOptions,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("kbn-xsrf", "xxx")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dat, &v)
	if err != nil {
		return err
	}

	return nil
}
