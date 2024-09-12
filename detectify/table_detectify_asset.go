package detectify

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAsset(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "detectify_asset",
		Description: "Table for querying Detectify root domain (asset) data. Subdomains not included.",
		List: &plugin.ListConfig{
			Hydrate: listAssets,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the asset."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the asset."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the asset was created."},
			{Name: "updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the asset was last updated."},
			{Name: "discovered", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the asset was discovered."},
			{Name: "last_seen", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the asset was last seen."},
			{Name: "token", Type: proto.ColumnType_STRING, Description: "Token associated with the asset."},
			{Name: "monitored", Type: proto.ColumnType_BOOL, Description: "Indicates if the asset is monitored."},
			{Name: "added_by", Type: proto.ColumnType_JSON, Description: "List of sources that added the asset."},
		},
	}
}

// // LIST FUNCTION
func listAssets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpoint := "/v2/assets/"

	paginatedResponse, err := paginatedResponse(ctx, d, endpoint)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_asset.listAssets", "connection_error", err)
		return nil, err
	}

	var allFindings []AssetItem

	for _, splitResponse := range paginatedResponse {
		var response AssetsResponse

		err = json.Unmarshal([]byte(splitResponse), &response)
		if err != nil {
			plugin.Logger(ctx).Error("detectify_asset.listAssets", "failed_unmarshal", err)
			return nil, err
		}

		for _, finding := range response.Assets {
			d.StreamListItem(ctx, finding)
			allFindings = append(allFindings, finding)
		}
	}

	return allFindings, nil
}

// // Custom Structs
type AssetItem struct {
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
	Discovered string   `json:"discovered"`
	LastSeen   string   `json:"last_seen"`
	Token      string   `json:"token"`
	Monitored  bool     `json:"monitored"`
	AddedBy    []string `json:"added_by"`
}

type AssetsResponse struct {
	Assets []AssetItem `json:"assets"`
}