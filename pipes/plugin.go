package pipes

import (
	"github.com/turbot/tailpipe-plugin-pipes/pipes_collection"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_source"
	"github.com/turbot/tailpipe-plugin-sdk/interfaces"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"time"
)

type Plugin struct {
	plugin.Base
}

func NewPlugin() (interfaces.TailpipePlugin, error) {
	p := &Plugin{}

	time.Sleep(10 * time.Second)
	// register collections which we support
	err := p.RegisterCollections(pipes_collection.NewAuditLogCollection)
	if err != nil {
		return nil, err
	}

	p.RegisterSources(pipes_source.NewAuditLogAPISource)

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "pipes"
}
