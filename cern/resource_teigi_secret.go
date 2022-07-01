package cern

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTeigiSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeigiSecretCreate,
		// Use the read func from the data source
		ReadContext:   resourceTeigiSecretRead,
		UpdateContext: resourceTeigiSecretUpdate,
		DeleteContext: resourceTeigiSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"host": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Description:   "Host where the secret is located",
				ConflictsWith: []string{"hostgroup", "service"},
				ValidateFunc:  validation.StringDoesNotContainAny("/"),
			},

			"hostgroup": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Description:   "Hostgroup where the secret is located",
				ConflictsWith: []string{"host", "service"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return SerializeHostgroup(old) == SerializeHostgroup(new)
				},
			},

			"service": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Description:   "Service where the secret is located",
				ConflictsWith: []string{"host", "hostgroup"},
				ValidateFunc:  validation.StringDoesNotContainAny("/"),
			},

			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key name which to retrieve",
			},

			"base64_encoded": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether or not the secret is base64 encoded",
			},

			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Secret string retrieved from Teigi",
			},
		},
	}
}

func parseTeigiSecretID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", "", fmt.Errorf("Unable to determine teigi secret ID %s", id)
	}

	return idParts[0], idParts[1], idParts[2], nil
}

func resourceTeigiSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).TeigiClient

	scope, entity, err := getScopeAndEntity(d)
	if err != nil {
		return diag.Errorf("Error getting scope and entity: %s", err)
	}

	key := d.Get("key").(string)
	log.Printf("[DEBUG] Creating Teigi create request for %s scope and %s entity for %s key",
		scope, entity, key)
	request := SecretRequest{
		Secret: d.Get("secret").(string),
	}
	if d.Get("base64_encoded").(bool) {
		request.Encoding = "b64"
	}

	err = client.Create(ctx, scope, entity, key, request)
	if err != nil {
		return diag.Errorf("Unable to create secret: %s", err)
	}

	return resourceTeigiSecretRead(ctx, d, meta)
}

func resourceTeigiSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).TeigiClient
	scope, entity, key, err := parseTeigiSecretID(d.Id())
	if err != nil {
		return diag.Errorf("Unable to parse teigi secret id: %s", err)
	}

	log.Printf("[DEBUG] Creating Teigi read request for %s scope and %s entity for %s key",
		scope, entity, key)
	secretResponse, err := client.Get(ctx, scope, entity, key)
	if err != nil {
		return diag.Errorf("Unable to get secret: %s", err)
	}

	if err := d.Set(scope, entity); err != nil {
		return diag.Errorf("Unable to set '%s': %s", scope, err)
	}
	if err := d.Set("key", key); err != nil {
		return diag.Errorf("Unable to set key: %s", err)
	}

	if err := d.Set("secret", secretResponse.Secret); err != nil {
		return diag.Errorf("Unable to set secret: %s", err)
	}

	base64_encoded := secretResponse.Encoding == "b64"
	if err := d.Set("base64_encoded", base64_encoded); err != nil {
		return diag.Errorf("Unable to set base64_encoded: %s", err)
	}
	return nil
}

func resourceTeigiSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceTeigiSecretCreate(ctx, d, meta)
}

func resourceTeigiSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).TeigiClient
	scope, entity, err := getScopeAndEntity(d)
	if err != nil {
		return diag.Errorf("Error getting scope and entity: %s", err)
	}

	key := d.Get("key").(string)
	log.Printf("[DEBUG] Creating Teigi delete request for %s scope and %s entity for %s key",
		scope, entity, key)
	err = client.Delete(ctx, scope, entity, key)
	if err != nil {
		return diag.Errorf("Unable to delete secret: %s", err)
	}

	return nil
}
