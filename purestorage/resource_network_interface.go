/*
   Copyright 2018 David Evans

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package purestorage

// The Pure API 1.x does not support setting network interfaces. Therefore this resource is disabled.
/*
import (
	"context"

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePureNetworkInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePureNetworkInterfaceCreate,
		ReadContext:   resourcePureNetworkInterfaceRead,
		UpdateContext: resourcePureNetworkInterfaceUpdate,
		DeleteContext: resourcePureNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the network interface",
				Required:    true,
				ForceNew:    true,
			},
			"address": {
				Type:         schema.TypeString,
				Description:  "Network address for the network interface",
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"gateway": {
				Type:         schema.TypeString,
				Description:  "Gateway address for the network interface",
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"netmask": {
				Type:         schema.TypeString,
				Description:  "Subnet mask i the form ddd.ddd.ddd.ddd (ex. 255.255.255.0)",
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable/disable network interface",
				Optional:    true,
				Default:     false,
			},
			"mtu": {
				Type:         schema.TypeInt,
				Description:  "mtu for the network interface",
				Optional:     true,
				Default:      1500,
				ValidateFunc: validation.IntBetween(568, 9000),
			},
			"mac": {
				Type:        schema.TypeString,
				Description: "mac address",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
			"speed": {
				Type:        schema.TypeInt,
				Description: "Interface speed",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
			"speed_gbs": {
				Type:        schema.TypeInt,
				Description: "Interface speed in Gb/s",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}
}

func resourcePureNetworkInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*flasharray.Client)
	name, _ := d.GetOk("name")
	data := make(map[string]interface{})

	if d.Get("enabled").(bool) {
		data["enabled"] = true
		data["mtu"] = d.Get("mtu").(int)
		if address, ok := d.GetOk("address"); ok {
			data["address"] = address
		}

		if netmask, ok := d.GetOk("netmask"); ok {
			data["netmask"] = netmask
		}

		if gateway, ok := d.GetOk("gateway"); ok {
			data["gateway"] = gateway
		}
	} else {
		data["enabled"] = false
	}

	netInterface, err := client.Networks.SetNetworkInterface(name.(string), data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", netInterface.Name)
	d.Set("address", netInterface.Address)
	d.Set("netmask", netInterface.Netmask)
	d.Set("gateway", netInterface.Gateway)
	d.Set("mtu", netInterface.Mtu)
	d.Set("speed", netInterface.Speed)
	d.Set("speed_gbs", netInterface.Speed/1000000000)
	d.Set("mac", netInterface.Hwaddr)
	d.SetId(netInterface.Name)

	return nil
}

func resourcePureNetworkInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	name := d.Id()

	netInterface, err := client.Networks.GetNetworkInterface(name)
	if err != nil {
		return diag.FromErr(err)
	} else if netInterface == nil {
		d.SetId("")
	} else {
		d.Set("name", netInterface.Name)
		d.Set("address", netInterface.Address)
		d.Set("netmask", netInterface.Netmask)
		d.Set("gateway", netInterface.Gateway)
		d.Set("mtu", netInterface.Mtu)
		d.Set("speed", netInterface.Speed)
		d.Set("speed_gbs", netInterface.Speed/1000000000)
		d.Set("mac", netInterface.Hwaddr)
		d.SetId(netInterface.Name)
	}

	return nil
}

func resourcePureNetworkInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)
	data := make(map[string]interface{})

	if d.HasChange("mtu") {
		data["mtu"] = d.Get("mtu").(int)
	}

	if d.HasChange("address") {
		data["address"] = d.Get("address").(string)
	}

	if d.HasChange("netmask") {
		data["netmask"] = d.Get("netmask").(string)
	}

	if d.HasChange("gateway") {
		data["gateway"] = d.Get("gateway").(string)
	}

	if len(data) > 0 || d.HasChange("enabled") {
		data["enabled"] = d.Get("enabled")
		netInterface, err := client.Networks.SetNetworkInterface(d.Id(), data)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("name", netInterface.Name)
		d.Set("address", netInterface.Address)
		d.Set("netmask", netInterface.Netmask)
		d.Set("gateway", netInterface.Gateway)
		d.Set("mtu", netInterface.Mtu)
		d.Set("speed", netInterface.Speed)
		d.Set("speed_gbs", netInterface.Speed/1000000000)
		d.Set("mac", netInterface.Hwaddr)

	}

	return nil
}

func resourcePureNetworkInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	if _, err := client.Networks.DisableNetworkInterface(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

*/
