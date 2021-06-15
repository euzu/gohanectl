package app

import (
	"github.com/coreos/go-systemd/daemon"
	"github.com/rs/zerolog/log"
	"time"
)

// sytemd monitoring
// Type=notify
func runSystemdMonitoring() {
	//readyness
	sent, err := daemon.SdNotify(false, "READY=1")
	if !sent && err != nil {
		log.Error().Msgf("Systemd monitoring failed to notify %v", err)
	}

	// liveness
	interval, err := daemon.SdWatchdogEnabled(false)
	if err == nil && interval != 0 {
		log.Info().Msgf("Watchdog enabled interval before inactivity is: %d", interval)
		ticker := time.NewTicker(interval / 3)
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-done:
					return
				case _ = <-ticker.C:
					daemon.SdNotify(false, "WATCHDOG=1")
				}
			}
		}()
	} else {
		if err != nil {
			log.Error().Msgf("Watchdog could not be enabled: %v", err)
		} else {
			log.Error().Msg("Watchdog could not be enabled")
		}
	}
}
