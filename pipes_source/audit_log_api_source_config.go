package pipes_source

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type AuditLogAPISourceConfig struct {
	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Token     string `json:"token" hcl:"token"`
	OrgHandle string `json:"org_handle" hcl:"org_handle"`
}

func (c *AuditLogAPISourceConfig) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}
