package mappers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-pipes/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type AuditLogMapper struct {
}

func NewAuditLogMapper() table.Mapper[*rows.AuditLog] {
	return &AuditLogMapper{}
}

func (m *AuditLogMapper) Map(_ context.Context, data any) ([]*rows.AuditLog, error) {
	input, ok := data.(pipes.AuditRecord)
	if !ok {
		return nil, fmt.Errorf("expected pipes.AuditRecord, got %T", data)
	}

	auditLog := &rows.AuditLog{}

	createdAt, err := time.Parse(time.RFC3339, input.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at field to time.Time: %w", err)
	}

	s, err := json.Marshal(input.Data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling row data: %w", err)
	}
	dataJsonString := types.JSONString(s)

	auditLog.ActionType = input.ActionType
	auditLog.ActorAvatarUrl = input.ActorAvatarUrl
	auditLog.ActorDisplayName = input.ActorDisplayName
	auditLog.ActorHandle = input.ActorHandle
	auditLog.ActorId = input.ActorId
	auditLog.ActorIp = input.ActorIp
	auditLog.CreatedAt = createdAt
	auditLog.Data = &dataJsonString
	auditLog.Id = input.Id
	auditLog.IdentityHandle = input.IdentityHandle
	auditLog.IdentityId = input.IdentityId
	auditLog.ProcessId = input.ProcessId
	auditLog.TargetHandle = input.TargetHandle
	auditLog.TargetId = input.TargetId
	auditLog.TenantId = input.TenantId

	return []*rows.AuditLog{auditLog}, nil
}

func (m *AuditLogMapper) Identifier() string {
	return "pipes_audit_log_mapper"
}
