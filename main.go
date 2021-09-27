package main

import (
	"context"
	"terraform-provider-fleet/pkg/terraform"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	tfsdk.Serve(context.Background(), terraform.New, tfsdk.ServeOpts{
		Name: "fleet",
	})

}
