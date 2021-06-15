package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"net/http"
	"strings"
)

func Routes(cfg config.IConfiguration, serviceFactory model.IServiceFactory) *chi.Mux {
	router := chi.NewRouter()

	corsConfig := createCorsConfig(cfg)
	router.Use(
		corsConfig.Handler,
		//middleware.Logger,
		utils.HttpLogger, //middleware.Logger,                // Log api request calls
		//middleware. DefaultCompress,   // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	router.Group(func(r chi.Router) {
		r.Route("/api", func(r2 chi.Router) {
			r2.Mount("/v1", apiRoutesV1(cfg, serviceFactory))
		})
	})

	return router
}

func PrintRoutes(router *chi.Mux) {
	walkFunc := func(method string, router string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route := strings.Replace(router, "*/", "", -1)
		path := fmt.Sprintf("%s %s", method, route)
		log.Debug().Msg(path)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panic().Err(err)
	}
}

func createCorsConfig(cfg config.IConfiguration) *cors.Cors {
	cfgAllowedOrigins := cfg.GetList(config.CorsAllowedOrigins)
	var allowedOrigins []string
	if cfgAllowedOrigins == nil {
		allowedOrigins = []string{}
	} else {
		allowedOrigins = make([]string, len(cfgAllowedOrigins))
		for i := range cfgAllowedOrigins {
			allowedOrigins[i] = cfgAllowedOrigins[i].(string)
		}
	}

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Language", "X-CSRF-Token", "X-Authorization", "Sec-WebSocket-Protocol"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
