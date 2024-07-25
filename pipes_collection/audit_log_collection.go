package pipes_collection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_source"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_types"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"os"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

type AuditLogCollection struct {
	// all collections must embed collection.Base
	collection.Base

	// the collection config
	Config *pipes_types.AuditLogCollectionConfig
}

func NewAuditLogCollection() plugin.Collection {
	return &AuditLogCollection{}
}

func (c *AuditLogCollection) Identifier() string {
	return "pipes_audit_log"
}

// GetRowStruct implements Collection
// return an instance of the row struct
func (c *AuditLogCollection) GetRowSchema() any {
	return pipes_types.AuditLogRow{}
}

func (c *AuditLogCollection) GetConfigSchema() any {
	return &pipes_types.AuditLogCollectionConfig{}
}

// Init implements Collection
func (c *AuditLogCollection) Init(ctx context.Context, configData []byte) error {
	// TEMP - this will actually parse (or the base will)
	// unmarshal the config
	config := &pipes_types.AuditLogCollectionConfig{
		Token: os.Getenv("PIPES_TOKEN"),
	}
	//err := json.Unmarshal(configData, config)
	//if err != nil {
	//	return fmt.Errorf("error unmarshalling config: %w", err)
	//}

	// todo #config- parse config as hcl
	c.Config = config
	// todo validate config

	// todo create source from config
	source, err := c.getSource(c.Config)
	if err != nil {
		return err
	}
	return c.AddSource(source)
}

func (c *AuditLogCollection) getSource(config *pipes_types.AuditLogCollectionConfig) (plugin.RowSource, error) {
	return pipes_source.NewAuditLogAPISource(config), nil
}

func (c *AuditLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
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
	if sourceEnrichmentFields == nil || sourceEnrichmentFields.TpConnection == "" {
		return nil, fmt.Errorf("Source must provide connection in sourceEnrichmentFields")
	}

	record := &pipes_types.AuditLogRow{
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
	record.TpCollection = "pipes_audit_log"
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
