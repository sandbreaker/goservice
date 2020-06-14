package config

import (
	"os"
	"time"

	"github.com/go-ini/ini"

	"github.com/sandbreaker/goservice/log"
)

var StaticConfig *Config

const (
	DefaultStaticCfgFile = ".cfg/static.cfg"
)

func loadStatic(env, role string, isDebugOn bool) {
	staticCfgFile := os.Getenv("GOSERVICE_STATIC_CONFIG")
	if staticCfgFile == "" {
		staticCfgFile = DefaultStaticCfgFile
	}

	file, err := ini.LoadSources(
		ini.LoadOptions{IgnoreInlineComment: true},
		staticCfgFile,
	)
	if err != nil {
		log.Error(err)
	}

	StaticConfig = &Config{
		env:         env,
		role:        role,
		dateUpdated: time.Now(),
		isDebugOn:   isDebugOn,
		file:        file,
	}
}
