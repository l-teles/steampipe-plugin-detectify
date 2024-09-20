package detectify

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TABLE DEFINITION

func tableConnector(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "detectify_connector",
		Description: "Table for querying Detectify Connectors configuration.",
		List: &plugin.ListConfig{
			Hydrate: listConnector,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID of the connector."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the connector."},
			{Name: "team_token", Type: proto.ColumnType_STRING, Description: "Team token associated with the connector."},
			{Name: "last_run", Type: proto.ColumnType_JSON, Description: "Details of the last run of the connector."},
			{Name: "provider", Type: proto.ColumnType_STRING, Description: "Provider of the connector."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the connector was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the connector was last updated."},
		},
	}
}

// LIST FUNCTION
func listConnector(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpoint := "/v3/connectors"

	paginatedResponse, err := paginatedResponseV3(ctx, d, endpoint)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_connector.listConnector", "connection_error", err)
		return nil, err
	}

	var allFindings []ConnectorItem

	for _, splitResponse := range paginatedResponse {
		var response ConnectorResponse

		err = json.Unmarshal([]byte(splitResponse), &response)
		if err != nil {
			plugin.Logger(ctx).Error("detectify_connector.listConnector", "failed_unmarshal", err)
			return nil, err
		}

		for _, finding := range response.ConnectorItems {
			d.StreamListItem(ctx, finding)
			allFindings = append(allFindings, finding)
		}
	}

	return allFindings, nil
}

// Custom Structs

// ConnectorItem represents an individual Connector item.
type ConnectorItem struct {
    ID                string `json:"id"`
    Name              string `json:"name"`
    TeamToken         string `json:"team_token"`
	LastRun			  LastRun `json:"last_run"`
    Provider          string `json:"provider"`
    CreatedAt         string `json:"created_at"`
    UpdatedAt         string `json:"updated_at"`
}

// LastRun represents the last run details of the connector.
type LastRun struct {
    Status      string `json:"status"`
    Error       string `json:"error"`
    CompletedAt string `json:"completed_at"`
}

// ConnectorResponse represents the response structure for ConnectorItem items.
type ConnectorResponse struct {
	ConnectorItems []ConnectorItem `json:"items"`
}
