package main

import (
	"github.com/turbot/tailpipe-plugin-pipes/pipes"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: pipes.NewPlugin,
	})
}
