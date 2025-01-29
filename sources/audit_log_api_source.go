package sources

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-pipes/config"
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const AuditLogAPISourceIdentifier = "pipes_audit_log_api"

// init function should register the source
func init() {
	row_source.RegisterRowSource[*AuditLogAPISource]()
}

// AuditLogAPISource source is responsible for collecting audit logs from Turbot Pipes API
type AuditLogAPISource struct {
	row_source.RowSourceImpl[*AuditLogAPISourceConfig, *config.PipesConnection]

	// shadow the collection state to use the reverse order collection state
	CollectionState *collection_state.ReverseOrderCollectionState[*AuditLogAPISourceConfig]
}

func (s *AuditLogAPISource) Init(ctx context.Context, params *row_source.RowSourceParams, opts ...row_source.RowSourceOption) error {
	// set the collection state ctor
	s.NewCollectionStateFunc = collection_state.NewReverseOrderCollectionState

	// call base init
	if err := s.RowSourceImpl.Init(ctx, params, opts...); err != nil {
		return err
	}

	// type assertion to store correctly typed collection state
	s.CollectionState = s.RowSourceImpl.CollectionState.(*collection_state.ReverseOrderCollectionState[*AuditLogAPISourceConfig])
	return nil
}

func (s *AuditLogAPISource) Identifier() string {
	return AuditLogAPISourceIdentifier
}

func (s *AuditLogAPISource) Collect(ctx context.Context) error {
	s.CollectionState.Start()
	defer s.CollectionState.End()

	var nextToken string

	// Create a default configuration
	configuration := pipes.NewConfiguration()

	// Add your Turbot Pipes user token as an auth header
	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", s.Connection.Token))

	// Create a client
	client := pipes.NewAPIClient(configuration)

	orgHandle := s.Connection.OrgHandle
	conn := client.GetConfig().Host
	if conn == "" {
		conn = "pipes.turbot.com"
	}
	conn = conn + ":" + orgHandle

	// populate enrichment fields the source is aware of
	// - in this case the connection
	sourceName := AuditLogAPISourceIdentifier
	sourceEnrichmentFields := &schema.SourceEnrichment{
		CommonFields: schema.CommonFields{
			TpSourceName:     &sourceName,
			TpSourceType:     AuditLogAPISourceIdentifier,
			TpSourceLocation: &conn,
		},
	}

	for {
		listReq := client.Orgs.ListAuditLogs(ctx, orgHandle)
		if nextToken != "" {
			listReq = listReq.NextToken(nextToken)
		}

		slog.Debug("Request with NextToken: ", "next_token", nextToken)

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
				var createdAt time.Time
				createdAt, err = time.Parse(time.RFC3339, item.CreatedAt)
				if err != nil {
					return fmt.Errorf("error parsing created_at field to time.Time: %w", err)
				}

				// check if we should collect this item, if not exit
				if createdAt.Before(s.FromTime) || !s.CollectionState.ShouldCollect(item.Id, createdAt) {
					return nil
				}

				// build a row from item and collect it
				row := &types.RowData{Data: item, SourceEnrichment: sourceEnrichmentFields}
				if err = s.CollectionState.OnCollected(item.Id, createdAt); err != nil {
					return fmt.Errorf("error updating collection state: %w", err)
				}
				if err = s.OnRow(ctx, row); err != nil {
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
