package pipes

import (
	"github.com/turbot/tailpipe-plugin-pipes/sources"
	"github.com/turbot/tailpipe-plugin-pipes/tables"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := plugin.NewPlugin("pipes")

	// register the tables, sources and mappers that we provide
	resources := &plugin.ResourceFunctions{
		Tables:  []func() table.Table{tables.NewAuditLogCollection},
		Sources: []func() row_source.RowSource{sources.NewAuditLogAPISource},
	}

	if err := p.RegisterResources(resources); err != nil {
		return nil, err
	}

	return p, nil
}
