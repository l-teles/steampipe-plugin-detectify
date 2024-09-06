package detectify

import (
	"context"
	"encoding/json"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)


//// TABLE DEFINITION

func tablePorts(_ context.Context) *plugin.Table {
    return &plugin.Table{
        Name:        "detectify_ports",
        Description: "Table for querying Detectify Port inventory.",
        List: &plugin.ListConfig{
            Hydrate: listPorts,
        },
        Columns: []*plugin.Column{
            {Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID of the port."},
            {Name: "team_id", Type: proto.ColumnType_STRING, Description: "ID of the team."},
            {Name: "asset_id", Type: proto.ColumnType_STRING, Description: "ID of the asset."},
            {Name: "domain_name", Type: proto.ColumnType_STRING, Description: "Domain name associated with the port."},
            {Name: "ip_address", Type: proto.ColumnType_STRING, Description: "IP address."},
            {Name: "port", Type: proto.ColumnType_INT, Description: "Port number."},
            {Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the port."},
            {Name: "first_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the port was first seen."},
            {Name: "disappeared_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the port disappeared."},        
    },
    }
}



//// LIST FUNCTION
func listPorts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
    endpoint := "/v3/ports"

    paginatedResponse, err := paginatedResponseV3(ctx, d, endpoint)
    if err != nil {
        plugin.Logger(ctx).Error("detectify_ports.listPorts", "connection_error", err)
        return nil, err
    }

    var allFindings []PortItem

    for _, splitResponse := range paginatedResponse {
        var response PortResponse

        err = json.Unmarshal([]byte(splitResponse), &response)
        if err != nil {
            plugin.Logger(ctx).Error("detectify_ports.listPorts", "failed_unmarshal", err)
            return nil, err
        }

        for i, finding := range response.PortItems {
            d.StreamListItem(ctx, finding)
            allFindings = append(allFindings, finding)
        
            // Convert finding to JSON (optional, for debugging purposes)
            findingData, err := json.MarshalIndent(finding, "", "  ")
            if err != nil {
                plugin.Logger(ctx).Error("Failed to marshal finding: %v", err)
                return nil, err
            }
        
            // Optional: Log the finding data for debugging
            plugin.Logger(ctx).Info("Finding data", "index", i, "data", string(findingData))
        }
    }

    return allFindings, nil
}


//// Custom Structs

// PortItem represents an individual port item.
type PortItem struct {
    ID            string  `json:"id"`
    TeamID        string  `json:"team_id"`
    AssetID       string  `json:"asset_id"`
    DomainName    string  `json:"domain_name"`
    IPAddress     string  `json:"ip_address"`
    Port          int     `json:"port"`
    Status        string  `json:"status"`
    FirstSeenAt   string  `json:"first_seen_at"`
    DisappearedAt *string `json:"disappeared_at"`
}

// PortResponse represents the response structure for PortItem items.
type PortResponse struct {
    PortItems []PortItem `json:"items"`
}