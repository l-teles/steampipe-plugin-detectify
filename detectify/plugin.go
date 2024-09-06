package detectify

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-detectify", 
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"detectify_finding":    tableFinding(ctx),
			"detectify_ips":        tableIps(ctx),
			"detectify_technologies":    tableTechnologies(ctx),
			"detectify_ports":    tablePorts(ctx),
			"detectify_breaches":    tableBreaches(ctx),
			"detectify_assets":    tableAssets(ctx),
			"detectify_profiles":    tableProfiles(ctx),
			"detectify_members":    tableMembers(ctx),
		},
	}
	return p
}
