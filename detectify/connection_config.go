package detectify

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type detectifyConfig struct {
	BaseUrl *string `cty:"base_url"`
	Token   *string `cty:"token"`
	Secret  *string `cty:"secret"`
	Tokenv3 *string `cty:"tokenv3"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"base_url": {
		Type: schema.TypeString,
	},
	"token": {
		Type: schema.TypeString,
	},
	"secret": {
		Type: schema.TypeString,
	},
	"tokenv3": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &detectifyConfig{}
}

func GetConfig(connection *plugin.Connection) detectifyConfig {
	if connection == nil || connection.Config == nil {
		return detectifyConfig{}
	}
	config, _ := connection.Config.(detectifyConfig)
	return config
}
