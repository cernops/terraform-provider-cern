package main

import (
	"github.com/hashicorp/terraform/plugin"
	"gitlab.cern.ch/batch-team/infra/terraform-provider-cern/cern"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cern.Provider})
}
