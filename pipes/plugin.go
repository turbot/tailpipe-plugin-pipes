package pipes

import (
	"time"

	"github.com/turbot/tailpipe-plugin-pipes/pipes_source"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_table"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type Plugin struct {
	plugin.PluginBase
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	time.Sleep(10 * time.Second)

	// register the tables, sources and mappers that we provide
	resources := &plugin.ResourceFunctions{
		Tables:  []func() table.Table{pipes_table.NewAuditLogCollection},
		Sources: []func() row_source.RowSource{pipes_source.NewAuditLogAPISource},
	}

	if err := p.RegisterResources(resources); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "pipes"
}
