package monitors

import (
	"github.com/MadEngineX/grafana-teams-sync/pkg/keycloak"
	"github.com/MadEngineX/grafana-teams-sync/pkg/processor"
	"github.com/MadEngineX/grafana-teams-sync/pkg/utils"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func MonitorKeycloakRoles(shutdownCh chan struct{}, wg *sync.WaitGroup, kk *keycloak.Keycloak, rolesRegexps []string) {
	defer wg.Done()

	loggerID := utils.GenerateLoggerID()

	for {
		select {
		case <-shutdownCh:
			log.WithField("loggerID", loggerID).Infof("MonitorKeycloakRoles: received shutdown signal")
			return
		default:
			cache := processor.NewCommonCache()

			for _, roleRegexp := range rolesRegexps {

				keycloakRoles, err := kk.GetRoles(roleRegexp)
				if err != nil {
					if strings.Contains(err.Error(), "unable to GetClientID") {
						log.WithField("loggerID", loggerID).Fatalf("Unable to GetRoles: %s", err.Error())
					}
					log.WithField("loggerID", loggerID).Errorf("Unable to GetRoles: %s", err.Error())
				}

				for _, role := range keycloakRoles {

					groupsInRole, err := kk.GetGroupsInRole(role.Name)
					if err != nil {
						log.WithField("loggerID", loggerID).Errorf("Unable to GetGroupsInRole: %s", err.Error())
					}
					for _, group := range groupsInRole {
						usersInRole, err := kk.GetUsersInGroup(group.ID)
						if err != nil {
							log.WithField("loggerID", loggerID).Errorf("Unable to GetUsersInGroup: %s", err.Error())
						}

						cacheRole := keycloak.CacheRole{
							Role:  role,
							Users: usersInRole,
						}

						if roleRegexp == kk.Cfg.RolesRegexRO {
							log.WithField("loggerID", loggerID).Debugf("MonitorKeycloakRoles: adding RO role: %s", role.Name)
							cache.KeycloakState.RolesRO = append(cache.KeycloakState.RolesRO, cacheRole)
						} else {
							log.WithField("loggerID", loggerID).Debugf("MonitorKeycloakRoles: adding keycloak RW role: %s", role.Name)
							cache.KeycloakState.RolesRW = append(cache.KeycloakState.RolesRW, cacheRole)
						}
					}

					usersInRole, err := kk.GetUsersInRole(role.Name)
					if err != nil {
						log.WithField("loggerID", loggerID).Errorf("Unable to GetUsersInRole: %s", err.Error())
					}

					cacheRole := keycloak.CacheRole{
						Role:  role,
						Users: usersInRole,
					}

					if roleRegexp == kk.Cfg.RolesRegexRO {
						log.WithField("loggerID", loggerID).Debugf("MonitorKeycloakRoles: adding RO role: %s", role.Name)
						cache.KeycloakState.RolesRO = append(cache.KeycloakState.RolesRO, cacheRole)
					} else {
						log.WithField("loggerID", loggerID).Debugf("MonitorKeycloakRoles: adding keycloak RW role: %s", role.Name)
						cache.KeycloakState.RolesRW = append(cache.KeycloakState.RolesRW, cacheRole)
					}

				}
			}

			kk.Cache.Store("keycloak cache", *cache)
			log.WithField("loggerID", loggerID).Debugf("MonitorKeycloakRoles: roles synced to cache, continue...")

			time.Sleep(kk.Cfg.KeycloakMonitorInterval)
		}
	}
}
