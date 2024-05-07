package grafana

import (
	"context"
	"github.com/MadEngineX/grafana-teams-sync/pkg/config"
	"net/url"
	"sync"

	gapi "github.com/grafana/grafana-api-golang-client"
)

const (
	Viewer = "Viewer"
	Editor = "Editor"
)

func New(ctx context.Context, cfg config.Config, commonCache *sync.Map) (*Grafana, error) {

	conf := gapi.Config{
		BasicAuth: url.UserPassword(cfg.GrafanaUser, cfg.GrafanaPassword),
	}

	client, err := gapi.New(cfg.GrafanaURL.String(), conf)
	if err != nil {
		return nil, err
	}

	return &Grafana{
		Client: client,
		Cache:  commonCache,
		Ctx:    ctx,
		Config: cfg,
	}, nil
}
