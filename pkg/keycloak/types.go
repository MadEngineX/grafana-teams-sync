package keycloak

import (
	"github.com/MadEngineX/grafana-teams-sync/pkg/config"
	"net/http"
	"net/url"
	"sync"
)

type Cache struct {
	RolesRO []CacheRole
	RolesRW []CacheRole
}

type CacheRole struct {
	Role  Role
	Users []User
}

type Group struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Path      string  `json:"path"`
	SubGroups []Group `json:"subgroups"`
}

type User struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Lastname         string `json:"lastName"`
	Firstname        string `json:"firstName"`
	Enabled          bool   `json:"enabled"`
	EmailVerified    bool   `json:"emailVerified"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
}

type Keycloak struct {
	URL          url.URL
	Username     string
	Password     string
	ClientName   string
	ClientSecret string
	Client       *http.Client
	Realm        string
	Cache        *sync.Map
	Cfg          config.Config
}

type Clients struct {
	Clients []Client
}

type Client struct {
	ID                                 string                             `json:"id"`
	ClientID                           string                             `json:"clientId"`
	Name                               string                             `json:"name"`
	Description                        string                             `json:"description"`
	RootURL                            string                             `json:"rootUrl"`
	AdminURL                           string                             `json:"adminUrl"`
	BaseURL                            string                             `json:"baseUrl"`
	SurrogateAuthRequired              bool                               `json:"surrogateAuthRequired"`
	Enabled                            bool                               `json:"enabled"`
	AlwaysDisplayInConsole             bool                               `json:"alwaysDisplayInConsole"`
	ClientAuthenticatorType            string                             `json:"clientAuthenticatorType"`
	Secret                             string                             `json:"secret"`
	RedirectUris                       []string                           `json:"redirectUris"`
	WebOrigins                         []string                           `json:"webOrigins"`
	NotBefore                          int                                `json:"notBefore"`
	BearerOnly                         bool                               `json:"bearerOnly"`
	ConsentRequired                    bool                               `json:"consentRequired"`
	StandardFlowEnabled                bool                               `json:"standardFlowEnabled"`
	ImplicitFlowEnabled                bool                               `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled          bool                               `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled             bool                               `json:"serviceAccountsEnabled"`
	AuthorizationServicesEnabled       bool                               `json:"authorizationServicesEnabled"`
	PublicClient                       bool                               `json:"publicClient"`
	FrontchannelLogout                 bool                               `json:"frontchannelLogout"`
	Protocol                           string                             `json:"protocol"`
	Attributes                         Attributes                         `json:"attributes"`
	AuthenticationFlowBindingOverrides AuthenticationFlowBindingOverrides `json:"authenticationFlowBindingOverrides"`
	FullScopeAllowed                   bool                               `json:"fullScopeAllowed"`
	NodeReRegistrationTimeout          int                                `json:"nodeReRegistrationTimeout"`
	ProtocolMappers                    []ProtocolMappers                  `json:"protocolMappers"`
	DefaultClientScopes                []string                           `json:"defaultClientScopes"`
	OptionalClientScopes               []string                           `json:"optionalClientScopes"`
	Access                             Access                             `json:"access"`
}
type Attributes struct {
	OidcCibaGrantEnabled                  string `json:"oidc.ciba.grant.enabled"`
	Oauth2DeviceAuthorizationGrantEnabled string `json:"oauth2.device.authorization.grant.enabled"`
	ClientSecretCreationTime              string `json:"client.secret.creation.time"`
	BackchannelLogoutSessionRequired      string `json:"backchannel.logout.session.required"`
	BackchannelLogoutRevokeOfflineTokens  string `json:"backchannel.logout.revoke.offline.tokens"`
}
type AuthenticationFlowBindingOverrides struct {
}
type Config struct {
	UserSessionNote  string `json:"user.session.note"`
	IDTokenClaim     string `json:"id.token.claim"`
	AccessTokenClaim string `json:"access.token.claim"`
	ClaimName        string `json:"claim.name"`
	JSONTypeLabel    string `json:"jsonType.label"`
}
type ProtocolMappers struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Protocol        string `json:"protocol"`
	ProtocolMapper  string `json:"protocolMapper"`
	ConsentRequired bool   `json:"consentRequired"`
	Config          Config `json:"config"`
}
type Access struct {
	View      bool `json:"view"`
	Configure bool `json:"configure"`
	Manage    bool `json:"manage"`
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Composite   bool   `json:"composite"`
	ClientRole  bool   `json:"clientRole"`
	ContainerID string `json:"containerId"`
	Description string `json:"description,omitempty"`
}
