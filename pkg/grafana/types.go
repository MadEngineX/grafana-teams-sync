package grafana

import (
	"context"
	"github.com/MadEngineX/grafana-teams-sync/pkg/config"
	gapi "github.com/grafana/grafana-api-golang-client"
	"sync"
)

type Grafana struct {
	Client *gapi.Client
	Cache  *sync.Map
	Ctx    context.Context
	Config config.Config
}

type Cache struct {
	Users   []gapi.UserSearch
	Folders []CacheFolder
	Teams   []CacheTeam
}

type CacheFolder struct {
	Folder      gapi.Folder
	Permissions []*gapi.FolderPermission
}

type CacheTeam struct {
	Team    gapi.Team
	Members []*gapi.TeamMember
}
