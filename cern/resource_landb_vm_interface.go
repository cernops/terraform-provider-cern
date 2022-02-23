package cern

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func landbVMInterfaceResource() *schema.Resource {
	return &schema.Resource{

		SchemaVersion: 1,

		Read:   landbVMInterfaceResourceRead,
		Create: landbVMInterfaceResourceCreate,
		Delete: landbVMInterfaceResourceDelete,
		Importer: &schema.ResourceImporter{
			State: landbVMInterfaceResourceImport,
		},

		Schema: map[string]*schema.Schema{
			"vm_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Virtual machine host name",
				ForceNew:    true,
			},
			"interface_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cern.ch",
				ForceNew: true,
			},
			"vm_cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_interface_options": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func landbVMInterfaceResourceCreate(d *schema.ResourceData, meta interface{}) error {
	landbClient, err := meta.(CernConfig).GetLandbClient()
	if err != nil {
		return err
	}

	vmInterfaceOptions := d.Get("vm_interface_options").(map[string]interface{})
	interfaceName := strings.ToUpper(fmt.Sprintf("%s.%s", d.Get("vm_name").(string), d.Get("interface_domain").(string)))

	interfaceRequest := VMAddInterfaceRequest{
		VMName:        d.Get("vm_name").(string),
		InterfaceName: interfaceName,
		VMClusterName: d.Get("vm_cluster_name").(string),
		VMInterfaceOptions: VMInterfaceOptions{
			IP:          vmInterfaceOptions["ip"].(string),
			AddressType: vmInterfaceOptions["address_type"].(string),
			ServiceName: vmInterfaceOptions["service_name"].(string),
		},
	}

	done, err := landbClient.VMAddInterface(context.TODO(), interfaceRequest)
	if err != nil || !done {
		return fmt.Errorf(
			"error creating VM interface %s for device %s: %s",
			interfaceName,
			d.Get("vm_name").(string),
			err)
	}
	d.SetId(interfaceName)
	return nil
}

func landbVMInterfaceResourceDelete(d *schema.ResourceData, meta interface{}) error {
	landbClient, err := meta.(CernConfig).GetLandbClient()
	if err != nil {
		return err
	}

	interfaceName := strings.ToUpper(fmt.Sprintf("%s.%s", d.Get("vm_name").(string), d.Get("interface_domain").(string)))
	done, err := landbClient.VMRemoveInterface(context.TODO(), d.Get("vm_name").(string), interfaceName)
	if err != nil || !done {
		return fmt.Errorf(
			"error deleting VM interface %s on device: %s: %s",
			d.Get("vm_name").(string),
			interfaceName,
			err)
	}
	return nil
}

func landbVMInterfaceResourceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func landbVMInterfaceResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
