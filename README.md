# grafana-teams-sync

Synchronize users from Keycloak Roles to Grafana Teams.

## Description 

__grafana-teams-sync__ - tracks Keycloak roles based on a specified regular expression. The service collects users in roles and concurrently monitors the state of Grafana. 

From Grafana, __grafana-teams-sync__ gathers information about Users, Teams, Permissions, and Folders. Additionally, service synchronize Keycloak state to Grafana.

How synchronization works:
1) For each Keycloak role satisfying the regex,  __grafana-teams-sync__ creates Grafana Folder.
2) For each such role, __grafana-teams-sync__ creates Grafana Team.
3) Permission is granted to the Grafana Folder for the Team.
4) Existing (*) Grafana Users are added to the Grafana Team.

(*) Due to API limitations, for an OIDC user to receive their permissions in Grafana, they need to log in and wait for the synchronization procedure.

Thanks to @rashaev for inspiration. 

## Docker images

Docker images are published on Dockerhub: [ksxack/grafana-teams-sync](https://hub.docker.com/r/ksxack/grafana-teams-sync)

## Configuration 

Environment variables:

| Name                          | Type          | Description                                                            |
|-------------------------------|---------------|------------------------------------------------------------------------|
| GRAFANA_URL                   | url.URL       | URL of the Grafana instance                                            |
| KEYCLOAK_URL                  | url.URL       | URL of the Keycloak instance                                           |
| LOG_LEVEL                     | string        | Logging level (e.g., info, debug)                                      |
| ROLES_REGEX_RO                | string        | ReadOnly Keycloak roles regex (e.g. "-ro")                             |
| ROLES_REGEX_RW                | string        | ReadWrite Keycloak roles regex (e.g. "-rw")                            |
| KEYCLOAK_MONITOR_INTERVAL     | time.Duration | How often should the Keycloak state in memory be updated, default:"5m" |
| GRAFANA_MONITOR_INTERVAL      | time.Duration | How often should the Grafana state in memory be updated, default:"5m"  |
| SYNC_INTERVAL                 | time.Duration | How often should sync process be launched, default:"5m"                |
| GRAFANA_USER                  | string        | Admin user (not OIDC)                                                  |
| GRAFANA_PASSWORD              | string        | Admin password                                                         |
| KEYCLOAK_REALM                | string        | Keycloak Realm with Grafana client                                     |
| KEYCLOAK_CLIENT_NAME          | string        | Grafana client name in Keycloak                                        |
| KEYCLOAK_CLIENT_SECRET        | string        | Grafana client secret in Keycloak                                      |
| KEYCLOAK_MASTER_CLIENT_NAME   | string        | Stub client name in Keycloak Master Realm (to obtain token)            |
| KEYCLOAK_MASTER_CLIENT_SECRET | string        | Stub client secret                                                     |
| KEYCLOAK_USER                 | string        | Keycloak admin user                                                    |
| KEYCLOAK_PASSWORD             | string        | Keycloak admin password                                                |

## ToDo

1. Now grafana-teams-sync is able only to add Users permissions and don't able to delete
2. Algorithm of synchronization process is now very weak and could be improved