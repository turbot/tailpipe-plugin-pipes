package pipes_collection

import (
	"encoding/json"
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

type AuditLogCollection struct {
	// all collections must embed collection.Base
	// this add observer and enrich functions
	collection.Base

	// the collection config
	Config AuditLogCollectionConfig
}

func NewAuditLogCollection(config AuditLogCollectionConfig, source plugin.Source) *AuditLogCollection {
	l := &AuditLogCollection{
		Config: config,
	}
	// Init sets the Source property on Base and adds us as an observer to it
	// It also sets us as the Enricher property on Base
	l.Base.Init(source, l)

	return l
}

func (a *AuditLogCollection) Identifier() string {
	return "pipes_audit_log"
}

func (a *AuditLogCollection) EnrichRow(row any, sourceEnrichmentFields map[string]any) (any, error) {
	// row must be an AuditRecord
	item, ok := row.(pipes.AuditRecord)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected AuditRecord", row)
	}
	// we expect sourceEnrichmentFields to be set
	if sourceEnrichmentFields == nil {
		return nil, fmt.Errorf("AuditLogCollection EnrichRow called with nil sourceEnrichmentFields")
	}
	// we expect connection to be set by the Source
	conn, ok := sourceEnrichmentFields[constants.TpConnection].(string)
	if !ok {
		return nil, fmt.Errorf("Source must provide connection in sourceEnrichmentFields")
	}

	record := &AuditLogRow{}

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
	record.TpCollection = "pipes_audit_log"
	record.TpConnection = conn
	record.TpYear = int32(tpTimestamp.Year())
	record.TpMonth = int32(tpTimestamp.Month())
	record.TpDay = int32(tpTimestamp.Day())

	// Record data
	record.ActionType = item.ActionType
	record.ActorAvatarUrl = item.ActorAvatarUrl
	record.ActorDisplayName = item.ActorDisplayName
	record.ActorHandle = item.ActorHandle
	record.ActorId = item.ActorId
	record.ActorIp = item.ActorIp
	record.CreatedAt = record.TpTimestamp
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
