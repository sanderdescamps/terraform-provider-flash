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
	"regexp"

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

func resourcePureAlertRecipient() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage alert recepients on Pure Flash array.",
		CreateContext: resourcePureAlertRecipientCreate,
		ReadContext:   resourcePureAlertRecipientRead,
		UpdateContext: resourcePureAlertRecipientUpdate,
		DeleteContext: resourcePureAlertRecipientDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:         schema.TypeString,
				Description:  "Email address",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(emailRegex), "not a valid email address format"),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable/disable email recipient",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourcePureAlertRecipientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*flasharray.Client)
	email := d.Get("email").(string)

	if alert, err := client.Alerts.CreateAlert(email, nil); err != nil {
		diag.FromErr(err)
	} else {
		d.Set("email", alert.Name)
		d.SetId(alert.Name)
	}

	if e, ok := d.GetOk("enabled"); ok && !e.(bool) {
		data := make(map[string]interface{})
		data["enabled"] = false
		if _, err := client.Alerts.SetAlert(email, data); err != nil {
			diag.FromErr(err)
		}
		d.Set("enabled", false)
	} else {
		d.Set("enabled", true)
	}

	return nil
}

func resourcePureAlertRecipientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	alert, err := client.Alerts.GetAlert(d.Id())
	if err != nil {
		return diag.FromErr(err)
	} else if alert == nil {
		d.SetId("")
		return nil
	}
	d.Set("email", alert.Name)
	d.Set("enabled", alert.Enabled)
	d.SetId(alert.Name)
	return nil
}

func resourcePureAlertRecipientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	if d.HasChange("enabled") {
		data := make(map[string]interface{})
		data["enabled"] = d.Get("enabled").(bool)
		if alert, err := client.Alerts.SetAlert(d.Id(), data); err != nil {
			diag.FromErr(err)
		} else {
			d.Set("enabled", alert.Enabled)
		}

	}

	return nil
}

func resourcePureAlertRecipientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	if _, err := client.Alerts.DeleteAlert(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
