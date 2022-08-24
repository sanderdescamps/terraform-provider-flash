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

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePureFlashArray() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePureFlashArrayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"revision": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePureFlashArrayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	flasharray, err := client.Array.Get(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(flasharray.ID)
	d.Set("name", flasharray.ArrayName)
	d.Set("version", flasharray.Version)
	d.Set("revision", flasharray.Revision)
	return nil
}
