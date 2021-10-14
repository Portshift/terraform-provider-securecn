package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-securecn/securecn"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: securecn.Provider(),
	})
}
