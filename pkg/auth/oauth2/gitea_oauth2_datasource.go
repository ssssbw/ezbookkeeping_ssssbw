package oauth2

import (
	"encoding/json"
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

type giteaUserInfoResponse struct {
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

// GiteaOAuth2DataSource represents Gitea OAuth 2.0 data source
type GiteaOAuth2DataSource struct {
	CommonOAuth2DataSource
	baseUrl string
}

// GetAuthUrl returns the authentication url of the Gitea data source
func (s *GiteaOAuth2DataSource) GetAuthUrl() string {
	// Reference: https://docs.gitea.com/development/oauth2-provider
	return s.baseUrl + "login/oauth/authorize"
}

// GetTokenUrl returns the token url of the Gitea data source
func (s *GiteaOAuth2DataSource) GetTokenUrl() string {
	// Reference: https://docs.gitea.com/development/oauth2-provider
	return s.baseUrl + "login/oauth/access_token"
}

// GetUserInfoRequest returns the user info request of the Gitea data source
func (s *GiteaOAuth2DataSource) GetUserInfoRequest() (*http.Request, error) {
	// Reference: https://gitea.com/api/swagger#/user/userGetCurrent
	req, err := http.NewRequest("GET", s.baseUrl+"api/v1/user", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetScopes returns the scopes required by the Gitea provider
func (s *GiteaOAuth2DataSource) GetScopes() []string {
	return []string{"read:user"}
}

// ParseUserInfo returns the user info by parsing the response body
func (s *GiteaOAuth2DataSource) ParseUserInfo(c core.Context, body []byte) (*OAuth2UserInfo, error) {
	userInfoResp := &giteaUserInfoResponse{}
	err := json.Unmarshal(body, &userInfoResp)

	if err != nil {
		log.Warnf(c, "[gitea_oauth2_datasource.ParseUserInfo] failed to parse user profile response body, because %s", err.Error())
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.Login == "" {
		log.Warnf(c, "[gitea_oauth2_datasource.ParseUserInfo] invalid user profile response body")
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	return &OAuth2UserInfo{
		UserName: userInfoResp.Login,
		Email:    userInfoResp.Email,
		NickName: userInfoResp.FullName,
	}, nil
}

// NewGiteaOAuth2Provider creates a new Gitea OAuth 2.0 provider instance
func NewGiteaOAuth2Provider(baseUrl string) OAuth2Provider {
	if baseUrl[len(baseUrl)-1] != '/' {
		baseUrl += "/"
	}

	return &CommonOAuth2Provider{
		dataSource: &GiteaOAuth2DataSource{
			baseUrl: baseUrl,
		},
	}
}
