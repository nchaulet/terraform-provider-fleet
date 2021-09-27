package fleetapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AgentPolicyRequest struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace"`
}

type AgentPolicyResponse struct {
	Item *struct {
		ID        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	} `json:"item"`
}

func (c *Client) PostAgentPolicies(ctx context.Context, request *AgentPolicyRequest) (*AgentPolicyResponse, error) {
	b, err := json.Marshal(request)
	if err != nil {
		// TODO better error handling
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/fleet/agent_policies", c.BaseURL), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := AgentPolicyResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
