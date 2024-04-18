package source

import (
	"context"
	"errors"
	"os"
	"sync"

	sdkconfig "github.com/turbot/tailpipe-plugin-sdk/config"
	"github.com/turbot/tailpipe-plugin-sdk/source"
)

type PipesAuditLogSourceConfig struct {
	Token string `json:"token"`
}

type PipesAuditLogSource struct {
	Config PipesAuditLogSourceConfig

	ctx            context.Context
	observers      []source.SourceObserver
	observersMutex sync.RWMutex
}

func (s *PipesAuditLogSource) Identifier() string {
	return "pipes_audit_log"
}

func (s *PipesAuditLogSource) Init(ctx context.Context) error {
	s.ctx = ctx
	return s.Validate()
}

func (s *PipesAuditLogSource) Context() context.Context {
	return s.ctx
}

func (s *PipesAuditLogSource) Validate() error {
	return nil
}

func (s *PipesAuditLogSource) AddObserver(observer source.SourceObserver) {
	s.observersMutex.Lock()
	defer s.observersMutex.Unlock()
	s.observers = append(s.observers, observer)
}

func (s *PipesAuditLogSource) RemoveObserver(observer source.SourceObserver) {
	s.observersMutex.Lock()
	defer s.observersMutex.Unlock()
	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *PipesAuditLogSource) LoadConfig(configRaw []byte) error {
	if err := sdkconfig.Load(configRaw, &s.Config); err != nil {
		return err
	}
	if s.Config.Token == "" {
		s.Config.Token = os.Getenv("PIPES_TOKEN")
	}
	return nil
}

func (s *PipesAuditLogSource) ValidateConfig() error {
	if s.Config.Token == "" {
		return errors.New("token is required")
	}
	return nil
}

func (s *PipesAuditLogSource) DiscoverArtifacts(ctx context.Context) error {
	for _, observer := range s.observers {
		observer.NotifyArtifactDiscovered(&source.ArtifactInfo{Name: "api"})
	}
	return nil
}

func (s *PipesAuditLogSource) DownloadArtifact(ctx context.Context, ai *source.ArtifactInfo) error {
	for _, observer := range s.observers {
		observer.NotifyArtifactDownloaded(&source.Artifact{ArtifactInfo: *ai})
	}
	return nil
}
