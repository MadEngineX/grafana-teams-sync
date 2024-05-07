package monitors

import (
	"github.com/MadEngineX/grafana-teams-sync/pkg/grafana"
	"github.com/MadEngineX/grafana-teams-sync/pkg/processor"
	"github.com/MadEngineX/grafana-teams-sync/pkg/utils"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func MonitorGrafanaUsers(shutdownCh chan struct{}, wg *sync.WaitGroup, gr *grafana.Grafana) {
	defer wg.Done()

	loggerID := utils.GenerateLoggerID()

	for {
		select {
		case <-shutdownCh:
			log.WithField("loggerID", loggerID).Infof("MonitorKeycloakRoles: received shutdown signal")
			return
		default:

			cache := processor.NewCommonCache()

			grafanaUsers, err := gr.Client.Users()
			if err != nil {
				log.WithField("loggerID", loggerID).Error(err)
			}

			cache.GrafanaState.Users = grafanaUsers

			grafanaTeams, err := gr.Client.SearchTeam("")
			if err != nil {
				log.WithField("loggerID", loggerID).Error(err)
			}

			for _, team := range grafanaTeams.Teams {
				teamMembers, err := gr.Client.TeamMembers(team.ID)
				if err != nil {
					log.Error(err)
				}

				cacheTeam := grafana.CacheTeam{
					Team:    *team,
					Members: teamMembers,
				}

				cache.GrafanaState.Teams = append(cache.GrafanaState.Teams, cacheTeam)
			}

			grafanaFolders, err := gr.Client.Folders()
			if err != nil {
				log.WithField("loggerID", loggerID).Error(err)
			}

			for _, folder := range grafanaFolders {
				grafanaFolderPermissions, err := gr.Client.FolderPermissions(folder.UID)
				if err != nil {
					log.WithField("loggerID", loggerID).Error(err)
				}

				cacheFolder := grafana.CacheFolder{
					Folder:      folder,
					Permissions: grafanaFolderPermissions,
				}

				cache.GrafanaState.Folders = append(cache.GrafanaState.Folders, cacheFolder)
			}

			gr.Cache.Store("grafana cache", *cache)
			log.WithField("loggerID", loggerID).Debugf("MonitorGrafanaUsers: state synced to cache, continue...")

			time.Sleep(gr.Config.GrafanaMonitorInterval)
		}
	}
}
