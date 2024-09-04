package pipes_types

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
)

// AuditLogRow is the struct containing the enriched data for an AuditRecord
type AuditLogRow struct {
	// embed required enrichment fields
	enrichment.CommonFields

	// Additional fields
	ActionType       string              `json:"action_type"`
	ActorAvatarUrl   string              `json:"actor_avatar_url"`
	ActorDisplayName string              `json:"actor_display_name"`
	ActorHandle      string              `json:"actor_handle"`
	ActorId          string              `json:"actor_id"`
	ActorIp          string              `json:"actor_ip"`
	CreatedAt        time.Time           `json:"created_at"`
	Data             *helpers.JSONString `json:"data"`
	Id               string              `json:"id"`
	IdentityHandle   string              `json:"identity_handle"`
	IdentityId       string              `json:"identity_id"`
	ProcessId        *string             `json:"process_id,omitempty"`
	TargetHandle     *string             `json:"target_handle,omitempty"`
	TargetId         *string             `json:"target_id,omitempty"`
	TenantId         string              `json:"tenant_id"`
}
