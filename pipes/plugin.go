package pipes

import (
	"context"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"log"
	"os"
)

type Plugin struct {
	plugin.Base
}

func (t *Plugin) Identifier() string {
	return "aws"
}

func NewPlugin(_ context.Context) *Plugin {
	return &Plugin{}
}

func (t *Plugin) Collect(req *proto.CollectRequest) error {
	ctx := context.Background()
	log.Println("[INFO] Collect")

	go t.doCollect(ctx, req)

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
	collection.Collect(ctx, onRow)
}
