// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm)
// aahframework.org/examples source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package security

import (
	"aahframe.work"
	"aahframe.work/config"
	"aahframe.work/security/authc"
	"aahframework.org/examples/form-based-auth/app/models"
)

var _ authc.Authenticator = (*FormAuthenticationProvider)(nil)

// FormAuthenticationProvider struct implements `authc.Authenticator` interface.
type FormAuthenticationProvider struct {
}

// Init method initializes the FormAuthenticationProvider, this method gets called
// during server start up.
func (fa *FormAuthenticationProvider) Init(cfg *config.Config) error {

	// NOTE: Init is called on application startup

	return nil
}

// GetAuthenticationInfo method is `authc.Authenticator` interface
func (fa *FormAuthenticationProvider) GetAuthenticationInfo(authcToken *authc.AuthenticationToken) (*authc.AuthenticationInfo, error) {

	// Form Auth Values from authcToken
	// 		authcToken.Identity => username
	// 		authcToken.Credential => passowrd

	user := models.FindUserByEmail(authcToken.Identity)
	if user == nil {
		// No subject exists, return nil and error
		return nil, authc.ErrSubjectNotExists
	}

	// User found, now create authentication info and return to the framework
	authcInfo := authc.NewAuthenticationInfo()
	authcInfo.Principals = append(authcInfo.Principals,
		&authc.Principal{
			Value:     user.Email,
			IsPrimary: true,
			Realm:     "inmemory",
		})
	authcInfo.Credential = []byte(user.Password)
	authcInfo.IsLocked = user.IsLocked
	authcInfo.IsExpired = user.IsExpried

	return authcInfo, nil
}

func PostAuthEvent(e *aah.Event) {
	ctx := e.Data.(*aah.Context)

	// Populate session info after authentication
	user := models.FindUserByEmail(ctx.Subject().PrimaryPrincipal().Value)
	ctx.Session().Set("FirstName", user.FirstName)
	ctx.Session().Set("LastName", user.LastName)
}
