package facebook

import (
	aah "aahframe.work"
	"aahframe.work/ahttp"
	"aahframe.work/config"
	ess "aahframe.work/essentials"
	"aahframe.work/security/authc"
	"aahframe.work/security/scheme"
	"golang.org/x/oauth2"
	"gopkg.in/resty.v1"
)

var _ authc.PrincipalProvider = (*SubjectPrincipalProvider)(nil)

type FacebookUserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SubjectPrincipalProvider struct implements interface authc.PrincipalProvider.
type SubjectPrincipalProvider struct{}

// Init method initializes the SubjectPrincipalProvider, Init method gets called
// during aah server start up.
func (a *SubjectPrincipalProvider) Init(appCfg *config.Config) error {
	// NOTE: Init is called on application startup
	return nil
}

// Principals method is called to get Subject's Principals information.
func (a *SubjectPrincipalProvider) Principal(keyName string, v ess.Valuer) ([]*authc.Principal, error) {
	token := v.Get(aah.KeyOAuth2Token).(*oauth2.Token)
	fbAuth := aah.App().SecurityManager().AuthScheme(keyName).(*scheme.OAuth2)
	fbClient := resty.NewWithClient(fbAuth.Client(token))

	resp, err := fbClient.R().
		SetHeader(ahttp.HeaderAccept, ahttp.ContentTypeJSON.Mime).
		SetResult(FacebookUserInfo{}).
		Get("https://graph.facebook.com/v3.0/me?fields=name,email")
	if err != nil {
		return nil, err
	}
	userInfo := resp.Result().(*FacebookUserInfo)

	principals := make([]*authc.Principal, 0)
	principals = append(principals, &authc.Principal{Realm: "Facebook", IsPrimary: true, Claim: "Email", Value: userInfo.Email})
	principals = append(principals, &authc.Principal{Realm: "Facebook", Claim: "ID", Value: userInfo.ID})
	principals = append(principals, &authc.Principal{Realm: "Facebook", Claim: "Name", Value: userInfo.Name})

	return principals, nil
}
