package tables

import (
	"fmt"
	"github.com/turbot/tailpipe-plugin-pipes/config"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-pipes/models"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const AuditLogTableIdentifier = "pipes_audit_log"

type AuditLogTable struct {
	// all tables must embed table.TableBase
	table.TableBase[*AuditLogTableConfig, *config.PipesConnection]
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
	return models.AuditLog{}
}

func (c *AuditLogTable) GetConfigSchema() parse.Config {
	return &AuditLogTableConfig{}
}

func (c *AuditLogTable) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// row should match expected type
	item, ok := row.(models.AuditLog)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected AuditRecord", row)
	}
	// we expect sourceEnrichmentFields to be set
	if sourceEnrichmentFields == nil {
		return nil, fmt.Errorf("AuditLogTable EnrichRow called with nil sourceEnrichmentFields")
	}
	// we expect name to be set by the Source
	if sourceEnrichmentFields.TpSourceName == "" {
		return nil, fmt.Errorf("AuditLogTable EnrichRow called with TpSourceName unset in sourceEnrichmentFields")
	}

	item.CommonFields = *sourceEnrichmentFields

	// id & Hive fields
	item.TpID = xid.New().String()
	item.TpPartition = AuditLogTableIdentifier // TODO: This should be the name from HCL config once passed in!
	item.TpIndex = item.IdentityHandle
	item.TpDate = item.CreatedAt.Format("2006-01-02")

	// Timestamps
	item.TpTimestamp = helpers.UnixMillis(item.CreatedAt.UnixNano() / int64(time.Millisecond))
	item.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

	// Other Enrichment Fields
	if item.ActorIp != "" {
		item.TpSourceIP = &item.ActorIp
		item.TpIps = append(item.TpIps, item.ActorIp)
	}

	if item.TargetId != nil {
		item.TpAkas = append(item.TpAkas, *item.TargetId)
		// TODO: Should item.ProcessId be added to TpAkas?
	}

	item.TpUsernames = append(item.TpUsernames, item.ActorHandle, item.ActorId)

	return item, nil
}
