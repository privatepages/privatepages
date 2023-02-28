package logger

import (
	"os"

	"package/main/internal/config"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var conf *config.Config

func init() {
	conf = config.Cfg
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)

	loglevel, err := logrus.ParseLevel(conf.Loglevel)
	if err == nil {
		Log.SetLevel(loglevel)
	} else {
		Log.SetLevel(logrus.ErrorLevel)
		Log.Error("Can`t parse LOG_LEVEL. Used default value: LOG_LEVEL=error")
	}

	Log.Infof("Used json logger and loglevel: %s", conf.Loglevel)
}
