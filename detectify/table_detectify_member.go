package detectify

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TABLE DEFINITION

func tableMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "detectify_member",
		Description: "Table for querying Detectify users data.",
		List: &plugin.ListConfig{
			Hydrate: listMembers,
		},
		Columns: []*plugin.Column{
			{Name: "user_token", Type: proto.ColumnType_STRING, Description: "Unique token of the user."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "First name of the user."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "Last name of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the user."},
			{Name: "authentication", Type: proto.ColumnType_STRING, Description: "Authentication method used by the user."},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "Role of the user."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the user was created."},
			{Name: "last_login", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the user last logged in."},
		},
	}
}

// LIST FUNCTION
func listMembers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpoint := "/v2/members/"

	paginatedResponse, err := connect(ctx, d, endpoint, nil)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_member.listMembers", "connection_error", err)
		return nil, err
	}

	var allFindings []MemberItem

	// Unmarshal the paginated response directly into a slice of ProfileItem
	err = json.Unmarshal([]byte(paginatedResponse), &allFindings)
	if err != nil {
		plugin.Logger(ctx).Error("detectify_member.listMembers", "failed_unmarshal", err)
		return nil, err
	}

	// Stream each finding
	for _, finding := range allFindings {
		d.StreamListItem(ctx, finding)
	}

	return allFindings, nil
}

// Custom Structs
type MemberItem struct {
	UserToken      string `json:"user_token"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Authentication string `json:"authentication"`
	Role           string `json:"role"`
	Created        string `json:"created"`
	LastLogin      string `json:"last_login"`
}
