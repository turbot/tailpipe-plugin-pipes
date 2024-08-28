package pipes

import (
	"time"

	"github.com/turbot/tailpipe-plugin-pipes/pipes_partition"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_source"
	"github.com/turbot/tailpipe-plugin-sdk/partition"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
)

type Plugin struct {
	plugin.PluginBase
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	time.Sleep(10 * time.Second)

	// register the partitions, sources and mappers that we provide
	resources := &plugin.ResourceFunctions{
		Partitions: []func() partition.Partition{pipes_partition.NewAuditLogCollection},
		Sources:    []func() row_source.RowSource{pipes_source.NewAuditLogAPISource},
	}

	if err := p.RegisterResources(resources); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "pipes"
}
