package plugin

import (
	"context"
	"errors"
	"sync"

	"github.com/turbot/tailpipe-plugin-pipes/collection"
	"github.com/turbot/tailpipe-plugin-pipes/source"

	sdkcollection "github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	sdksource "github.com/turbot/tailpipe-plugin-sdk/source"
)

type PipesPlugin struct {
	ctx context.Context

	// observers is a list of observers that will be notified of events.
	observers      []plugin.PluginObserver
	observersMutex sync.RWMutex
}

func (p *PipesPlugin) Identifier() string {
	return "pipes"
}

func (p *PipesPlugin) Init(ctx context.Context) error {
	p.ctx = ctx
	return nil
}

func (p *PipesPlugin) Context() context.Context {
	return p.ctx
}

func (p *PipesPlugin) Validate() error {
	return errors.ErrUnsupported
}

func (p *PipesPlugin) AddObserver(observer plugin.PluginObserver) {
	p.observersMutex.Lock()
	defer p.observersMutex.Unlock()
	p.observers = append(p.observers, observer)
}

func (p *PipesPlugin) RemoveObserver(observer plugin.PluginObserver) {
	p.observersMutex.Lock()
	defer p.observersMutex.Unlock()
	for i, o := range p.observers {
		if o == observer {
			p.observers = append(p.observers[:i], p.observers[i+1:]...)
			break
		}
	}
}

func (p *PipesPlugin) Sources() map[string]sdksource.Plugin {
	return map[string]sdksource.Plugin{
		"pipes_audit_log": &source.PipesAuditLogSource{},
	}
}

func (p *PipesPlugin) Collections() map[string]sdkcollection.Plugin {
	return map[string]sdkcollection.Plugin{
		"pipes_audit_log": &collection.PipesAuditLogCollection{},
	}
}
