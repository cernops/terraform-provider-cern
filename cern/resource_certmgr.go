package cern

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func certMgrResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertMgrCreate,
		ReadContext:   resourceCertMgrRead,
		DeleteContext: resourceCertMgrDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Hostname to stage the certificate",
			},
			"id_cert": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"requestor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceCertMgrCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).CertMgrClient
	hostname := d.Get("hostname").(string)
	log.Printf("[DEBUG] Creating CertMgr request for %s", hostname)

	resp, err := client.Do(ctx, hostname)
	if err != nil {
		return diag.Errorf("Error staging certificate: %s", err)
	}

	d.SetId(hostname)

	if err := d.Set("id_cert", resp.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}
	if err := d.Set("requestor", resp.Requestor); err != nil {
		return diag.Errorf("Error setting requestor: %s", err)
	}
	if err := d.Set("start", resp.Start); err != nil {
		return diag.Errorf("Error setting start: %s", err)
	}
	if err := d.Set("end", resp.End); err != nil {
		return diag.Errorf("Error setting end: %s", err)
	}

	return nil
}

func resourceCertMgrRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := d.Set("hostname", d.Id()); err != nil {
		return diag.Errorf("Unable to set hostname: %s", err)
	}
	return nil
}

func resourceCertMgrDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
