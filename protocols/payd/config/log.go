package config

import (
	"github.com/labstack/gommon/log"
)

// SetupLog will setup the logger.
func SetupLog(cfg *Logging) {
	switch cfg.Level {
	case LogDebug:
		log.SetLevel(log.DEBUG)
	case LogError:
		log.SetLevel(log.ERROR)
	case LogWarn:
		log.SetLevel(log.WARN)
	default:
		log.SetLevel(log.INFO)
	}
}
