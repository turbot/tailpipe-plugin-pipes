package pipes

import (
	"context"

	"github.com/turbot/tailpipe-plugin-pipes/pipes_collection"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_source"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"os"
)

type Plugin struct {
	plugin.Base
}

func (t *Plugin) Identifier() string {
	return "aws"
}

func (t *Plugin) Collect(req *proto.CollectRequest) error {

	go t.doCollect(context.Background(), req)

	return nil
}

func (t *Plugin) doCollect(ctx context.Context, req *proto.CollectRequest) {
	// todo config parsing, identify collection type etc.

	// TODO parse config and use to build collection
	//  tactical - create collection
	config := pipes_collection.AuditLogConfig{Token: os.Getenv("PIPES_TOKEN")}
	// TODO source
	var source = pipes_source.NewAuditLogAPISource(config)
	var col = pipes_collection.NewAuditLog(config, source)

	// add ourselves as an observer
	col.AddObserver(t)

	// signal we have started
	t.OnStarted(req)

	// tell the collection to start collecting - this is a blocking call
	err := col.Collect(ctx, req)

	// signal we have completed
	t.OnComplete(req, err)
}
