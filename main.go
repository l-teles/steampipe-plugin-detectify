package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/l-teles/steampipe-plugin-detectify/detectify"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: detectify.Plugin})
}
