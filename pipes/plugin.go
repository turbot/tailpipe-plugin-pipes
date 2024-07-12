package pipes

import (
	"context"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_collection"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_types"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"os"
)

type Plugin struct {
	plugin.Base
}

// ctor
//
//	func NewPlugin(_ context.Context) *Plugin {
//		return &Plugin{}
//	}
func (t *Plugin) Identifier() string {
	return "pipes"
}

//// GetSchema returns the schema (i.e. an instance of the row struct) for all collections
//// it is used primarily to validate the row structs provide the required fields
//func (t *Plugin) GetSchema(collection string) map[string]any {
//	collections := []plugin.Collection{
//		&pipes_collection.AuditLogCollection{},
//	}
//
//	return map[string]any{
//		pipes_collection.AuditLogCollection{}.Identifier(): pipes_collection.AuditLogRow{},
//	}
//}

func (t *Plugin) Collect(req *proto.CollectRequest) error {
	go t.doCollect(context.Background(), req)

	return nil
}

func (t *Plugin) doCollect(ctx context.Context, req *proto.CollectRequest) {
	// todo config parsing, identify collection type etc.

	// TODO parse config and use to build collection
	//  tactical - create collection
	config := &pipes_types.AuditLogCollectionConfig{Token: os.Getenv("PIPES_TOKEN")}

	var col = pipes_collection.NewAuditLogCollection()

	// TEMP call init
	col.Init(config)

	// add ourselves as an observer
	col.AddObserver(t)

	// signal we have started
	t.OnStarted(req)

	// tell the collection to start collecting - this is a blocking call
	err := col.Collect(ctx, req)

	// signal we have completed
	t.OnComplete(req, err)
}
