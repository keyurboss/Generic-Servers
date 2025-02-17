package apis

import (
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/env"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/utility"
)

var AccessTokenService *utility.TokenService
var RefreshTokenService *utility.TokenService

func init() {
	AccessTokenService = &utility.TokenService{
		SigningKey: []byte(env.Env.ACCESS_TOKEN_KEY),
	}
	RefreshTokenService = &utility.TokenService{
		SigningKey: []byte(env.Env.REFRESH_TOKEN_KEY),
	}
	println("Token Service Initialized")
}
