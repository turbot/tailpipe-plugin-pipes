package pipes

import (
	"github.com/turbot/tailpipe-plugin-pipes/pipes_collection"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

type Plugin struct {
	plugin.Base
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	//time.Sleep(10 * time.Second)
	// register collections which we support
	err := p.RegisterCollections(pipes_collection.NewAuditLogCollection)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "pipes"
}
