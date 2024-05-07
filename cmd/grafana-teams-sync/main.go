package main

import (
	"context"
	"github.com/MadEngineX/grafana-teams-sync/pkg/config"
	"github.com/MadEngineX/grafana-teams-sync/pkg/grafana"
	"github.com/MadEngineX/grafana-teams-sync/pkg/keycloak"
	"github.com/MadEngineX/grafana-teams-sync/pkg/monitors"
	"github.com/MadEngineX/grafana-teams-sync/pkg/processor"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.Info("Starting...")

	shutdownCh := make(chan struct{})
	wg := sync.WaitGroup{}
	cfg := config.AppConfig

	serverCtx := context.Background()
	rolesRegexps := []string{cfg.RolesRegexRO, cfg.RolesRegexRW}

	commonCacheStruct := processor.NewCommonCache()
	commonCache := &sync.Map{}
	commonCache.Store("cache", *commonCacheStruct)

	// Create channel for graceful shutdown signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Init Keycloak client
	keycloakClient, err := keycloak.New(serverCtx, *cfg, commonCache)
	if err != nil {
		log.Fatalf("Unable to initialize Keycloak client: %s", err.Error())
	}

	// Init Grafana client
	grafanaClient, err := grafana.New(serverCtx, *cfg, commonCache)
	if err != nil {
		log.Fatalf("Unable to initialize Grafana client: %s", err.Error())
	}

	syncProcessor := processor.New(serverCtx, *cfg, commonCache, keycloakClient, grafanaClient)

	wg.Add(1)

	go monitors.MonitorKeycloakRoles(shutdownCh, &wg, keycloakClient, rolesRegexps)
	go monitors.MonitorGrafanaUsers(shutdownCh, &wg, grafanaClient)

	go syncProcessor.SyncState(shutdownCh, &wg, cfg.RolesRegexRO)
	go syncProcessor.SyncState(shutdownCh, &wg, cfg.RolesRegexRW)

	// Ожидаем сигнала для завершения
	select {
	case sig := <-stopChan:
		log.Warnf("Received signal: %s\n", sig)
		close(shutdownCh) // Отправляем сигнал завершения в канал
		wg.Wait()         // Дожидаемся завершения всех горутин
		log.Infof("Graceful shutdown completed")
	}
}
