package tables

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-pipes/models"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const AuditLogTableIdentifier = "pipes_audit_log"

type AuditLogTable struct {
	// all tables must embed table.TableBase
	table.TableBase[*AuditLogTableConfig]
}

func NewAuditLogCollection() table.Table {
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
	// row must be an AuditRecord
	item, ok := row.(pipes.AuditRecord)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected AuditRecord", row)
	}
	// we expect sourceEnrichmentFields to be set
	if sourceEnrichmentFields == nil {
		return nil, fmt.Errorf("AuditLogTable EnrichRow called with nil sourceEnrichmentFields")
	}
	// we expect connection to be set by the Source
	if sourceEnrichmentFields == nil || sourceEnrichmentFields.TpIndex == "" {
		return nil, fmt.Errorf("Source must provide connection in sourceEnrichmentFields")
	}

	record := &models.AuditLog{
		CommonFields: *sourceEnrichmentFields,
	}

	// Record standardization
	record.TpID = xid.New().String()
	record.TpSourceType = "pipes_audit_log"
	tpTimestamp, err := time.Parse(time.RFC3339, item.CreatedAt)
	if err != nil {
		return nil, err
	}
	record.TpTimestamp = helpers.UnixMillis(tpTimestamp.UnixNano() / int64(time.Millisecond))
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
	if record.ActorIp != "" {
		record.TpSourceIP = &item.ActorIp
		record.TpIps = append(record.TpIps, item.ActorIp)
	}
	if item.TargetId != nil {
		record.TpAkas = append(record.TpAkas, *item.TargetId)
	}
	record.TpUsernames = append(record.TpUsernames, item.ActorHandle)

	// Set hive fields
	record.TpPartition = "pipes_audit_log"
	record.TpDate = tpTimestamp.Format("2006-01-02")

	// Record data
	record.ActionType = item.ActionType
	record.ActorAvatarUrl = item.ActorAvatarUrl
	record.ActorDisplayName = item.ActorDisplayName
	record.ActorHandle = item.ActorHandle
	record.ActorId = item.ActorId
	record.ActorIp = item.ActorIp
	record.CreatedAt = tpTimestamp
	s, err := json.Marshal(item.Data)
	if err != nil {
		panic(err)
	}
	js := helpers.JSONString(s)
	record.Data = &js
	record.Id = item.Id
	record.IdentityHandle = item.IdentityHandle
	record.IdentityId = item.IdentityId
	record.ProcessId = item.ProcessId
	record.TargetHandle = item.TargetHandle
	record.TargetId = item.TargetId
	record.TenantId = item.TenantId

	return record, nil
}
