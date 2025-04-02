package audit_log

import (
	"context"
	"log/slog"
	"time"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-sdk/error_types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type AuditLogMapper struct {
}

func (m *AuditLogMapper) Map(_ context.Context, data any, _ ...mappers.MapOption[*AuditLog]) (*AuditLog, error) {
	input, ok := data.(pipes.AuditRecord)
	if !ok {
		slog.Error("unable to map audit log record: expected pipes.AuditRecord, got %T", data)
		return nil, error_types.NewRowErrorWithMessage("unable to map row, invalid type received")
	}

	auditLog := &AuditLog{}

	createdAt, err := time.Parse(time.RFC3339, input.CreatedAt)
	if err != nil {
		slog.Error("audit log %s failed mapping created_at field to time.Time: %v", input.Id, err)
		return nil, error_types.NewRowErrorWithFields([]string{}, []string{"created_at"})
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
	auditLog.IdentityHandle = types.SafeString(input.IdentityHandle)
	auditLog.IdentityId = types.SafeString(input.IdentityId)
	auditLog.ProcessId = input.ProcessId
	auditLog.TargetHandle = input.TargetHandle
	auditLog.TargetId = input.TargetId
	auditLog.TenantId = input.TenantId

	return auditLog, nil
}

func (m *AuditLogMapper) Identifier() string {
	return "pipes_audit_log_mapper"
}
