package pipes

import (
	"context"
	"fmt"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_collection"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_types"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"log/slog"
	"os"
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

func (t *Plugin) Collect(ctx context.Context, req *proto.CollectRequest) error {
	go func() {
		if err := t.doCollect(context.Background(), req); err != nil {
			// TODO #err handle error
			slog.Error("doCollect failed", "error", err)
		}
	}()

	return nil
}

func (t *Plugin) doCollect(ctx context.Context, req *proto.CollectRequest) error {
	// todo config parsing, identify collection type etc.

	// TODO parse config and use to build collection
	//  tactical - create collection
	config := &pipes_types.AuditLogCollectionConfig{Token: os.Getenv("PIPES_TOKEN")}

	var col = pipes_collection.NewAuditLogCollection()

	// TEMP call init
	err := col.Init(config)
	if err != nil {
		return fmt.Errorf("error initializing collection: %w", err)
	}

	// add ourselves as an observer
	if err := col.AddObserver(t); err != nil {
		return fmt.Errorf("error adding observer: %w", err)
	}

	// signal we have started
	if err := t.OnStarted(ctx, req); err != nil {
		return fmt.Errorf("error signalling started: %w", err)
	}

	// tell the collection to start collecting - this is a blocking call
	if err := col.Collect(ctx, req); err != nil {
		return fmt.Errorf("error collecting: %w", err)
	}

	// signal we have completed
	return t.OnCompleted(ctx, req, err)
}
