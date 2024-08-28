package pipes_source

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const AuditLogAPISourceIdentifier = "pipes_audit_log_api"

// AuditLogAPISource source is responsible for collecting audit logs from Turbot Pipes API
type AuditLogAPISource struct {
	row_source.RowSourceBase[*AuditLogAPISourceConfig]
}

func (s *AuditLogAPISource) GetPagingData() (json.RawMessage, error) {
	return nil, nil
}

func NewAuditLogAPISource() row_source.RowSource {
	return &AuditLogAPISource{}
}

func (s *AuditLogAPISource) Identifier() string {
	return AuditLogAPISourceIdentifier
}

func (s *AuditLogAPISource) GetConfigSchema() parse.Config {
	return &AuditLogAPISourceConfig{}
}

func (s *AuditLogAPISource) Collect(ctx context.Context) error {
	// Create a default configuration
	configuration := pipes.NewConfiguration()

	// Add your Turbot Pipes user token as an auth header
	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", s.Config.Token))

	// Create a client
	client := pipes.NewAPIClient(configuration)

	nextToken := ""

	// TODO - fix me
	orgHandle := "turbot-ops"

	conn := client.GetConfig().Host
	if conn == "" {
		conn = "pipes.turbot.com"
	}
	conn = conn + ":" + orgHandle

	// populate enrichment fields the the source is aware of
	// - in this case the connection
	sourceEnrichmentFields := &enrichment.CommonFields{TpIndex: conn}

	for {
		listReq := client.Orgs.ListAuditLogs(ctx, orgHandle)
		if nextToken != "" {
			listReq = listReq.NextToken(nextToken)
		}

		fmt.Println("Request with NextToken: ", nextToken)

		listReq = listReq.Limit(100)

		response, _, err := listReq.Execute()
		if err != nil {
			// Do something with the error
			panic(err)
		}

		if response.HasItems() {

			fmt.Printf("Response item count: %d\n", len(*response.Items))

			// TODO PAGING DATA
			for _, item := range *response.Items {
				// populate artifact data
				row := &types.RowData{Data: item, Metadata: sourceEnrichmentFields}
				if err := s.OnRow(ctx, row, nil); err != nil {
					return fmt.Errorf("error processing row: %w", err)
				}
			}
		}

		if response.HasNextToken() {
			nextToken = *response.NextToken
		} else {
			break
		}

	}

	fmt.Printf("Done!\n")

	return nil

}
