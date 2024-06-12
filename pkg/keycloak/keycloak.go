package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MadEngineX/grafana-teams-sync/pkg/config"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"sync"
)

func (k *Keycloak) GetRoles(rolesRegex string) ([]Role, error) {

	clientID, err := k.GetClientID()
	if err != nil {
		return nil, fmt.Errorf("unable to GetClientID: %s", err.Error())
	}

	req, err := http.NewRequest("GET", k.URL.String()+"/admin/realms/"+k.Realm+"/clients/"+clientID+"/roles?max=1000&search="+rolesRegex, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to http.NewRequest: %s", err.Error())
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to k.Client.Do: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to io.ReadAll: %s", err.Error())
	}

	var roles []Role

	if len(body) == 0 {
		return nil, fmt.Errorf("keycloak returns 0 roles, check roles regexp")
	}

	err = json.Unmarshal(body, &roles)
	if err != nil {
		return nil, fmt.Errorf("unable to json.Unmarshal: %s", err.Error())
	}

	return roles, nil

}

func (k *Keycloak) GetUsersInRole(roleName string) ([]User, error) {

	clientID, err := k.GetClientID()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", k.URL.String()+"/admin/realms/"+k.Realm+"/clients/"+clientID+"/roles/"+roleName+"/users", nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var users []User

	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (k *Keycloak) GetGroupsInRole(roleName string) ([]Group, error) {

	clientID, err := k.GetClientID()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", k.URL.String()+"/admin/realms/"+k.Realm+"/clients/"+clientID+"/roles/"+roleName+"/groups", nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var groups []Group

	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil

}

func (k *Keycloak) GetUsersInGroup(groupID string) ([]User, error) {

	req, err := http.NewRequest("GET", k.URL.String()+"/admin/realms/"+k.Realm+"/groups/"+groupID+"/members", nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var users []User

	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (k *Keycloak) GetClientID() (string, error) {
	req, err := http.NewRequest("GET", k.URL.String()+"/admin/realms/"+k.Realm+"/clients?clientId="+k.ClientName, nil)
	if err != nil {
		return "", err
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var clients []Client

	err = json.Unmarshal(body, &clients)
	if err != nil {
		return "", fmt.Errorf("GetClientID error while json.Unmarshal: %s", err.Error())
	}

	if len(clients) == 0 {
		return "", fmt.Errorf("client not found: %s", k.ClientName)
	}

	return clients[0].ID, nil
}

func New(ctx context.Context, cfg config.Config, commonCache *sync.Map) (*Keycloak, error) {

	conf := &oauth2.Config{
		ClientID:     cfg.KeycloakMasterClientName,
		ClientSecret: cfg.KeycloakMasterClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.KeycloakURL.String() + "/realms/master/protocol/openid-connect/auth",
			TokenURL: cfg.KeycloakURL.String() + "/realms/master/protocol/openid-connect/token",
		},
	}

	token, err := conf.PasswordCredentialsToken(ctx, cfg.KeycloakUser, cfg.KeycloakPassword)
	if err != nil {
		return nil, err
	}

	client := conf.Client(ctx, token)

	return &Keycloak{
		URL:          cfg.KeycloakURL,
		Username:     cfg.KeycloakUser,
		Password:     cfg.KeycloakPassword,
		ClientName:   cfg.KeycloakClientName,
		ClientSecret: cfg.KeycloakClientSecret,
		Client:       client,
		Realm:        cfg.KeycloakRealm,
		Cache:        commonCache,
		Cfg:          cfg,
	}, nil
}
