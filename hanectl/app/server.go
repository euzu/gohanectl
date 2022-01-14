package app

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/rest"
	"gohanectl/hanectl/websocket"
	"net"
	"net/http"
	"time"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func createListener(cfg config.IConfiguration) *net.Listener {
	host := cfg.GetStr(config.ListenHost, config.DefHost)
	port := cfg.GetInt(config.ListenPort, config.DefPort)
	serverUri := fmt.Sprintf("%s:%d", host, port)
	lhttp, err := net.Listen("tcp", serverUri)
	if err != nil {
		log.Fatal().Msgf("Failed to start web server: %v", err)
	}
	log.Info().Msgf("Server listening at http://%s", serverUri)
	lhttp = tcpKeepAliveListener{lhttp.(*net.TCPListener)}
	return &lhttp
}

func startServer(cfg config.IConfiguration, serviceFactory model.IServiceFactory) *net.Listener {
	router := rest.Routes(cfg, serviceFactory)
	websocket.StartWebSocket(router, serviceFactory.GetConfigService().SetWebsocketStatus)
	startWebfileServer(cfg, router)
	rest.PrintRoutes(router)
	lhttp := createListener(cfg)
	log.Fatal().Err(http.Serve(*lhttp, router))
	return lhttp
}
