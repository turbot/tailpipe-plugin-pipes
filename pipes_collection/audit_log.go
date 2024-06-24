package pipes_collection

import (
	"context"
	"encoding/json"
	"fmt"
	helpers "github.com/turbot/tailpipe-plugin-pipes/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipes-sdk-go"
)

type AuditLog struct {
	collection.Base
	Config AuditLogConfig
	source plugin.Source
}

func NewAuditLog(config AuditLogConfig, source plugin.Source) *AuditLog {
	l := &AuditLog{
		Config: config,
		source: source,
	}
	// set the enrich func
	l.EnrichFunc = l.enrichRow

	// add ourselves as an observer to our source
	l.source.AddObserver(l)

	return l
}

func (a *AuditLog) Collect(ctx context.Context, req *proto.CollectRequest) error {
	// tell our source to collect - we will receive row
	return a.source.Collect(ctx, req)
}

func (a *AuditLog) enrichRow(row any, conn string) (any, error) {
	// row must be an AuditRecord
	item, ok := row.(pipes.AuditRecord)
	if !ok {

		return nil, fmt.Errorf("invalid row type")
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
