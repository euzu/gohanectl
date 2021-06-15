package rest

import (
	"github.com/go-chi/jwtauth"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/utils"
)

func createJwtAuth(cfg config.IConfiguration) *jwtauth.JWTAuth {
	secret := cfg.GetStr(config.JwtSecret, "")
	if secret == "" {
		secret = utils.RandomString(32)
	}
	return jwtauth.New("HS256", []byte(secret), nil)
}
