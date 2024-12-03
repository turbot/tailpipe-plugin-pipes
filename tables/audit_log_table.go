package tables

import (
	"fmt"
	"time"

	"github.com/rs/xid"

	"github.com/turbot/tailpipe-plugin-pipes/mappers"
	"github.com/turbot/tailpipe-plugin-pipes/rows"
	"github.com/turbot/tailpipe-plugin-pipes/sources"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// init registers the table
func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.AuditLog, *AuditLogTableConfig, *AuditLogTable]()
}

const AuditLogTableIdentifier = "pipes_audit_log"

type AuditLogTable struct {
}

func (c *AuditLogTable) Identifier() string {
	return AuditLogTableIdentifier
}

func (c *AuditLogTable) GetSourceMetadata(_ *AuditLogTableConfig) []*table.SourceMetadata[*rows.AuditLog] {
	return []*table.SourceMetadata[*rows.AuditLog]{
		{
			SourceName: sources.AuditLogAPISourceIdentifier,
			Mapper:     &mappers.AuditLogMapper{},
		},
	}
}

func (c *AuditLogTable) EnrichRow(row *rows.AuditLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.AuditLog, error) {
	// we expect sourceEnrichmentFields to be set
	if sourceEnrichmentFields == nil {
		return nil, fmt.Errorf("AuditLogTable EnrichRow called with nil sourceEnrichmentFields")
	}
	// we expect name to be set by the Source
	if sourceEnrichmentFields.TpSourceName == nil {
		return nil, fmt.Errorf("AuditLogTable EnrichRow called with TpSourceName unset in sourceEnrichmentFields")
	}

	row.CommonFields = *sourceEnrichmentFields

	// id & Hive fields
	row.TpID = xid.New().String()
	row.TpIndex = row.IdentityHandle
	row.TpDate = row.CreatedAt.Truncate(24 * time.Hour)

	// Timestamps
	row.TpTimestamp = row.CreatedAt
	row.TpIngestTimestamp = time.Now()

	// Other Enrichment Fields
	if row.ActorIp != "" {
		row.TpSourceIP = &row.ActorIp
		row.TpIps = append(row.TpIps, row.ActorIp)
	}

	if row.TargetId != nil {
		row.TpAkas = append(row.TpAkas, *row.TargetId)
	}

	row.TpUsernames = append(row.TpUsernames, row.ActorHandle, row.ActorId)

	return row, nil
}
