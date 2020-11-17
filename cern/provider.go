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
				Default:     "ldap://xldap.cern.ch:389",
				Optional:    true,
				Description: "LDAP server to connect to",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cern_egroup": dataSourceCernEgroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		LdapServer: d.Get("ldap_server").(string),
	}
	return config, nil
}
