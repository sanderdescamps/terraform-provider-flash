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

import (
	"context"
	"fmt"

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePureDnsSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePureDnsSettingsCreateUpdate,
		ReadContext:   resourcePureDnsSettingsRead,
		UpdateContext: resourcePureDnsSettingsCreateUpdate,
		DeleteContext: resourcePureDnsSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nameservers": {
				Type:        schema.TypeList,
				Description: "A list of up to three DNS server IP addresses",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
				MaxItems: 3,
				Required: true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Network address for the network interface",
				Required:    false,
				Optional:    true,
			},
		},
	}
}

func resourcePureDnsSettingsCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*flasharray.Client)
	data := make(map[string]interface{})

	if domain, ok := d.GetOk("domain"); ok {
		data["domain"] = domain
	} else {
		data["domain"] = ""
	}

	if nameservers, ok := d.GetOk("nameservers"); ok {
		data["nameservers"] = nameservers
	} else {
		data["nameservers"] = []string{}
	}

	if dnsSettings, err := client.Networks.SetDNS(data); err != nil {
		return diag.FromErr(err)
	} else {
		d.Set("nameservers", dnsSettings.Nameservers)
		d.Set("domain", dnsSettings.Domain)
		d.SetId(fmt.Sprintf("dns-settings-%s", client.Target))
	}

	return nil
}

func resourcePureDnsSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	if dnsSettings, err := client.Networks.GetDNS(); err != nil {
		return diag.FromErr(err)
	} else {
		d.Set("nameservers", dnsSettings.Nameservers)
		d.Set("domain", dnsSettings.Domain)
		d.SetId(fmt.Sprintf("dns-settings-%s", client.Target))
	}

	return nil
}

func resourcePureDnsSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := m.(*flasharray.Client)

	// data := make(map[string]interface{})
	// data["nameservers"] = []string{}
	// data["domain"] = ""

	// if dnsSettings, err := client.Networks.SetDNS(data); err != nil {
	// 	return diag.FromErr(err)
	// } else {
	// 	d.Set("nameservers", dnsSettings.Nameservers)
	// 	d.Set("domain", dnsSettings.Domain)
	// 	d.SetId(fmt.Sprintf("dns-settings-%s", client.Target))
	// }

	d.SetId("")
	return nil
}
