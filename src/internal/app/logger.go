package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	loglevel, err := logrus.ParseLevel(conf.Loglevel)
	if err == nil {
		log.SetLevel(loglevel)
	} else {
		log.SetLevel(logrus.ErrorLevel)
		log.Error("Can`t parse LOG_LEVEL. Used default value: LOG_LEVEL=error")
	}

	log.Infof("Used json logger and loglevel: %s", conf.Loglevel)
}
