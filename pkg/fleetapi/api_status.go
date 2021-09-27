package fleetapi

import (
	"context"
	"fmt"
	"net/http"
)

type StatusRes struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

func (c *Client) GetStatus(ctx context.Context) (*StatusRes, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/status", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := StatusRes{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
