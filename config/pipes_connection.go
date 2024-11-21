package config

import (
	"fmt"
)

const PluginName = "pipes"

type PipesConnection struct {
	Token     string `json:"token" hcl:"token"`
	OrgHandle string `json:"org_handle" hcl:"org_handle"`
}

func (c *PipesConnection) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	if c.OrgHandle == "" {
		return fmt.Errorf("org_handle is required")
	}
	return nil
}

func (c *PipesConnection) Identifier() string {
	return PluginName
}
