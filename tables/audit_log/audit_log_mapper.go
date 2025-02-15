package audit_log

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/table"

	"github.com/turbot/pipes-sdk-go"
)

type AuditLogMapper struct {
}

func (m *AuditLogMapper) Map(_ context.Context, data any, _ ...table.MapOption[*AuditLog]) (*AuditLog, error) {
	input, ok := data.(pipes.AuditRecord)
	if !ok {
		return nil, fmt.Errorf("expected pipes.AuditRecord, got %T", data)
	}

	auditLog := &AuditLog{}

	createdAt, err := time.Parse(time.RFC3339, input.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at field to time.Time: %w", err)
	}

	auditLog.ActionType = input.ActionType
	auditLog.ActorAvatarUrl = input.ActorAvatarUrl
	auditLog.ActorDisplayName = input.ActorDisplayName
	auditLog.ActorHandle = input.ActorHandle
	auditLog.ActorId = input.ActorId
	auditLog.ActorIp = input.ActorIp
	auditLog.CreatedAt = createdAt
	auditLog.Data = input.Data
	auditLog.Id = input.Id
	auditLog.IdentityHandle = input.IdentityHandle
	auditLog.IdentityId = input.IdentityId
	auditLog.ProcessId = input.ProcessId
	auditLog.TargetHandle = input.TargetHandle
	auditLog.TargetId = input.TargetId
	auditLog.TenantId = input.TenantId

	return auditLog, nil
}

func (m *AuditLogMapper) Identifier() string {
	return "pipes_audit_log_mapper"
}
