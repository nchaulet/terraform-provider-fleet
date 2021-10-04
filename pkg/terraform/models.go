package terraform

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AgentPolicy struct {
	ID              types.String `tfsdk:"id"`
	EnrollmentToken types.String `tfsdk:"enrollment_token"`
	ConfigJSON      types.String `tfsdk:"config_json"`

	PackagePolicies types.List `tfsdk:"package_policies"`
}

type AgentPolicyWithPackagePolicy struct {
	ID              types.String `tfsdk:"id"`
	EnrollmentToken types.String `tfsdk:"enrollment_token"`
	ConfigJSON      types.String `tfsdk:"config_json"`

	PackagePolicies []PackagePolicy `tfsdk:"package_policies"`
}

type PackagePolicy struct {
	ID   string `tfsdk:"id"`
	Name string `tfsdk:"name"`
}
