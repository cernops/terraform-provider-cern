package cern

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func rogerResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRogerCreate,
		ReadContext:   resourceRogerRead,
		UpdateContext: resourceRogerUpdate,
		DeleteContext: resourceRogerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Hostname to update the roger state with",
			},
			"appstate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_alarmed": {
				Type:        schema.TypeBool,
				Description: "Toggle on or off application alarms",
				Optional:    true,
				Default:     false,
			},
			"hw_alarmed": {
				Type:        schema.TypeBool,
				Description: "Toggle on or off hardware alarms",
				Optional:    true,
				Default:     false,
			},
			"nc_alarmed": {
				Type:        schema.TypeBool,
				Description: "Toggle on or off no contact alarms",
				Optional:    true,
				Default:     false,
			},
			"os_alarmed": {
				Type:        schema.TypeBool,
				Description: "Toggle on or off operating system alarms",
				Optional:    true,
				Default:     false,
			},
			"message": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expires": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expires_dt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time_dt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by_puppet": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceRogerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).RogerClient

	request := RogerRequest{
		Hostname: d.Get("hostname").(string),
		Expires:  d.Get("expires").(string),
		Message:  d.Get("message").(string),
	}
	appAlarmed := d.Get("app_alarmed").(bool)
	if appAlarmed {
		request.AppAlarmed = &appAlarmed
	}
	hwAlarmed := d.Get("hw_alarmed").(bool)
	if appAlarmed {
		request.HwAlarmed = &hwAlarmed
	}
	ncAlarmed := d.Get("nc_alarmed").(bool)
	if ncAlarmed {
		request.NcAlarmed = &ncAlarmed
	}
	osAlarmed := d.Get("os_alarmed").(bool)
	if ncAlarmed {
		request.OsAlarmed = &osAlarmed
	}

	log.Printf("[DEBUG] Creating roger create request for %s", request.Hostname)
	err := client.Create(ctx, request)
	if err != nil {
		return diag.Errorf("Error creating roger state: %s", err)
	}

	d.SetId(request.Hostname)

	return resourceRogerRead(ctx, d, meta)
}

func resourceRogerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).RogerClient
	hostname := d.Id()
	log.Printf("[DEBUG] Creating roger read request for %s", hostname)
	resp, err := client.Get(ctx, hostname)
	if err != nil {
		return diag.Errorf("Error reading roger state: %s", err)
	}

	if err := d.Set("hostname", hostname); err != nil {
		return diag.Errorf("Unable to set hostname: %s", err)
	}
	if err := d.Set("appstate", resp.AppState); err != nil {
		return diag.Errorf("Error setting appstate: %s", err)
	}
	if err := d.Set("app_alarmed", resp.AppAlarmed); err != nil {
		return diag.Errorf("Error setting app_alarmed: %s", err)
	}
	if err := d.Set("hw_alarmed", resp.HwAlarmed); err != nil {
		return diag.Errorf("Error setting hw_alarmed: %s", err)
	}
	if err := d.Set("os_alarmed", resp.OsAlarmed); err != nil {
		return diag.Errorf("Error setting os_alarmed: %s", err)
	}
	if err := d.Set("message", resp.Message); err != nil {
		return diag.Errorf("Error setting message: %s", err)
	}
	if err := d.Set("expires", resp.Expires); err != nil {
		return diag.Errorf("Error setting expires: %s", err)
	}
	if err := d.Set("expires_dt", resp.ExpiresDt); err != nil {
		return diag.Errorf("Error setting expires_dt: %s", err)
	}
	if err := d.Set("update_time", resp.UpdateTime); err != nil {
		return diag.Errorf("Error setting update_time: %s", err)
	}
	if err := d.Set("update_time_dt", resp.UpdateTimeDt); err != nil {
		return diag.Errorf("Error setting update_time_dt: %s", err)
	}
	if err := d.Set("updated_by", resp.UpdatedBy); err != nil {
		return diag.Errorf("Error setting updated_by: %s", err)
	}
	if err := d.Set("updated_by_puppet", resp.UpdatedByPuppet); err != nil {
		return diag.Errorf("Error setting updated_by_puppet: %s", err)
	}

	return nil
}

func resourceRogerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).RogerClient

	request := RogerRequest{
		Hostname: d.Get("hostname").(string),
		Expires:  d.Get("expires").(string),
		Message:  d.Get("message").(string),
	}
	appAlarmed := d.Get("app_alarmed").(bool)
	if appAlarmed {
		request.AppAlarmed = &appAlarmed
	}
	hwAlarmed := d.Get("hw_alarmed").(bool)
	if appAlarmed {
		request.HwAlarmed = &hwAlarmed
	}
	ncAlarmed := d.Get("nc_alarmed").(bool)
	if ncAlarmed {
		request.NcAlarmed = &ncAlarmed
	}
	osAlarmed := d.Get("os_alarmed").(bool)
	if ncAlarmed {
		request.OsAlarmed = &osAlarmed
	}

	log.Printf("[DEBUG] Creating roger update request for %s", request.Hostname)
	err := client.Update(ctx, request)
	if err != nil {
		return diag.Errorf("Error updating roger state: %s", err)
	}

	return resourceRogerRead(ctx, d, meta)
}

func resourceRogerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config).RogerClient
	hostname := d.Get("hostname").(string)

	log.Printf("[DEBUG] Creating roger delete request for %s", hostname)
	_, err := client.Delete(ctx, hostname)
	if err != nil {
		return diag.Errorf("Error deleting roger state: %s", err)
	}

	return nil
}
