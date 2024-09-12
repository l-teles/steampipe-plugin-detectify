package detectify

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TABLE DEFINITION

func tableProfile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "detectify_profile",
		Description: "Table for querying Detectify scanning profiles data.",
		List: &plugin.ListConfig{
			Hydrate: listProfiles,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the profile."},
			{Name: "endpoint", Type: proto.ColumnType_STRING, Description: "Endpoint of the profile."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the profile was created."},
			{Name: "token", Type: proto.ColumnType_STRING, Description: "Token associated with the profile."},
			{Name: "latest_scan", Type: proto.ColumnType_JSON, Description: "Timestamp when the latest scan started."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the profile."},
		},
	}
}

// LIST FUNCTION
func listProfiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpoint := "/v2/profiles/"

	paginatedResponse, err := connect(ctx, d, endpoint, nil)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_profile.listProfiles", "connection_error", err)
		return nil, err
	}

	var allFindings []ProfileItem

	// Unmarshal the paginated response directly into a slice of ProfileItem
	err = json.Unmarshal([]byte(paginatedResponse), &allFindings)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_profile.listProfiles", "failed_unmarshal", err)
		return nil, err
	}

	// Stream each finding
	for _, finding := range allFindings {
		d.StreamListItem(ctx, finding)
	}

	return allFindings, nil
}

// Custom Structs

// PortItem represents an individual profile item.
type ProfileItem struct {
	Name       string     `json:"name"`
	Endpoint   string     `json:"endpoint"`
	Created    string     `json:"created"`
	Token      string     `json:"token"`
	LatestScan LatestScan `json:"latest_scan"`
	Status     string     `json:"status"`
}

// LatestScan represents the  structure for LatestScan items.
type LatestScan struct {
	Started string `json:"started"`
	Ended   string `json:"ended"`
	Status  string `json:"status"`
}