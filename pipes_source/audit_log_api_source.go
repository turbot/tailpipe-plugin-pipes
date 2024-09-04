package pipes_source

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

func NewAuditLogAPISource() row_source.RowSource {
	return &AuditLogAPISource{}
}

func (s *AuditLogAPISource) Init(ctx context.Context, configData *parse.Data, opts ...row_source.RowSourceOption) error {
	// set the collection state ctor
	s.SetCollectionStateFunc(NewAuditLogAPICollectionState)

	// call base init
	return s.RowSourceBase.Init(ctx, configData, opts...)
}

func (s *AuditLogAPISource) Identifier() string {
	return AuditLogAPISourceIdentifier
}

func (s *AuditLogAPISource) GetConfigSchema() parse.Config {
	return &AuditLogAPISourceConfig{}
}

func (s *AuditLogAPISource) Collect(ctx context.Context) error {
	// NOTE: The API only allows fetching from newest to oldest, so we need to collect in reverse order until we've hit a previously obtain item.
	collectionState := s.CollectionState.(*AuditLogAPICollectionState)
	collectionState.StartCollection() // sets previous state to current state as we manipulate the current state

	var nextToken string

	// Create a default configuration
	configuration := pipes.NewConfiguration()

	// Add your Turbot Pipes user token as an auth header
	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", s.Config.Token))

	// Create a client
	client := pipes.NewAPIClient(configuration)

	orgHandle := s.Config.OrgHandle
	conn := client.GetConfig().Host
	if conn == "" {
		conn = "pipes.turbot.com"
	}
	conn = conn + ":" + orgHandle

	// populate enrichment fields the source is aware of
	// - in this case the connection
	sourceEnrichmentFields := &enrichment.CommonFields{TpIndex: conn, TpSourceType: AuditLogAPISourceIdentifier}

	for {
		listReq := client.Orgs.ListAuditLogs(ctx, orgHandle)
		if nextToken != "" {
			listReq = listReq.NextToken(nextToken)
		}

		slog.Debug("Request with NextToken: ", nextToken)

		listReq = listReq.Limit(100)

		response, _, err := listReq.Execute()
		if err != nil {
			return fmt.Errorf("error obtaining audit logs: %v", err)
		}

		// Checks we have items, and that we have not processed all items previously
		if response.HasItems() {
			items := *response.Items

			for _, item := range items {
				// get time as time opposed to string
				createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
				if err != nil {
					return fmt.Errorf("error parsing created_at field to time.Time: %w", err)
				}

				// check if we've hit previous item - return false if we have, return from function
				if !collectionState.ShouldCollectRow(createdAt, item.Id) {
					return nil
				}
				// populate artifact data
				row := &types.RowData{Data: item, Metadata: sourceEnrichmentFields}

				// update collection state
				collectionState.Upsert(createdAt, item.Id)
				collectionStateJSON, err := s.GetCollectionStateJSON()
				if err != nil {
					return fmt.Errorf("error serialising collectionState data: %w", err)
				}

				if err := s.OnRow(ctx, row, collectionStateJSON); err != nil {
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

	return nil
}
