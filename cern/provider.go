package cern

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider defines the schema of the CERN provider seen by Terraform
func Provider() terraform.ResourceProvider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"ldap_server": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("CERN_LDAP_SERVER", "ldap://xldap.cern.ch:389"),
				Optional:    true,
			},
			"landb_endpoint": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("CERN_LANDB_ENDPOINT", "https://network.cern.ch/sc/soap/soap.fcgi?v=6"),
				Optional:    true,
			},
			"landb_username": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("CERN_LANDB_USERNAME", ""),
				Optional:    true,
			},
			"landb_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CERN_LANDB_PASSWORD", ""),
			},
			"teigi_endpoint": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CERN_TEIGI_ENDPOINT", "https://woger.cern.ch:8201"),
				Description: "Teigi API url that we can use",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cern_egroup":       dataSourceCernEgroup(),
			"cern_teigi_secret": dataSourceTeigiSecret(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cern_landb_vm":           landbVMResource(),
			"cern_landb_vm_card":      landbVMCardResource(),
			"cern_landb_vm_interface": landbVMInterfaceResource(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// This LanDB client is initialised with a token that should be valid for
	// a few hours. A renovation mechanism has not been implemented yet.
	landbClient, err := NewLandbClient(
		d.Get("landb_endpoint").(string),
		d.Get("landb_username").(string),
		d.Get("landb_password").(string),
	)
	if err != nil {
		return nil, err
	}

	// Teigi client
	teigiClient, err := NewTeigiClient(d.Get("teigi_endpoint").(string))
	if err != nil {
		return nil, err
	}

	// Initialise Terraform provider configuration
	config := &Config{
		LdapServer:  d.Get("ldap_server").(string),
		LandbClient: landbClient,
		TeigiClient: teigiClient,
	}

	return config, nil
}
