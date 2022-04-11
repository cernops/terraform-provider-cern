package cern

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeigiSecret() *schema.Resource {
	return &schema.Resource{
		Read: teigiSecretDataSourceRead,

		Schema: map[string]*schema.Schema{
			"hostgroup": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostgroup where the secret is located",
			},

			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key name which to retrieve",
			},

			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Secret string retrieved from Teigi",
			},
		},
	}
}

func teigiSecretDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config).TeigiClient

	hostgroup := d.Get("hostgroup").(string)
	key := d.Get("key").(string)
	log.Printf("[DEBUG] Creating Teigi request for %s hostgroup for %s key", hostgroup, key)

	secret, msg, err := client.Get(hostgroup, key)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("%s: %v", msg, err))
	}

	d.SetId(key)
	if err := d.Set("secret", secret.Secret); err != nil {
		return fmt.Errorf("Unable to set secret: %s", err)
	}

	return nil
}
