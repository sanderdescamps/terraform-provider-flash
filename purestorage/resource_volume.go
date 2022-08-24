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
	"log"

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePureVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePureVolumeCreate,
		ReadContext:   resourcePureVolumeRead,
		UpdateContext: resourcePureVolumeUpdate,
		DeleteContext: resourcePureVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePureVolumeImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Computed: true,
			},
			"serial": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// resourcePureVolumeCreate creates a Pure Volume on a FlashArray according
// to the schema Resource Data provided.
// If the size parameter is provided, a new Volume of that size will be created.
// If the source parameter is provided, a new Volume that is a copy of the source
// volume will be created.
func resourcePureVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	var v *flasharray.Volume
	var err error

	n, _ := d.GetOk("name")
	s, _ := d.GetOk("source")
	if s.(string) == "" {
		z, _ := d.GetOk("size")
		if v, err = client.Volumes.CreateVolume(n.(string), z.(int)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if v, err = client.Volumes.CopyVolume(n.(string), s.(string), false); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(v.Name)
	return resourcePureVolumeRead(ctx, d, m)
}

// resourcePureVolumeRead sets the values for the given volume ID
func resourcePureVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)

	vol, _ := client.Volumes.GetVolume(d.Id(), nil)

	if vol == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", vol.Name)
	d.Set("size", vol.Size)
	d.Set("serial", vol.Serial)
	d.Set("created", vol.Created)
	d.Set("source", vol.Source)
	return nil
}

// resourcePureVolumeUpdate will update the attributes of the volume.
//
// If a new source is provided, a snapshot of the current volume will be
// taken before the source volume is copied over the current volume. This
// should help protect from any accidental overwrites.
//
// If a new size is provided, it must be larger than the current size.  Only
// extending volumes is supported at this time, since truncating volumes can
// lead to data loss.
func resourcePureVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)

	client := m.(*flasharray.Client)
	var v *flasharray.Volume
	var err error

	if d.HasChange("name") {
		if v, err = client.Volumes.RenameVolume(d.Id(), d.Get("name").(string)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(v.Name)
		d.Set("name", d.Get("name").(string))
	}

	if d.HasChange("source") {
		snapshot, err := client.Volumes.CreateSnapshot(d.Id(), "")
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] Created volume snapshot %s before overwriting volume %s.", snapshot.Name, d.Id())
		if _, err = client.Volumes.CopyVolume(d.Id(), d.Get("source").(string), true); err != nil {
			return diag.FromErr(err)
		}
		d.Set("source", d.Get("source").(string))
	}

	if d.HasChange("size") {
		oldVol, err := client.Volumes.GetVolume(d.Id(), nil)
		z, _ := d.GetOk("size")
		if z.(int) > oldVol.Size {
			if _, err = client.Volumes.ExtendVolume(d.Id(), z.(int)); err != nil {
				return diag.FromErr(err)
			}
		}
		if z.(int) < oldVol.Size {
			return diag.Errorf("error: New size must be larger than current size. Truncating volumes not supported")
		}
	}

	return resourcePureVolumeRead(ctx, d, m)
}

// resourcePureVolumeDelete will delete the volume specified.
// The volume will NOT be eradicated. This is to reduce the chance of
// data loss.  The volume's timer will start for 24 hours, at that time
// the volume will be eradicated.
func resourcePureVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*flasharray.Client)
	_, err := client.Volumes.DeleteVolume(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// resourcePureVolumeImport imports a volume into Terraform.
func resourcePureVolumeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(*flasharray.Client)

	vol, err := client.Volumes.GetVolume(d.Id(), nil)

	if err != nil {
		return nil, err
	}

	d.Set("name", vol.Name)
	d.Set("size", vol.Size)
	d.Set("serial", vol.Serial)
	d.Set("created", vol.Created)
	d.Set("source", vol.Source)
	return []*schema.ResourceData{d}, nil
}
