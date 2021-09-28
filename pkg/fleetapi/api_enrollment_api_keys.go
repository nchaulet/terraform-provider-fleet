package fleetapi

import (
	"context"
	"fmt"
	"net/http"
)

type EnrollmentApiKey struct {
	PolicyId string `json:"policy_id"`
	Active   bool   `json:"active"`
	ApiKey   string `json:"api_key"`
}

type GetEnrollmentApiKeysResponse struct {
	List []*EnrollmentApiKey `json:"list"`
}

func (c *Client) GetEnrollmentTokens(ctx context.Context, policyID string) (*GetEnrollmentApiKeysResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/fleet/enrollment-api-keys?kuery=policy_id:\"%s\"", c.BaseURL, policyID), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := GetEnrollmentApiKeysResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
