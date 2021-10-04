package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-fleet/pkg/fleetapi"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceAgentPolicyType struct{}

// Order Resource schema
func (r resourceAgentPolicyType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Agent policy",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"enrollment_token": {
				Type:      types.StringType,
				Computed:  true,
				Sensitive: true,
			},
			"package_policies": {
				Computed: true,
				// Optional: true,
				// Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				// 	"name": {
				// 		Type:     types.StringType,
				// 		Computed: true,
				// 	},
				// 	"id": {
				// 		Type:     types.StringType,
				// 		Computed: true,
				// 	},
				// }, tfsdk.ListNestedAttributesOptions{}),
				Type: types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name": types.StringType,
						"id":   types.StringType,
					},
				}},
			},
			"config_json": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

// New resource instance
func (r resourceAgentPolicyType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceAgentPolicy{
		p: *(p.(*provider)),
	}, nil
}

type resourceAgentPolicy struct {
	p provider
}

func (r resourceAgentPolicy) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
}

// Create a new resource
func (r resourceAgentPolicy) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan AgentPolicy
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var raw map[string]interface{}

	err := json.Unmarshal([]byte(plan.ConfigJSON.Value), &raw)
	if err != nil {
		resp.AddError("Error reading config_json", fmt.Sprintf("... details ... %s", err))
	}

	var namespace string
	if raw["namespace"] == nil {
		namespace = "default"
	} else {
		namespace = raw["namespace"].(string)
	}

	var name string
	if raw["name"] == nil {
		name = "Default"
	} else {
		name = raw["name"].(string)
	}

	postAgentPoliciesRes, err := r.p.client.PostAgentPolicies(ctx, &fleetapi.AgentPolicyRequest{
		Namespace: namespace,
		Name:      name,
	})

	if err != nil {
		resp.AddError("Error creating agent policy", fmt.Sprintf("... details ... %s", err))
		return
	}

	enrollmentTokenRes, err := r.p.client.GetEnrollmentTokens(ctx, postAgentPoliciesRes.Item.ID)
	if err != nil {
		resp.AddError("Error retrieving agent policy enrollment token", fmt.Sprintf("... details ... %s", err))
		return
	}

	if len(enrollmentTokenRes.List) <= 0 {
		resp.AddError("Enrollment token not found for agent policy", fmt.Sprintf("Policy id %s", postAgentPoliciesRes.Item.ID))
		return
	}

	policyID := postAgentPoliciesRes.Item.ID
	var packagePolicies = raw["package_policies"].([](interface{}))
	var agentPolicyPackagePolicies []PackagePolicy
	for _, packagePolicy := range packagePolicies {

		packagePolicyRequest := packagePolicy.(map[string]interface{})
		packagePolicyRequest["policy_id"] = policyID
		packagePolicyRes, err := r.p.client.PostPackagePolicies(ctx, &packagePolicyRequest)

		if err != nil {
			resp.AddError("Error creating package policy", fmt.Sprintf("... details ... %s", err))
			return
		}

		agentPolicyPackagePolicies = append(agentPolicyPackagePolicies, PackagePolicy{
			ID:   packagePolicyRes.Item.ID,
			Name: packagePolicyRequest["name"].(string),
		},
		)
	}

	var result = AgentPolicyWithPackagePolicy{
		ID:              types.String{Value: policyID},
		EnrollmentToken: types.String{Value: enrollmentTokenRes.List[0].ApiKey},
		ConfigJSON:      types.String{Value: plan.ConfigJSON.Value},
		PackagePolicies: agentPolicyPackagePolicies,
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceAgentPolicy) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
}

// // Update resource
func (r resourceAgentPolicy) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan AgentPolicy
	var state AgentPolicyWithPackagePolicy
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	policyID := state.ID.Value

	var raw map[string]interface{}

	err := json.Unmarshal([]byte(plan.ConfigJSON.Value), &raw)
	if err != nil {
		resp.AddError("Error reading config_json", fmt.Sprintf("... details ... %s", err))
	}

	// Delete old package policy
	var packagePolicyIds []string
	for _, statePackagePolicy := range state.PackagePolicies {
		packagePolicyIds = append(packagePolicyIds, statePackagePolicy.ID)
	}

	if len(packagePolicyIds) > 0 {
		_, err := r.p.client.DeletePackagePolicies(ctx, &fleetapi.DeleteAgentPolicyRequest{
			PackagePolicyIDS: packagePolicyIds,
		})
		if err != nil {
			resp.AddError("Error deleting package policies", fmt.Sprintf("... details ... %s", err))
		}
	}

	// resp.AddWarning("test1", fmt.Sprintf("%#v", state.PackagePolicies))

	// Create new one
	var packagePolicies = raw["package_policies"].([](interface{}))
	var agentPolicyPackagePolicies []PackagePolicy
	for _, packagePolicy := range packagePolicies {

		packagePolicyRequest := packagePolicy.(map[string]interface{})
		packagePolicyRequest["policy_id"] = policyID
		packagePolicyRes, err := r.p.client.PostPackagePolicies(ctx, &packagePolicyRequest)

		if err != nil {
			resp.AddError("Error creating package policy", fmt.Sprintf("... details ... %s", err))
			return
		}

		agentPolicyPackagePolicies = append(agentPolicyPackagePolicies, PackagePolicy{
			ID:   packagePolicyRes.Item.ID,
			Name: packagePolicyRequest["name"].(string),
		},
		)
	}

	var result = AgentPolicyWithPackagePolicy{
		ID:              state.ID,
		EnrollmentToken: state.EnrollmentToken,
		ConfigJSON:      plan.ConfigJSON,
		PackagePolicies: agentPolicyPackagePolicies,
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceAgentPolicy) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}

func (r resourceAgentPolicy) ModifyPlan(ctx context.Context, req tfsdk.ModifyResourcePlanRequest, resp *tfsdk.ModifyResourcePlanResponse) {
	var plan AgentPolicy
	var state AgentPolicy
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	if !plan.ConfigJSON.Equal(state.ConfigJSON) {
		var result = AgentPolicy{
			ID:              plan.ID,
			EnrollmentToken: plan.EnrollmentToken,
			ConfigJSON:      plan.ConfigJSON,
			PackagePolicies: types.List{Unknown: true},
		}

		diags := resp.Plan.Set(ctx, result)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

}
