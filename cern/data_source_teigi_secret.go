package cern

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceTeigiSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeigiSecretRead,

		Schema: map[string]*schema.Schema{
			"host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Host where the secret is located",
				ConflictsWith: []string{"hostgroup", "service"},
				ValidateFunc:  validation.StringDoesNotContainAny("/"),
			},

			"hostgroup": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Hostgroup where the secret is located",
				ConflictsWith: []string{"host", "service"},
			},

			"service": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Service where the secret is located",
				ConflictsWith: []string{"host", "hostgroup"},
				ValidateFunc:  validation.StringDoesNotContainAny("/"),
			},

			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key name which to retrieve",
			},

			"base64_encoded": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the secret is base64 encoded",
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

func getScopeAndEntity(d *schema.ResourceData) (string, string, error) {
	var scope string
	var entity string

	hostgroup := d.Get("hostgroup").(string)
	host := d.Get("host").(string)
	service := d.Get("service").(string)

	if hostgroup != "" {
		scope = "hostgroup"
		entity = hostgroup
	} else if host != "" {
		scope = "host"
		entity = host
	} else if service != "" {
		scope = "service"
		entity = service
	} else {
		return "", "", fmt.Errorf("one of the variable 'hostgroup', 'host' or 'service' should be set")
	}

	return scope, entity, nil
}

func dataSourceTeigiSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).TeigiClient

	scope, entity, err := getScopeAndEntity(d)
	if err != nil {
		return diag.Errorf("Error getting scope and entity: %s", err)
	}

	key := d.Get("key").(string)
	log.Printf("[DEBUG] Creating Teigi read request for %s scope and %s entity for %s key",
		scope, entity, key)

	secretResponse, err := client.Get(ctx, scope, entity, key)
	if err != nil {
		return diag.Errorf("Unable to get secret: %s", err)
	}

	d.SetId(scope + "/" + entity + "/" + key)
	if err = d.Set("secret", secretResponse.Secret); err != nil {
		return diag.Errorf("Unable to set secret: %s", err)
	}

	base64_encoded := secretResponse.Encoding == "b64"
	if err := d.Set("base64_encoded", base64_encoded); err != nil {
		return diag.Errorf("Unable to set base64_encoded: %s", err)
	}
	return nil
}
