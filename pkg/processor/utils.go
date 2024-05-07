package processor

import (
	"github.com/MadEngineX/grafana-teams-sync/pkg/grafana"
)

func folderExists(folder string, grafanaFolders []grafana.CacheFolder) bool {
	for _, grafanaFolder := range grafanaFolders {
		if folder == grafanaFolder.Folder.Title {
			return true
		}
	}
	return false
}

func teamExists(team string, grafanaTeams []grafana.CacheTeam) bool {
	for _, grafanaTeam := range grafanaTeams {
		if team == grafanaTeam.Team.Name {
			return true
		}
	}
	return false
}

func permissionExists(team, folder string, grafanaFolders []grafana.CacheFolder) bool {
	for _, grafanaFolder := range grafanaFolders {
		if folder == grafanaFolder.Folder.Title {
			for _, permission := range grafanaFolder.Permissions {
				if team == permission.PermissionName {
					return true
				}
			}
		}
	}
	return false
}

func userAlreadyInTeam(username, teamName string, grafanaTeams []grafana.CacheTeam) bool {
	for _, grafanaTeam := range grafanaTeams {
		if teamName == grafanaTeam.Team.Name {
			for _, teamMember := range grafanaTeam.Members {
				if teamMember.Login == username {
					return true
				}
			}
		}
	}
	return false
}

func getTeamID(team string, grafanaTeams []grafana.CacheTeam) int64 {
	for _, grafanaTeam := range grafanaTeams {
		if team == grafanaTeam.Team.Name {
			return grafanaTeam.Team.ID
		}
	}
	return -1
}

func getFolderUID(folderName string, grafanaFolders []grafana.CacheFolder) string {
	for _, folder := range grafanaFolders {
		if folderName == folder.Folder.Title {
			return folder.Folder.UID
		}
	}
	return ""
}
