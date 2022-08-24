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
	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePureHostgroup() *schema.Resource {
	return &schema.Resource{
		Create: resourcePureHostgroupCreate,
		Read:   resourcePureHostgroupRead,
		Update: resourcePureHostgroupUpdate,
		Delete: resourcePureHostgroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePureHostgroupImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hosts": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Default:  nil,
			},
			"volume": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"lun": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourcePureHostgroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*flasharray.Client)
	var hgroup *flasharray.Hostgroup
	var err error

	var hosts []string
	if h, ok := d.GetOk("hosts"); ok {
		for _, element := range h.([]interface{}) {
			hosts = append(hosts, element.(string))
		}
	}
	data := map[string][]string{"hostlist": hosts}
	if hgroup, err = client.Hostgroups.CreateHostgroup(d.Get("name").(string), data); err != nil {
		return err
	}
	d.SetId(hgroup.Name)
	d.Set("name", d.Get("name").(string))
	d.Set("hosts", hosts)

	var connectedVolumes []string
	if cv, ok := d.GetOk("connected_volumes"); ok {
		for _, element := range cv.([]interface{}) {
			connectedVolumes = append(connectedVolumes, element.(string))
		}
	}

	if cv := d.Get("volume").(*schema.Set).List(); len(cv) > 0 {
		for _, volume := range cv {
			vol, _ := volume.(map[string]interface{})
			data := make(map[string]interface{})
			if vol["lun"] != 0 {
				data["lun"] = vol["lun"].(int)
			}
			if _, err := client.Hostgroups.ConnectHostgroup(hgroup.Name, vol["vol"].(string), data); err != nil {
				return err
			}
		}
	}

	return resourcePureHostgroupRead(d, m)
}

func resourcePureHostgroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*flasharray.Client)

	h, _ := client.Hostgroups.GetHostgroup(d.Id(), nil)

	if h == nil {
		d.SetId("")
		return nil
	}

	if volumes, _ := client.Hostgroups.ListHostgroupConnections(h.Name); volumes != nil {
		if err := d.Set("volume", flattenHgroupVolume(volumes)); err != nil {
			return err
		}
	}

	d.Set("name", h.Name)
	d.Set("hosts", h.Hosts)
	return nil
}

func resourcePureHostgroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*flasharray.Client)
	var hgroup *flasharray.Hostgroup
	var err error

	if d.HasChange("name") {
		if hgroup, err = client.Hostgroups.RenameHostgroup(d.Id(), d.Get("name").(string)); err != nil {
			return err
		}
		d.SetId(hgroup.Name)
		d.Set("name", d.Get("name").(string))
	}

	if d.HasChange("hosts") {
		var hosts []string
		for _, element := range d.Get("hosts").([]interface{}) {
			hosts = append(hosts, element.(string))
		}
		data := map[string][]string{"hostlist": hosts}
		if _, err = client.Hostgroups.SetHostgroup(d.Id(), data); err != nil {
			return err
		}
	}

	if d.HasChange("volume") {
		o, n := d.GetChange("volume")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		disconnectVolumes := os.Difference(ns).List()
		connectVolumes := ns.Difference(os).List()

		if len(connectVolumes) > 0 {
			for _, volume := range connectVolumes {
				data := make(map[string]interface{})
				vol := volume.(map[string]interface{})
				if vol["lun"] != 0 {
					data["lun"] = vol["lun"].(int)
				}
				if _, err = client.Hostgroups.ConnectHostgroup(d.Id(), vol["vol"].(string), data); err != nil {
					return err
				}
			}
		}

		if len(disconnectVolumes) > 0 {
			for _, volume := range disconnectVolumes {
				vol := volume.(map[string]interface{})
				if _, err = client.Hostgroups.DisconnectHostgroup(d.Id(), vol["vol"].(string)); err != nil {
					return err
				}
			}
		}
	}

	return resourcePureHostgroupRead(d, m)
}

func resourcePureHostgroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*flasharray.Client)

	volumes := d.Get("volume").(*schema.Set).List()
	for _, volume := range volumes {
		vol := volume.(map[string]interface{})
		if _, err := client.Hostgroups.DisconnectHostgroup(d.Id(), vol["vol"].(string)); err != nil {
			return err
		}
	}

	var hosts []string
	data := map[string][]string{"hostlist": hosts}
	_, err := client.Hostgroups.SetHostgroup(d.Id(), data)
	if err != nil {
		return err
	}

	_, err = client.Hostgroups.DeleteHostgroup(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourcePureHostgroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(*flasharray.Client)

	h, err := client.Hostgroups.GetHostgroup(d.Id(), nil)

	if err != nil {
		return nil, err
	}

	if volumes, _ := client.Hostgroups.ListHostgroupConnections(h.Name); volumes != nil {
		if err := d.Set("volume", flattenHgroupVolume(volumes)); err != nil {
			return nil, err
		}
	}

	d.Set("name", h.Name)
	d.Set("hosts", h.Hosts)
	return []*schema.ResourceData{d}, nil
}
