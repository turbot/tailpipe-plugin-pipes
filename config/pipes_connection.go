package config

import "github.com/turbot/tailpipe-plugin-sdk/parse"

type PipesConnection struct {
}

func NewPipesConnection() parse.Config {
	return &PipesConnection{}
}

func (c *PipesConnection) Validate() error {
	return nil
}
