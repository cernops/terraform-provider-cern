package cern

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func landbVMCardResource() *schema.Resource {
	return &schema.Resource{

		SchemaVersion: 1,

		Read:   landbVMCardResourceRead,
		Create: landbVMCardResourceCreate,
		Delete: landbVMCardResourceDelete,
		Importer: &schema.ResourceImporter{
			State: landbVMCardResourceImport,
		},

		Schema: map[string]*schema.Schema{
			"vm_name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				ForceNew:    true,
				Description: "Virtual machine host name",
			},
			"hardware_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"card_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Ethernet",
			},
		},
	}
}

func landbVMCardResourceCreate(d *schema.ResourceData, meta interface{}) error {
	landbClient, err := meta.(CernConfig).GetLandbClient()
	if err != nil {
		return err
	}

	interfaceCard := InterfaceCard{
		HardwareAddress: d.Get("hardware_address").(string),
		CardType:        d.Get("card_type").(string),
	}

	hwAddr, err := landbClient.VMAddCard(context.TODO(), d.Get("vm_name").(string), interfaceCard)
	if err != nil || hwAddr != interfaceCard.HardwareAddress {
		return fmt.Errorf("error creating VM card %s: %s", d.Get("vm_name").(string), err)
	}
	d.SetId(hwAddr)
	return nil
}

func landbVMCardResourceDelete(d *schema.ResourceData, meta interface{}) error {
	landbClient, err := meta.(CernConfig).GetLandbClient()
	if err != nil {
		return err
	}

	done, err := landbClient.VMRemoveCard(context.TODO(), d.Get("vm_name").(string), d.Get("hardware_address").(string))
	if err != nil || !done {
		return fmt.Errorf(
			"error deleting VM card %s on device: %s: %s",
			d.Get("vm_name").(string),
			d.Get("hardware_address").(string),
			err)
	}
	return nil
}

func landbVMCardResourceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func landbVMCardResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
