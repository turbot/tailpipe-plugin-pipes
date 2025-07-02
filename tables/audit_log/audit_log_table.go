package audit_log

import (
	"time"

	"github.com/rs/xid"

	"github.com/turbot/tailpipe-plugin-pipes/sources/audit_log_api"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const AuditLogTableIdentifier = "pipes_audit_log"

type AuditLogTable struct {
}

func (c *AuditLogTable) Identifier() string {
	return AuditLogTableIdentifier
}

func (c *AuditLogTable) GetSourceMetadata() ([]*table.SourceMetadata[*AuditLog], error) {
	return []*table.SourceMetadata[*AuditLog]{
		{
			SourceName: audit_log_api.AuditLogAPISourceIdentifier,
			Mapper:     &AuditLogMapper{},
		},
	}, nil
}

func (c *AuditLogTable) EnrichRow(row *AuditLog, sourceEnrichmentFields schema.SourceEnrichment) (*AuditLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// id & Hive fields
	row.TpID = xid.New().String()
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

func (c *AuditLogTable) GetDescription() string {
	return "Turbot Pipes audit logs detail administrative actions taken."
}
