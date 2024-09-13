package detectify

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TABLE DEFINITION

func tableTechnology(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "detectify_technology",
		Description: "Table for querying Detectify Technologies inventory.",
		List: &plugin.ListConfig{
			Hydrate: listTechnologies,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID of the item."},
			{Name: "asset_id", Type: proto.ColumnType_STRING, Description: "Asset ID associated with the item."},
			{Name: "team_id", Type: proto.ColumnType_STRING, Description: "Team ID associated with the item."},
			{Name: "domain_name", Type: proto.ColumnType_STRING, Description: "Domain name associated with the item."},
			{Name: "service_protocol", Type: proto.ColumnType_STRING, Description: "Service protocol (e.g., https, http)."},
			{Name: "port", Type: proto.ColumnType_INT, Description: "Port number."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the service or technology."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Version of the service or technology."},
			{Name: "categories", Type: proto.ColumnType_JSON, Description: "Categories associated with the service or technology."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Indicates if the item is active."},
			{Name: "first_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the item was first seen."},
			{Name: "disappeared_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the item disappeared."},
		},
	}
}

// LIST FUNCTION
func listTechnologies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpoint := "/v3/technologies"

	paginatedResponse, err := paginatedResponseV3(ctx, d, endpoint)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_technology.listTechnologies", "connection_error", err)
		return nil, err
	}

	var allFindings []TechItem

	for _, splitResponse := range paginatedResponse {
		var response TechResponse

		err = json.Unmarshal([]byte(splitResponse), &response)
		if err != nil {
			plugin.Logger(ctx).Error("detectify_technology.listTechnologies", "failed_unmarshal", err)
			return nil, err
		}

		for _, finding := range response.TechItems {
			d.StreamListItem(ctx, finding)
			allFindings = append(allFindings, finding)
		}
	}

	return allFindings, nil
}

// Custom Structs

// TechItem represents an individual IP item.
type TechItem struct {
	ID              string   `json:"id"`
	AssetID         string   `json:"asset_id"`
	TeamID          string   `json:"team_id"`
	DomainName      string   `json:"domain_name"`
	ServiceProtocol string   `json:"service_protocol"`
	Port            int      `json:"port"`
	Name            string   `json:"name"`
	Version         *string  `json:"version"`
	Categories      []string `json:"categories"`
	Active          bool     `json:"active"`
	FirstSeenAt     string   `json:"first_seen_at"`
	DisappearedAt   *string  `json:"disappeared_at"`
}

// TechResponse represents the response structure for IP items.
type TechResponse struct {
	TechItems []TechItem `json:"items"`
}
