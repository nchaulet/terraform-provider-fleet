package fleetapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Package struct {
	Name    string `json:"name,omitempty"`
	Title   string `json:"title,omitempty"`
	Version string `json:"version,omitempty"`
}

type PackagePolicyInput struct {
	Type           string `json:"type,omitempty"`
	PolicyTemplate string `json:"policy_template,omitempty"`
	Enabled        bool   `json:"enabled,omitempty"`
}

type PostPackagePolicyResponse struct {
	Item *struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"item"`
}

func (c *Client) PostPackagePolicies(ctx context.Context, request *map[string]interface{}) (*PostPackagePolicyResponse, error) {
	b, err := json.Marshal(request)
	if err != nil {
		// TODO better error handling
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/fleet/package_policies", c.BaseURL), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := PostPackagePolicyResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type DeleteAgentPolicyRequest struct {
	PackagePolicyIDS []string `json:"packagePolicyIds"`
}

func (c *Client) DeletePackagePolicies(ctx context.Context, request *DeleteAgentPolicyRequest) (*[]map[string]interface{}, error) {
	b, err := json.Marshal(request)
	if err != nil {
		// TODO better error handling
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/fleet/package_policies/delete", c.BaseURL), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := []map[string]interface{}{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
