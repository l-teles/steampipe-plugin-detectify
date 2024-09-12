package detectify

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TABLE DEFINITION

func tablePolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "detectify_policy",
		Description: "Table for querying Detectify Policies inventory.",
		List: &plugin.ListConfig{
			Hydrate: listPolicy,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID of the item."},
			{Name: "policy_id", Type: proto.ColumnType_STRING, Description: "ID of the policy."},
			{Name: "policy_name", Type: proto.ColumnType_STRING, Description: "Name of the policy."},
			{Name: "asset_id", Type: proto.ColumnType_STRING, Description: "ID of the asset."},
			{Name: "asset_name", Type: proto.ColumnType_STRING, Description: "Name of the asset."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "Severity level."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Indicates if the item is active."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the item."},
			{Name: "status_updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the status was last updated."},
			{Name: "first_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the item was first seen."},
			{Name: "disappeared_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the item disappeared."},
		},
	}
}

// LIST FUNCTION
func listPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpoint := "/v3/breaches"

	paginatedResponse, err := paginatedResponseV3(ctx, d, endpoint)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_policy.listPolicy", "connection_error", err)
		return nil, err
	}

	var allFindings []PolicyItem

	for _, splitResponse := range paginatedResponse {
		var response PolicyResponse

		err = json.Unmarshal([]byte(splitResponse), &response)
		if err != nil {
			plugin.Logger(ctx).Error("detectify_policy.listPolicy", "failed_unmarshal", err)
			return nil, err
		}

		for _, finding := range response.PolicyItems {
			d.StreamListItem(ctx, finding)
			allFindings = append(allFindings, finding)
		}
	}

	return allFindings, nil
}

// Custom Structs

// PolicyItem represents an individual Policy item.
type PolicyItem struct {
	ID              string  `json:"id"`
	PolicyID        string  `json:"policy_id"`
	PolicyName      string  `json:"policy_name"`
	AssetID         string  `json:"asset_id"`
	AssetName       string  `json:"asset_name"`
	Severity        string  `json:"severity"`
	Active          bool    `json:"active"`
	Status          string  `json:"status"`
	StatusUpdatedAt string  `json:"status_updated_at"`
	FirstSeenAt     string  `json:"first_seen_at"`
	DisappearedAt   *string `json:"disappeared_at"`
}

// PolicyResponse represents the response structure for Policy items.
type PolicyResponse struct {
	PolicyItems []PolicyItem `json:"items"`
}
