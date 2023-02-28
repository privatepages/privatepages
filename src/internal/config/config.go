package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a structure containing configuration fields for this application.
type Config struct {
	Loglevel            string `env:"LOG_LEVEL"              env-default:"error"`
	HTTPListenIPPort    string `env:"HTTP_LISTEN"            env-default:":80"`
	HTTPUploadMaxSize   int    `env:"HTTP_UPLOAD_MAX_SIZE"   env-default:"10"`
	APISecret           string `env:"API_SECRET"` // stored in secret
	ArtifactStoragePath string `env:"ARTIFACT_STORAGE_PATH"  env-default:"/data"`
	CheckIndexPage      bool   `env:"CHECK_INDEX_PAGE"       env-default:"true"`
}

// Cfg contains pointer to config object
var Cfg *Config

func init() {
	Cfg = &Config{}
	err := cleanenv.ReadEnv(Cfg)
	if err != nil {
		fmt.Printf("Something went wrong while reading the configuration: %s", err)
		os.Exit(1)
	}
}
