package cern

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider defines the schema of the CERN provider seen by Terraform
func Provider() *schema.Provider {
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
				DefaultFunc: schema.EnvDefaultFunc("CERN_TEIGI_ENDPOINT", "https://woger-direct.cern.ch:8201"),
				Description: "Teigi API url that we can use",
			},
			"certmgr_endpoint": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CERN_TEIGI_ENDPOINT", "https://hector.cern.ch:8008"),
				Description: "Certmgr API url that we can use",
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
			"cern_roger":              rogerResource(),
			"cern_certmgr":            certMgrResource(),
			"cern_teigi_secret":       resourceTeigiSecret(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Teigi client
	teigiClient, err := NewTeigiClient(d.Get("teigi_endpoint").(string))
	if err != nil {
		return nil, err
	}
	// CertMgr client
	certMgrClient, err := NewCertMgrClient(d.Get("certmgr_endpoint").(string))
	if err != nil {
		return nil, err
	}

	// Roger client
	rogerClient, err := NewRogerClient(d.Get("teigi_endpoint").(string))
	if err != nil {
		return nil, err
	}

	// Initialise Terraform provider configuration
	config := &config{
		LdapServer:    d.Get("ldap_server").(string),
		LandbEndpoint: d.Get("landb_endpoint").(string),
		LandbUsername: d.Get("landb_username").(string),
		LandbPassword: d.Get("landb_password").(string),
		TeigiClient:   teigiClient,
		RogerClient:   rogerClient,
		CertMgrClient: certMgrClient,
	}

	return config, nil
}
