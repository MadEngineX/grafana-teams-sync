package config

import (
	"github.com/caarlos0/env/v9"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

var AppConfig = loadConfig()

type Config struct {
	GrafanaURL                 url.URL       `env:"GRAFANA_URL,notEmpty"`
	GrafanaUser                string        `env:"GRAFANA_USER,notEmpty"`
	GrafanaPassword            string        `env:"GRAFANA_PASSWORD,notEmpty"`
	KeycloakURL                url.URL       `env:"KEYCLOAK_URL,notEmpty"`
	KeycloakUser               string        `env:"KEYCLOAK_USER,notEmpty"`
	KeycloakPassword           string        `env:"KEYCLOAK_PASSWORD,notEmpty"`
	KeycloakMasterClientName   string        `env:"KEYCLOAK_MASTER_CLIENT_NAME,notEmpty"`
	KeycloakMasterClientSecret string        `env:"KEYCLOAK_MASTER_CLIENT_SECRET,notEmpty"`
	KeycloakClientName         string        `env:"KEYCLOAK_CLIENT_NAME,notEmpty"`
	KeycloakClientSecret       string        `env:"KEYCLOAK_CLIENT_SECRET,notEmpty"`
	KeycloakRealm              string        `env:"KEYCLOAK_REALM,notEmpty" envDefault:"master"`
	RolesRegexRO               string        `env:"ROLES_REGEX_RO,notEmpty"`
	RolesRegexRW               string        `env:"ROLES_REGEX_RW,notEmpty"`
	KeycloakMonitorInterval    time.Duration `env:"KEYCLOAK_MONITOR_INTERVAL,notEmpty" envDefault:"5m"`
	GrafanaMonitorInterval     time.Duration `env:"GRAFANA_MONITOR_INTERVAL,notEmpty" envDefault:"5m"`
	SyncInterval               time.Duration `env:"SYNC_INTERVAL,notEmpty" envDefault:"5m"`

	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
}

func loadConfig() *Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Unable to parse envs: %s", err.Error())
	}

	cfg.initLogger()

	return &cfg
}

func (c *Config) initLogger() {
	const defaultLogLevel = "info"

	logLevelValue, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Warnf(
			"Invalid 'LOG_LEVEL': [%v], will use default LOG_LEVEL: [%v]. "+
				"Allowed values: trace, debug, info, warn, warning, error, fatal, panic", c.LogLevel, defaultLogLevel,
		)

		logLevelValue = log.InfoLevel
	}
	log.SetLevel(logLevelValue)
}
