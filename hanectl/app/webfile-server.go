package app

import (
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"net/http"
	"os"
	"strings"
)

func startWebfileServer(cfg config.IConfiguration, r chi.Router) {

	webFilesDir := cfg.GetStr(config.WebFilesDir, config.DefWebFilesDir)
	if _, err := os.Stat(webFilesDir); err != nil {
		log.Fatal().Msgf("Cant find web directory %s", webFilesDir)
	}
	path := config.DefWebUrl
	root := http.Dir(webFilesDir)

	fs := http.StripPrefix(path, http.FileServer(root))
	path += "*"

	// rewrite is for angular direct links
	rewrite := func(path string) string {
		if !strings.Contains(path, ".") {
			return "/"
		}
		return path
	}
	r.Group(func(r2 chi.Router) {
		r2.Get(path, func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = rewrite(r.URL.Path)
			fs.ServeHTTP(w, r)
		})
	})
}
