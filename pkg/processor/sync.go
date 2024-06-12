package processor

import (
	"context"
	"fmt"
	"github.com/MadEngineX/grafana-teams-sync/pkg/config"
	"github.com/MadEngineX/grafana-teams-sync/pkg/grafana"
	"github.com/MadEngineX/grafana-teams-sync/pkg/keycloak"
	"github.com/MadEngineX/grafana-teams-sync/pkg/utils"
	gapi "github.com/grafana/grafana-api-golang-client"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func (s *SyncProcessor) SyncState(shutdownCh chan struct{}, wg *sync.WaitGroup, rolesRegex string) {
	defer wg.Done()

	loggerID := utils.GenerateLoggerID()

	for {
		select {
		case <-shutdownCh:
			log.WithField("loggerID", loggerID).Infof("SyncProcessor: received shutdown signal")
			return
		default:

			keycloakCache, err := s.loadCache("keycloak cache")
			if err != nil {
				log.WithField("loggerID", loggerID).Errorf("Unable to load kecylaok state from cache, will try in next iteration: %s", err.Error())
				time.Sleep(s.Config.SyncInterval)
				continue
			}

			grafanaCache, err := s.loadCache("grafana cache")
			if err != nil {
				log.WithField("loggerID", loggerID).Errorf("Unable to load grafana state from cache, will try in next iteration: %s", err.Error())
				time.Sleep(s.Config.SyncInterval)
				continue
			}

			log.WithField("loggerID", loggerID).Debugf("SyncProcessorState start sync states for roles: %s", rolesRegex)

			err = s.syncRolesToFolders(keycloakCache, grafanaCache, loggerID, rolesRegex)
			if err != nil {
				if strings.Contains(err.Error(), "already exists") {
					log.WithField("loggerID", loggerID).Warnf("Ubable to syncRolesToFolders: %s", err.Error())
				} else {
					log.WithField("loggerID", loggerID).Errorf("Ubable to syncRolesToFolders: %s", err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
			}

			err = s.syncRolesToTeams(keycloakCache, grafanaCache, loggerID, rolesRegex)
			if err != nil {
				if strings.Contains(err.Error(), "name taken") {
					log.WithField("loggerID", loggerID).Warnf("Ubable to syncRolesToTeams: %s", err.Error())
				} else {
					log.WithField("loggerID", loggerID).Errorf("Ubable to syncRolesToTeams: %s", err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
			}

			err = s.syncTeamsPermissions(keycloakCache, grafanaCache, loggerID, rolesRegex)
			if err != nil {
				if strings.Contains(err.Error(), "already exists") {
					log.WithField("loggerID", loggerID).Warnf("Ubable to syncTeamsPermissions: %s", err.Error())
				} else {
					log.WithField("loggerID", loggerID).Errorf("Ubable to syncTeamsPermissions: %s", err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
			}

			err = s.syncUsersToTeams(keycloakCache, grafanaCache, loggerID, rolesRegex)
			if err != nil {
				if strings.Contains(err.Error(), "already added") {
					log.WithField("loggerID", loggerID).Warnf("Ubable to syncUsersToTeams: %s", err.Error())
				} else {
					log.WithField("loggerID", loggerID).Errorf("Ubable to syncUsersToTeams: %s", err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
			}

			time.Sleep(10 * time.Second)
		}
	}
}

func (s *SyncProcessor) loadCache(name string) (CommonCache, error) {
	cache := CommonCache{}
	cacheInterface, ok := s.Cache.Load(name)
	if !ok {
		return cache, fmt.Errorf("SyncProcessor: unable to gr.Cache.Load")
	}

	cache, ok = cacheInterface.(CommonCache)
	if !ok {
		return cache, fmt.Errorf("SyncProcessor: unable to assert cache")
	}
	return cache, nil
}

// sync KeycloakRoles to GrafanaFolders
func (s *SyncProcessor) syncRolesToFolders(keycloakCache, grafanaCache CommonCache, loggerID, roleRegex string) error {
	var keycloakRoles []keycloak.CacheRole

	if roleRegex == s.Config.RolesRegexRO {
		keycloakRoles = keycloakCache.KeycloakState.RolesRO
	} else {
		keycloakRoles = keycloakCache.KeycloakState.RolesRW
	}
	for _, role := range keycloakRoles {

		log.WithField("loggerID", loggerID).Debugf("Syncing role: %s to Grafana Folder", role.Role.Name)

		folderName := strings.Replace(role.Role.Name, roleRegex, "", 1)

		if folderExists(folderName, grafanaCache.GrafanaState.Folders) {
			log.WithField("loggerID", loggerID).Debugf("Folder %s already exists in Grafana", folderName)
			continue
		}

		grafanaFolder, err := s.GrafanaClient.Client.NewFolder(folderName)
		if err != nil {
			return fmt.Errorf("unable to create folder s.GrafanaClient.Client.NewFolder: %s", err.Error())
		}
		log.WithField("loggerID", loggerID).Debugf("Created new Folder in Grafana: ID %d, Name %s", grafanaFolder.ID, grafanaFolder.Title)

	}

	return nil

}

// sync KeycloakRoles to GrafanaTeams
func (s *SyncProcessor) syncRolesToTeams(keycloakCache, grafanaCache CommonCache, loggerID, roleRegex string) error {
	var keycloakRoles []keycloak.CacheRole

	if roleRegex == s.Config.RolesRegexRO {
		keycloakRoles = keycloakCache.KeycloakState.RolesRO
	} else {
		keycloakRoles = keycloakCache.KeycloakState.RolesRW
	}
	for _, role := range keycloakRoles {

		log.WithField("loggerID", loggerID).Debugf("Syncing role: %s to Grafana Team", role.Role.Name)

		teamName := role.Role.Name

		if teamExists(teamName, grafanaCache.GrafanaState.Teams) {
			log.WithField("loggerID", loggerID).Debugf("Team %s already exists in Grafana", teamName)
			continue
		}

		grafanaTeam, err := s.GrafanaClient.Client.AddTeam(teamName, "")
		if err != nil {
			return fmt.Errorf("unable to create team s.GrafanaClient.Client.AddTeam: %s", err.Error())
		}
		log.WithField("loggerID", loggerID).Debugf("Created new Team in Grafana: ID %d, Name %s", grafanaTeam, teamName)

	}

	return nil

}

// sync GrafanaPermissions to GrafanaTeams
func (s *SyncProcessor) syncTeamsPermissions(keycloakCache, grafanaCache CommonCache, loggerID, roleRegex string) error {
	var keycloakRoles []keycloak.CacheRole

	if roleRegex == s.Config.RolesRegexRO {
		keycloakRoles = keycloakCache.KeycloakState.RolesRO
	} else {
		keycloakRoles = keycloakCache.KeycloakState.RolesRW
	}
OUTER:
	for _, role := range keycloakRoles {

		log.WithField("loggerID", loggerID).Debugf("Syncing role: %s to Grafana TeamsPermissions", role.Role.Name)

		folderName := strings.Replace(role.Role.Name, roleRegex, "", 1)
		teamName := role.Role.Name

		if permissionExists(teamName, folderName, grafanaCache.GrafanaState.Folders) {
			log.WithField("loggerID", loggerID).Debugf("Team %s already exists in Grafana", teamName)
			continue
		}

		permissions := gapi.PermissionItems{}
		item := gapi.PermissionItem{}

		teamID := getTeamID(teamName, grafanaCache.GrafanaState.Teams)
		folderUID := getFolderUID(folderName, grafanaCache.GrafanaState.Folders)

		if roleRegex == s.Config.RolesRegexRO {
			item = gapi.PermissionItem{
				TeamID:     teamID,
				Permission: 1,
			}
		} else {

			item = gapi.PermissionItem{
				TeamID:     teamID,
				Permission: 2,
			}
		}

		permissions.Items = append(permissions.Items, &item)

		currentPermissions, err := s.GrafanaClient.Client.FolderPermissions(folderUID)
		if err != nil {
			return fmt.Errorf("unable to s.GrafanaClient.Client.FolderPermissions: %s", err.Error())
		}

		for _, currentPermission := range currentPermissions {
			if currentPermission.TeamID == item.TeamID && currentPermission.Permission == item.Permission {
				continue OUTER

			} else {
				permissions.Items = append(permissions.Items, &gapi.PermissionItem{
					TeamID:     currentPermission.TeamID,
					Permission: currentPermission.Permission,
				})
			}
		}

		err = s.GrafanaClient.Client.UpdateFolderPermissions(folderUID, &permissions)
		if err != nil {
			return fmt.Errorf("unable to s.GrafanaClient.Client.UpdateFolderPermissions: %s", err.Error())
		}
		log.WithField("loggerID", loggerID).Debugf("Created Permissions for Folder in Grafana: Folder %s, Team %s", folderName, teamName)

	}

	return nil
}

// sync GrafanaUsers to GrafanaTeams according to KeycloakRoles
func (s *SyncProcessor) syncUsersToTeams(keycloakCache, grafanaCache CommonCache, loggerID, roleRegex string) error {
	var keycloakRoles []keycloak.CacheRole

	if roleRegex == s.Config.RolesRegexRO {
		keycloakRoles = keycloakCache.KeycloakState.RolesRO
	} else {
		keycloakRoles = keycloakCache.KeycloakState.RolesRW
	}

	for _, grafanaUser := range grafanaCache.GrafanaState.Users {
		for _, keycloakRole := range keycloakRoles {

			teamName := keycloakRole.Role.Name

			for _, keycloakUser := range keycloakRole.Users {
				if keycloakUser.Username == grafanaUser.Login {

					teamID := getTeamID(teamName, grafanaCache.GrafanaState.Teams)

					if userAlreadyInTeam(keycloakUser.Username, teamName, grafanaCache.GrafanaState.Teams) {
						log.WithField("loggerID", loggerID).Debugf("User %s already exists in Grafana Team %s", grafanaUser.Login, teamName)
						continue
					}

					log.WithField("loggerID", loggerID).Debugf("Adding Member to Grafana Team: Team %s, User %s,", teamName, grafanaUser.Login)
					err := s.GrafanaClient.Client.AddTeamMember(teamID, grafanaUser.ID)
					if err != nil {
						return fmt.Errorf("unable to s.GrafanaClient.Client.AddTeamMember: %s", err.Error())
					}

				}
			}
		}
	}

	log.WithField("loggerID", loggerID).Debugf("Synced GrafanaUsers to GrafanaTeams")

	return nil
}

func New(ctx context.Context, cfg config.Config, commonCache *sync.Map, keycloakClient *keycloak.Keycloak, grafanaClient *grafana.Grafana) *SyncProcessor {

	return &SyncProcessor{
		KeycloakClient: keycloakClient,
		GrafanaClient:  grafanaClient,
		Cache:          commonCache,
		Ctx:            ctx,
		Config:         cfg,
	}
}
