package detectify

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TABLE DEFINITION

func tableIp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "detectify_ip",
		Description: "Table for querying Detectify IPs inventory.",
		List: &plugin.ListConfig{
			Hydrate: listIps,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID of the IP item."},
			{Name: "ip_address", Type: proto.ColumnType_STRING, Description: "IP address."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Indicates if the IP is active."},
			{Name: "enriched", Type: proto.ColumnType_BOOL, Description: "Indicates if the IP is enriched."},
			{Name: "domain_name", Type: proto.ColumnType_STRING, Description: "Domain name associated with the IP."},
			{Name: "asset_id", Type: proto.ColumnType_STRING, Description: "Asset ID associated with the IP."},
			{Name: "team_id", Type: proto.ColumnType_STRING, Description: "Team ID associated with the IP."},
			{Name: "ip_version", Type: proto.ColumnType_STRING, Description: "IP version (e.g., IPv4, IPv6)."},
			{Name: "first_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the IP was first seen."},
			{Name: "disappeared_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the IP disappeared."},
			{Name: "autonomous_system", Type: proto.ColumnType_JSON, Description: "Detailed information about the autonomous system."},
			{Name: "geolocation", Type: proto.ColumnType_JSON, Description: "Detailed information about the geolocation."},
		},
	}
}

// LIST FUNCTION
func listIps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpoint := "/v3/ips"

	paginatedResponse, err := paginatedResponseV3(ctx, d, endpoint)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_ip.listIps", "connection_error", err)
		return nil, err
	}

	var allFindings []IPItem

	for _, splitResponse := range paginatedResponse {
		var response IPResponse

		err = json.Unmarshal([]byte(splitResponse), &response)
		if err != nil {
			plugin.Logger(ctx).Error("detectify_ip.listIps", "failed_unmarshal", err)
			return nil, err
		}

		for _, finding := range response.IPItems {
			d.StreamListItem(ctx, finding)
			allFindings = append(allFindings, finding)
		}
	}

	return allFindings, nil
}

// // Custom Structs
// AutonomousSystem represents the autonomous system information.
type AutonomousSystem struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Number int    `json:"number"`
}

// Geolocation represents the geolocation information.
type Geolocation struct {
	Continent     string `json:"continent"`
	ContinentName string `json:"continent_name"`
	Country       string `json:"country"`
	CountryName   string `json:"country_name"`
}

// IPItem represents an individual IP item.
type IPItem struct {
	ID               string           `json:"id"`
	IPAddress        string           `json:"ip_address"`
	Active           bool             `json:"active"`
	Enriched         bool             `json:"enriched"`
	DomainName       string           `json:"domain_name"`
	AssetID          string           `json:"asset_id"`
	TeamID           string           `json:"team_id"`
	IPVersion        string           `json:"ip_version"`
	FirstSeenAt      time.Time        `json:"first_seen_at"`
	DisappearedAt    time.Time        `json:"disappeared_at"`
	AutonomousSystem AutonomousSystem `json:"autonomous_system"`
	Geolocation      Geolocation      `json:"geolocation"`
}

// IPResponse represents the response structure for IP items.
type IPResponse struct {
	IPItems []IPItem `json:"items"`
}