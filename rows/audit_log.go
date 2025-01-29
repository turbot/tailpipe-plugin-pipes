package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

// AuditLog is the struct containing the enriched data for an AuditRecord
type AuditLog struct {
	// embed required enrichment fields
	schema.CommonFields

	// Additional fields
	ActionType       string                 `json:"action_type"`
	ActorAvatarUrl   string                 `json:"actor_avatar_url"`
	ActorDisplayName string                 `json:"actor_display_name"`
	ActorHandle      string                 `json:"actor_handle"`
	ActorId          string                 `json:"actor_id"`
	ActorIp          string                 `json:"actor_ip"`
	CreatedAt        time.Time              `json:"created_at"`
	Data             map[string]interface{} `json:"data" parquet:"type=JSON"`
	Id               string                 `json:"id"`
	IdentityHandle   string                 `json:"identity_handle"`
	IdentityId       string                 `json:"identity_id"`
	ProcessId        *string                `json:"process_id,omitempty"`
	TargetHandle     *string                `json:"target_handle,omitempty"`
	TargetId         *string                `json:"target_id,omitempty"`
	TenantId         string                 `json:"tenant_id"`
}

func (a *AuditLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"action_type":        "The type of action performed in the audit log, such as 'workspace.create', 'workspace.reboot', or 'workspace.delete'.",
		"actor_avatar_url":   "The URL of the avatar image associated with the actor who performed the action.",
		"actor_display_name": "The display name of the actor who performed the action.",
		"actor_handle":       "The unique handle or username associated with the actor.",
		"actor_id":           "The unique identifier of the actor who initiated the action.",
		"actor_ip":           "The IP address from which the action was performed.",
		"created_at":         "The date and time when the audit log entry was created, in ISO 8601 format.",
		"data":               "Additional metadata related to the action, in JSON format.",
		"id":                 "The unique identifier of the audit log entry.",
		"identity_handle":    "The handle associated with the identity that was affected by the action.",
		"identity_id":        "The unique identifier of the identity that was affected by the action.",
		"process_id":         "The unique identifier of the process that triggered the action, if applicable.",
		"target_handle":      "The handle of the target resource that was acted upon.",
		"target_id":          "The unique identifier of the target resource that was acted upon.",
		"tenant_id":          "The unique identifier of the tenant or organization associated with the action.",

		// Override table specific tp_* column descriptions
		"tp_index": "The org or user handle the logs were collected from.",
	}
}
