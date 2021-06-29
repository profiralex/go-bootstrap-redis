package logs

import (
	"github.com/profiralex/go-bootstrap-redis/pkg/config"
	log "github.com/sirupsen/logrus"
)

func Init(cfg config.Config) {
	debugLevel, err := log.ParseLevel(cfg.AppConfig.DebugLevel)
	if err != nil {
		log.Warnf("Unknown debug level %s defaulting to warning level", cfg.AppConfig.DebugLevel)
		debugLevel = log.WarnLevel
	}
	log.SetLevel(debugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}
