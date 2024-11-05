package rows

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
)

// AuditLog is the struct containing the enriched data for an AuditRecord
type AuditLog struct {
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

func (a *AuditLog) MapFromPipesAuditRecord(record pipes.AuditRecord) error {
	createdAt, err := time.Parse(time.RFC3339, record.CreatedAt)
	if err != nil {
		return fmt.Errorf("error parsing created_at field to time.Time: %w", err)
	}

	s, err := json.Marshal(record.Data)
	if err != nil {
		return fmt.Errorf("error marshalling row data: %w", err)
	}
	dataJsonString := helpers.JSONString(s)

	a.ActionType = record.ActionType
	a.ActorAvatarUrl = record.ActorAvatarUrl
	a.ActorDisplayName = record.ActorDisplayName
	a.ActorHandle = record.ActorHandle
	a.ActorId = record.ActorId
	a.ActorIp = record.ActorIp
	a.CreatedAt = createdAt
	a.Data = &dataJsonString
	a.Id = record.Id
	a.IdentityHandle = record.IdentityHandle
	a.IdentityId = record.IdentityId
	a.ProcessId = record.ProcessId
	a.TargetHandle = record.TargetHandle
	a.TargetId = record.TargetId
	a.TenantId = record.TenantId

	return nil
}
