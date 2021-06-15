package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"net/http"
)

const (
	id    = "id"
	token = "token"
)

func apiRoutesV1(cfg config.IConfiguration, serviceFactory model.IServiceFactory) *chi.Mux {
	jwtAuth := createJwtAuth(cfg)
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(
			render.SetContentType(render.ContentTypeJSON), // set content-type header as application/json
		)
		authApi(r, jwtAuth, serviceFactory.GetUserService())
	})

	router.Group(func(r chi.Router) {
		r.Use(
			render.SetContentType(render.ContentTypeJSON), // set content-type header as application/json
			jwtauth.Verifier(jwtAuth),
			jwtAuthenticator,
		)
		configApi(r, serviceFactory)
		deviceApi(r, serviceFactory)
	})

	router.Group(func(r chi.Router) {
		r.Use(
			render.SetContentType(render.ContentTypeJSON), // set content-type header as application/json
		)
		commandApi(r, serviceFactory, cfg)
	})

	return router
}

func configApi(rapi chi.Router, serviceFactory model.IServiceFactory) {
	rapi.Get("/status", func(w http.ResponseWriter, r *http.Request) { serverStatus(w, r, serviceFactory) })
	rapi.Get("/config/room", func(w http.ResponseWriter, r *http.Request) { configRooms(w, r, serviceFactory) })
	rapi.Post("/user/setting", func(w http.ResponseWriter, r *http.Request) { setUserSettings(w, r, serviceFactory) })
	rapi.Get("/user/setting", func(w http.ResponseWriter, r *http.Request) { getUserSettings(w, r, serviceFactory) })
}

func deviceApi(rapi chi.Router, serviceFactory model.IServiceFactory) {
	rapi.Get("/devices", func(w http.ResponseWriter, r *http.Request) { deviceConfig(w, r, serviceFactory) })
	rapi.Get("/devices/status", func(w http.ResponseWriter, r *http.Request) { deviceStates(w, r, serviceFactory) })
	rapi.Get(fmt.Sprintf("/device/{%s}", id), func(w http.ResponseWriter, r *http.Request) { deviceState(w, r, serviceFactory) })
	rapi.Post("/device", func(w http.ResponseWriter, r *http.Request) { deviceCommand(w, r, serviceFactory) })
}

func commandApi(rapi chi.Router, serviceFactory model.IServiceFactory, cfg config.IConfiguration) {
	rapi.Get(fmt.Sprintf("/cmd/{%s}/{%s}", token, id), func(w http.ResponseWriter, r *http.Request) { command(w, r, serviceFactory, cfg) })
	rapi.Get(fmt.Sprintf("/groupcmd/{%s}/{%s}", token, id), func(w http.ResponseWriter, r *http.Request) { groupCommand(w, r, serviceFactory, cfg) })
}
