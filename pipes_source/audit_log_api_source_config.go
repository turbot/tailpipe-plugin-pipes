package pipes_source

import "fmt"

type AuditLogAPISourceConfig struct {
	Token string `json:"token"`
}

func (c *AuditLogAPISourceConfig) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}
