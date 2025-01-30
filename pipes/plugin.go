package pipes

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-pipes/config"
	"github.com/turbot/tailpipe-plugin-pipes/sources/audit_log_api"
	"github.com/turbot/tailpipe-plugin-pipes/tables/audit_log"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type Plugin struct {
	plugin.PluginImpl
}

func init() {
	// Register tables, with type parameters:
	// 1. row struct
	// 2. table implementation
	table.RegisterTable[*audit_log.AuditLog, *audit_log.AuditLogTable]()

	// register sources
	row_source.RegisterRowSource[*audit_log_api.AuditLogAPISource]()
}

func NewPlugin() (_ plugin.TailpipePlugin, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = helpers.ToError(r)
		}
	}()

	p := &Plugin{
		PluginImpl: plugin.NewPluginImpl(config.PluginName),
	}

	return p, nil
}
