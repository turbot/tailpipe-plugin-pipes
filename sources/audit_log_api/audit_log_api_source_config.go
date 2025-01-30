package audit_log_api

import (
	"github.com/hashicorp/hcl/v2"
)

type AuditLogAPISourceConfig struct {
	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`
}

func (c *AuditLogAPISourceConfig) Validate() error {
	return nil
}

func (c *AuditLogAPISourceConfig) Identifier() string {
	return AuditLogAPISourceIdentifier
}
