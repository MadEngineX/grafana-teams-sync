package processor

import (
	"context"
	"github.com/MadEngineX/grafana-teams-sync/pkg/config"
	"github.com/MadEngineX/grafana-teams-sync/pkg/grafana"
	"github.com/MadEngineX/grafana-teams-sync/pkg/keycloak"
	"sync"
)

type CommonCache struct {
	KeycloakState keycloak.Cache
	GrafanaState  grafana.Cache
}

func NewCommonCache() *CommonCache {
	keycloakCache := keycloak.Cache{}
	grafanaCache := grafana.Cache{}
	return &CommonCache{
		KeycloakState: keycloakCache,
		GrafanaState:  grafanaCache,
	}
}

type SyncProcessor struct {
	GrafanaClient  *grafana.Grafana
	KeycloakClient *keycloak.Keycloak
	Cache          *sync.Map
	Ctx            context.Context
	Config         config.Config
	LoggerID       string
}
