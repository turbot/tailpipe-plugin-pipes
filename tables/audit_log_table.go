package tables

import (
	"fmt"
	"time"

	"github.com/rs/xid"

	"github.com/turbot/tailpipe-plugin-pipes/config"
	"github.com/turbot/tailpipe-plugin-pipes/rows"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const AuditLogTableIdentifier = "pipes_audit_log"

type AuditLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[*rows.AuditLog, *AuditLogTableConfig, *config.PipesConnection]
}

func NewAuditLogTable() table.Table {
	return &AuditLogTable{}
}

func (c *AuditLogTable) Identifier() string {
	return AuditLogTableIdentifier
}

// GetRowSchema implements Table
// return an instance of the row struct
func (c *AuditLogTable) GetRowSchema() any {
	return rows.AuditLog{}
}

func (c *AuditLogTable) GetConfigSchema() parse.Config {
	return &AuditLogTableConfig{}
}

func (c *AuditLogTable) EnrichRow(row *rows.AuditLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.AuditLog, error) {
	// we expect sourceEnrichmentFields to be set
	if sourceEnrichmentFields == nil {
		return nil, fmt.Errorf("AuditLogTable EnrichRow called with nil sourceEnrichmentFields")
	}
	// we expect name to be set by the Source
	if sourceEnrichmentFields.TpSourceName == "" {
		return nil, fmt.Errorf("AuditLogTable EnrichRow called with TpSourceName unset in sourceEnrichmentFields")
	}

	row.CommonFields = *sourceEnrichmentFields

	// id & Hive fields
	row.TpID = xid.New().String()
	row.TpIndex = row.IdentityHandle
	row.TpDate = row.CreatedAt.Format("2006-01-02")

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
		// TODO: Should row.ProcessId be added to TpAkas?
	}

	row.TpUsernames = append(row.TpUsernames, row.ActorHandle, row.ActorId)

	return row, nil
}
