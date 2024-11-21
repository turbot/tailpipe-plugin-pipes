package tables

type AuditLogTableConfig struct {
}

func (c *AuditLogTableConfig) Validate() error {
	return nil
}

func (c *AuditLogTableConfig) Identifier() string {
	return AuditLogTableIdentifier
}
