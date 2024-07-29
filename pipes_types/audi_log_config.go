package pipes_types

import "fmt"

type AuditLogCollectionConfig struct {
	Token string `json:"token"`
}

func (c AuditLogCollectionConfig) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}
