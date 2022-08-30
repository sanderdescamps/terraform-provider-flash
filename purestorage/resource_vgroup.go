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

func resourcePureVolumegroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePureVolumegroupCreate,
		ReadContext:   resourcePureVolumegroupRead,
		UpdateContext: resourcePureVolumegroupUpdate,
		DeleteContext: resourcePureVolumegroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the volume group",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w\-\d]+$`), "can only contain letters, numbers and '-'"),
			},
		},
	}
}

func resourcePureVolumegroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	if vgroup, err := client.Vgroups.CreateVgroup(d.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	} else {
		d.Set("name", vgroup.Name)
		d.SetId(vgroup.Name)
	}

	return nil
}

func resourcePureVolumegroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	if vgroup, err := client.Vgroups.GetVgroup(d.Id()); err != nil {
		d.SetId("")
		return nil
	} else {
		d.Set("name", vgroup.Name)
		d.SetId(vgroup.Name)
	}

	return nil
}

func resourcePureVolumegroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*flasharray.Client)

	if d.HasChange("name") {
		c, n := d.GetChange("name")
		if v, err := client.Vgroups.RenameVgroup(c.(string), n.(string)); err != nil {
			return diag.FromErr(err)
		} else {
			d.SetId(v.Name)
			d.Set("name", d.Get("name").(string))
		}

	}

	return resourcePureVolumegroupRead(ctx, d, m)
}

// resourcePureVolumeDelete will delete the volumegroup specified.
func resourcePureVolumegroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)
	_, err := client.Vgroups.DestroyVgroup(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
