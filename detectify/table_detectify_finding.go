package detectify

import (
	"context"
	"encoding/json"
	"time"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
    "os"
)


//// TABLE DEFINITION

func tableFinding(_ context.Context) *plugin.Table {
    return &plugin.Table{
        Name:        "detectify_finding",
        Description: "Table for querying Detectify findings data, including status and tags.",
        List: &plugin.ListConfig{
            Hydrate: listFindings,
        },
        Columns: []*plugin.Column{
            {Name: "version", Type: proto.ColumnType_STRING, Description: "Version of the finding."},
            {Name: "uuid", Type: proto.ColumnType_STRING, Description: "Unique ID of this finding."},
            {Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the finding."},
            {Name: "severity", Type: proto.ColumnType_STRING, Description: "Severity of the finding."},
            {Name: "source", Type: proto.ColumnType_JSON, Description: "Source of the finding."},
            {Name: "scan_source", Type: proto.ColumnType_STRING, Description: "Source of the scan."},
            {Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the finding."},
            {Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the finding was last updated."},
            {Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the finding was created."},
            {Name: "asset_token", Type: proto.ColumnType_STRING, Description: "Token of the asset associated with the finding."},
            {Name: "asset", Type: proto.ColumnType_JSON, Description: "Details of the asset associated with the finding."},
            {Name: "location", Type: proto.ColumnType_STRING, Description: "Location of the finding."},
            {Name: "definition", Type: proto.ColumnType_JSON, Description: "Definition of the finding."},
            {Name: "request", Type: proto.ColumnType_JSON, Description: "Request details that triggered the finding."},
            {Name: "response", Type: proto.ColumnType_JSON, Description: "Response details of the finding."},
            {Name: "tags", Type: proto.ColumnType_JSON, Description: "Tags associated with the finding."},
            {Name: "links", Type: proto.ColumnType_JSON, Description: "Links related to the finding."},
            {Name: "cvss_scores", Type: proto.ColumnType_JSON, Description: "CVSS scores of the finding."},
            {Name: "details", Type: proto.ColumnType_JSON, Description: "Detailed information about the finding."},
            {Name: "references", Type: proto.ColumnType_JSON, Description: "References related to the finding."},
            {Name: "cwe", Type: proto.ColumnType_INT, Description: "Common Weakness Enumeration (CWE) identifier."},
            {Name: "host", Type: proto.ColumnType_STRING, Description: "Host associated with the finding."},
        },
    }
}


//// LIST FUNCTION
func listFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
    endpoint := "/v2/vulnerabilities/"

    paginatedResponse, err := paginatedResponse(ctx, d, endpoint)
    if err != nil {
        plugin.Logger(ctx).Error("detectify_finding.listFindings", "connection_error", err)
        return nil, err
    }

    // Write paginatedResponse to a file for debugging
    file, err := os.Create("/Users/luisteles/Downloads/paginated_response1.json")
    if err != nil {
        plugin.Logger(ctx).Error("Failed to create file: %v", err)
        return nil, err
    }
    defer file.Close()

    data, err := json.MarshalIndent(paginatedResponse, "", "  ")
    if err != nil {
        plugin.Logger(ctx).Error("Failed to marshal paginatedResponse: %v", err)
        return nil, err
    }

    if _, err := file.Write(data); err != nil {
        plugin.Logger(ctx).Error("Failed to write to file: %v", err)
        return nil, err
    }

    var allFindings []Finding

    for _, splitResponse := range paginatedResponse {
        var response FindingsResponse

        err = json.Unmarshal([]byte(splitResponse), &response)
        if err != nil {
            plugin.Logger(ctx).Error("detectify_finding.listFindings", "failed_unmarshal", err)
            return nil, err
        }

        for i, finding := range response.Findings {
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
type Finding struct {
    Version    string          `json:"version"`
    UUID       string          `json:"uuid"`
    Title      string          `json:"title"`
    Severity   string          `json:"severity"`
    Source     Source          `json:"source"`
    ScanSource string          `json:"scan_source"`
    Status     string          `json:"status"`
    UpdatedAt  time.Time       `json:"updated_at"`
    CreatedAt  time.Time       `json:"created_at"`
    AssetToken string          `json:"asset_token"`
    Asset      Asset           `json:"asset"`
    Location   string          `json:"location"`
    Definition Definition      `json:"definition"`
    Request    Request         `json:"request"`
    Response   Response        `json:"response"`
    Tags       []Tag           `json:"tags"`
    Links      Links           `json:"links"`
    CvssScores CvssScores      `json:"cvss_scores"`
    Details    Details         `json:"details"`
    References []Reference     `json:"references"`
    CWE        int             `json:"cwe"`
    Host       string          `json:"host"`
}

type Source struct {
    Value string `json:"value"`
}

type Asset struct {
    Token string `json:"token"`
    Name  string `json:"name"`
}

type Definition struct {
    Title         string    `json:"title"`
    Description   string    `json:"description"`
    Risk          string    `json:"risk"`
    ModuleVersion string    `json:"module_version"`
    ModuleRelease time.Time `json:"module_release"`
    IsCrowdsourced bool     `json:"is_crowdsourced"`
}

type Request struct {
    Method  string   `json:"method"`
    URL     string   `json:"url"`
    Headers []Header `json:"headers"`
    Body    string   `json:"body"`
}

type Header struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

type Response struct {
    StatusCode int      `json:"status_code"`
    Headers    []Header `json:"headers"`
    Body       string   `json:"body"`
}

type Tag struct {
    UUID string `json:"uuid"`
    Name string `json:"name"`
}

type Links struct {
    DetailsPage string `json:"details_page"`
}

type CvssScores struct {
    Cvss20 Cvss `json:"cvss_2_0"`
    Cvss30 Cvss `json:"cvss_3_0"`
    Cvss31 Cvss `json:"cvss_3_1"`
}

type Cvss struct {
    Score    float64 `json:"score"`
    Vector   string  `json:"vector"`
    Severity string  `json:"severity"`
}

type Details struct {
    HTML []HTMLDetail `json:"html"`
}

type HTMLDetail struct {
    Value string `json:"value"`
}

type Reference struct {
    UUID   string `json:"uuid"`
    Link   string `json:"link"`
    Name   string `json:"name"`
    Source string `json:"source"`
}

type FindingsResponse struct {
    Findings []Finding `json:"vulnerabilities"`
}

