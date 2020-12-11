package cern

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func landbVMResource() *schema.Resource {
	return &schema.Resource{

		SchemaVersion: 1,

		Read:   landbVMResourceRead,
		Create: landbVMResourceCreate,
		Update: landbVMResourceUpdate,
		Delete: landbVMResourceDelete,
		Importer: &schema.ResourceImporter{
			State: landbVMResourceImport,
		},

		Schema: map[string]*schema.Schema{
			"device_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Virtual machine host name",
			},
			"location": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"manufacturer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Terraform managed virtual machine",
			},
			"tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operating_system": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"landb_manager_person": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"responsible_person": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_person": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipv6_ready": {
				Required: true,
				Type:     schema.TypeBool,
			},
			"manager_locked": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func landbVMResourceCreate(d *schema.ResourceData, meta interface{}) error {
	landbClient := meta.(*Config).LandbClient
	operatingSystem := d.Get("operating_system").(map[string]interface{})
	location := d.Get("location").(map[string]interface{})
	responsible := d.Get("responsible_person").(map[string]interface{})
	landbManager := d.Get("landb_manager_person").(map[string]interface{})
	user := d.Get("user_person").(map[string]interface{})
	deviceInput := DeviceInput{
		DeviceName: d.Get("device_name").(string),
		Location: Location{
			Building: location["building"].(string),
			Floor:    location["floor"].(string),
			Room:     location["room"].(string),
		},
		Manufacturer: d.Get("manufacturer").(string),
		Model:        d.Get("model").(string),
		Description:  d.Get("description").(string),
		Tag:          d.Get("tag").(string),
		OperatingSystem: OperatingSystem{
			Name:    operatingSystem["name"].(string),
			Version: operatingSystem["version"].(string),
		},
		LandbManagerPerson: PersonInput{
			Name:       landbManager["name"].(string),
			FirstName:  landbManager["first_name"].(string),
			Department: landbManager["department"].(string),
			Group:      landbManager["group"].(string),
		},
		ResponsiblePerson: PersonInput{
			Name:       responsible["name"].(string),
			FirstName:  responsible["first_name"].(string),
			Department: responsible["department"].(string),
			Group:      responsible["group"].(string),
		},
		UserPerson: PersonInput{
			Name:       user["name"].(string),
			FirstName:  user["first_name"].(string),
			Department: user["department"].(string),
			Group:      user["group"].(string),
		},
		IPv6Ready: d.Get("ipv6_ready").(bool),
	}
	createOptions := VMCreateOptions{}

	done, err := landbClient.VMCreate(context.TODO(), deviceInput, createOptions)
	if err != nil || !done {
		return fmt.Errorf("error creating VM %s: %s", d.Get("device_name").(string), err)
	}

	d.SetId(deviceInput.DeviceName)
	return nil
}

func landbVMResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func landbVMResourceRead(d *schema.ResourceData, meta interface{}) error {
	// When the resource is in the state, this call allows to read the
	// remote API and update the values on the local state.
	//
	// At the moment, it is not likely that other entity modifies LanDB for
	// our instances, so there is no need to sync here. If that ever changes,
	// this is the place to implement the sync back to the Terraform state.
	return nil
}

func landbVMResourceDelete(d *schema.ResourceData, meta interface{}) error {
	landbClient := meta.(*Config).LandbClient
	done, err := landbClient.VMDestroy(context.TODO(), d.Get("device_name").(string))
	if err != nil || !done {
		return fmt.Errorf("error deleting VM %s: %s", d.Get("device_name").(string), err)
	}
	return nil
}

func landbVMResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
