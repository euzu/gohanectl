package rest

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
	"gohanectl/hanectl/auth"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"golang.org/x/net/context"
	"net/http"
)

type contextKey struct {
	name string
}

// Context keys
var (
	ClaimsCtxKey = &contextKey{"Claims"}
)

const (
	UsernameClaimKey    = "username"
	AuthoritiesClaimKey = "authorities"
)

type AuthToken struct {
	Token string `json:"token"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authApi(router chi.Router, jwtAuth *jwtauth.JWTAuth, userService model.IUserService) {
	router.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var authRequest AuthRequest
		err := render.DecodeJSON(r.Body, &authRequest)
		_ = r.Body.Close()
		if err == nil && authRequest.Username != "" && authRequest.Password != "" {
			token := authorizeUser(&authRequest, jwtAuth, userService)
			if token != nil {
				render.JSON(w, r, token)
			} else {
				http.Error(w, fmt.Sprintf("%s %s", http.StatusText(http.StatusUnauthorized), err), http.StatusUnauthorized)
			}
		} else {
			http.Error(w, fmt.Sprintf("%s %s", http.StatusText(http.StatusBadRequest), err), http.StatusBadRequest)
		}
	})
	router.Get("/auth/authenticated", func(w http.ResponseWriter, r *http.Request) {
		reqToken := jwtauth.TokenFromHeader(r)
		if utils.IsNotBlank(reqToken) {
			if token, err := jwtAuth.Decode(reqToken); err == nil {
				if err2 := jwt.Validate(token); err2 == nil {
					render.JSON(w, r, true)
					return
				}
			}
		}
		render.JSON(w, r, false)
	})
	//router.Post("/pwdreset", func(w http.ResponseWriter, r *http.Request) {
	//	render.HTML(w, r, "todo")
	//})
}

func authorizeUser(authRequest *AuthRequest, jwtAuth *jwtauth.JWTAuth, userService model.IUserService) *AuthToken {

	user, err := userService.FindByUsername(authRequest.Username)
	if err == nil && user != nil {
		if user.Enabled && auth.CheckPasswordHash(authRequest.Password, user.Password) {
			user.Password = ""
			_, token, _ := jwtAuth.Encode(model.Dictionary{UsernameClaimKey: user.Username, AuthoritiesClaimKey: user.Authorities})
			return &AuthToken{Token: token}
		}
	}
	return nil
}

func jwtAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if err := jwt.Validate(token); err != nil {
			http.Error(w, jwtauth.ErrorReason(err).Error(), 401)
			return
		}

		// Token is authenticated, pass it through
		ctx := context.WithValue(r.Context(), ClaimsCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//func adminAuthenticator(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		claims, _ := r.Context().Value(ClaimsCtxKey).(jwt.MapClaims)
//		if authorities, ok := (claims)[AuthoritiesClaimKey]; ok {
//			if auths, ok := authorities.([]interface{}); ok {
//				for _, item := range auths {
//					if item.(string) == auth.AdminRole {
//						next.ServeHTTP(w, r)
//						return
//					}
//				}
//			}
//		}
//		http.Error(w, http.StatusText(401), 401)
//	})
//}
//
//func getCurrentUsername(r *http.Request) (string, error) {
//	if claims, ok := r.Context().Value(ClaimsCtxKey).(jwt.MapClaims); ok {
//		if username, ok := (claims)[UsernameClaimKey]; ok {
//			return username.(string), nil
//		}
//	}
//	return "", errors.New("no username found")
//}

func getCurrentUser(r *http.Request) (string, []string, error) {
	if claims, ok := r.Context().Value(ClaimsCtxKey).(model.Dictionary); ok {
		var uname string
		var auths []string
		if username, ok := (claims)[UsernameClaimKey]; ok {
			uname = username.(string)
		}
		if authorities, ok := (claims)[AuthoritiesClaimKey]; ok {
			authoritiesList := authorities.([]interface{})
			for i := range authoritiesList {
				auths = append(auths, authoritiesList[i].(string))
			}
		}
		return uname, auths, nil
	}
	return "", nil, errors.New("no username found")
}
