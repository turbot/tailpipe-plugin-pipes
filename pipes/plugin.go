package pipes

import (
	"context"
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
	collection := &PipesAuditLogCollection{
		Config: PipesAuditLogCollectionConfig{
			Token: os.Getenv("PIPES_TOKEN"),
		},
	}
	onRow := func(row any) {
		t.OnRow(row, req)
	}

	t.OnStarted(req)

	err := collection.Collect(ctx, onRow)

	t.OnComplete(req, err)

}
