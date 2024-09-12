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
			"detectify_ip":         tableIp(ctx),
			"detectify_technology": tableTechnology(ctx),
			"detectify_port":       tablePort(ctx),
			"detectify_policy":     tablePolicy(ctx),
			"detectify_asset":      tableAsset(ctx),
			"detectify_asset_full": tableAssetFull(ctx),
			"detectify_profile":    tableProfile(ctx),
			"detectify_member":     tableMember(ctx),
		},
	}
	return p
}
