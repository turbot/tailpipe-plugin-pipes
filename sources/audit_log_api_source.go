package sources

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
	"github.com/turbot/tailpipe-plugin-sdk/config_data"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const AuditLogAPISourceIdentifier = "pipes_audit_log_api"

// init function should register the source
func init() {
	row_source.RegisterRowSource[*AuditLogAPISource]()
}

// AuditLogAPISource source is responsible for collecting audit logs from Turbot Pipes API
type AuditLogAPISource struct {
	row_source.RowSourceImpl[*AuditLogAPISourceConfig]
}

func (s *AuditLogAPISource) Init(ctx context.Context, configData config_data.ConfigData, opts ...row_source.RowSourceOption) error {
	// set the collection state ctor
	s.NewCollectionStateFunc = collection_state.NewTimeRangeCollectionState

	// call base init
	return s.RowSourceImpl.Init(ctx, configData, opts...)
}

func (s *AuditLogAPISource) Identifier() string {
	return AuditLogAPISourceIdentifier
}

func (s *AuditLogAPISource) Collect(ctx context.Context) error {
	// NOTE: The API only allows fetching from newest to oldest, so we need to collect in reverse order until we've hit a previously obtained item.
	collectionState := s.CollectionState.(*collection_state.TimeRangeCollectionState[*AuditLogAPISourceConfig])
	// TODO: #config the below should be settable via a config option
	collectionState.IsChronological = false
	collectionState.HasContinuation = false
	// TODO: #collectionState is there a way we can call StartCollection/EndCollection from elsewhere to enforce it?
	collectionState.StartCollection()

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
	sourceName := AuditLogAPISourceIdentifier
	sourceEnrichmentFields := &enrichment.CommonFields{
		TpSourceName:     &sourceName,
		TpSourceType:     AuditLogAPISourceIdentifier,
		TpSourceLocation: &conn,
	}

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

				// check if we've hit previous item - end collection and return if we have
				// TODO: #collectionState this will fill until we hit record in previous state, but what if we have gaps? [incoming data] -> [data]ENDS-HERE -> [gap] -> [data]
				if !collectionState.ShouldCollectRow(createdAt, item.Id) {
					collectionState.EndCollection()
					return nil
				}

				// populate artifact data
				row := &types.RowData{Data: item, Metadata: sourceEnrichmentFields}

				// update collection state
				collectionState.Upsert(createdAt, item.Id, nil)
				collectionStateJSON, err := s.GetCollectionStateJSON()
				if err != nil {
					return fmt.Errorf("error serialising collectionState data: %w", err)
				}

				if err := s.OnRow(ctx, row, collectionStateJSON); err != nil {
					// TODO K #errorHandling - this does not bubble up
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

	collectionState.EndCollection()
	return nil
}
