package pkg

import (
	log "github.com/sirupsen/logrus"
)

func InitLog(config *Config) {
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
