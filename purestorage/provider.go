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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the terraform resource provider called by main.go
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PURE_USERNAME", ""),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PURE_PASSWORD", ""),
			},

			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PURE_APITOKEN", ""),
			},

			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PURE_TARGET", ""),
			},

			"rest_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"verify_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ssl_cert": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"user_agent": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Terraform",
			},

			"request_kwargs": {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  nil,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"purefa_flasharray": dataSourcePureFlashArray(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"purefa_flasharray": schema.DataSourceResourceShim(
				"purefa_flasharray",
				dataSourcePureFlashArray(),
			),
			"purefa_volume":          resourcePureVolume(),
			"purefa_host":            resourcePureHost(),
			"purefa_hostgroup":       resourcePureHostgroup(),
			"purefa_protectiongroup": resourcePureProtectiongroup(),
			"purefa_vgroup":          resourcePureVolumegroup(),
			// "purefa_network_interface": resourcePureNetworkInterface(),
			"purefa_dns_settings":    resourcePureDnsSettings(),
			"purefa_alert_recipient": resourcePureAlertRecipient(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c, err := NewConfig(d)
	if err != nil {
		return nil, err
	}

	return c.Client()
}
