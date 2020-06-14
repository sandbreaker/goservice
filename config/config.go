package config

import (
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/palantir/stacktrace"
	"github.com/palantir/stacktrace/cleanpath"
	"github.com/sirupsen/logrus"

	"github.com/sandbreaker/goservice/log"
)

var isDebugOn bool

const (
	DefaultEnv  = "dev"
	DefaultRole = "default"

	EnvProduction = "prod"
	EnvStaging    = "staging"
	EnvDev        = "dev"
)

type Config struct {
	env         string
	role        string
	dateUpdated time.Time
	isDebugOn   bool
	file        *ini.File
}

func init() {
	// environment
	env := os.Getenv("GOSERVICE_ENV")
	if env == "" {
		env = DefaultEnv
	}

	// logger config
	log.Logger.Formatter = &logrus.TextFormatter{}
	if env == EnvProduction {
		log.Logger.Level = logrus.InfoLevel
	} else {
		log.Logger.Out = os.Stdout
		log.Logger.Level = logrus.DebugLevel
	}

	// role
	role := os.Getenv("GOSERVICE_ROLE")
	if role == "" {
		role = DefaultRole
	}

	// check and see if overrid debug on
	debugOn()

	// errors and stack traces cleanse
	stacktrace.CleanPath = func(path string) string {
		path = cleanpath.RemoveGoPath(path)
		path = strings.TrimPrefix(path, "github.com/")
		path = strings.TrimPrefix(path, "bitbucket.org/")
		return path
	}
	stacktrace.DefaultFormat = stacktrace.FormatFull

	// load static
	loadStatic(env, role, isDebugOn)

	// load dynamic

	// load version
}

func debugOn() {
	if os.Getenv("GOSERVICE_DEBUG") == "true" {
		log.Logger.Out = os.Stdout
		log.Logger.Level = logrus.DebugLevel
		isDebugOn = true
	}
}

func (c *Config) GetEnv() string {
	return c.env
}

func (c *Config) GetRole() string {
	return c.role
}

func (c *Config) IsDebugOn() bool {
	return c.isDebugOn
}

func (c *Config) GetInt(key string, defaultValue int) int {
	if c.file == nil || !c.file.Section(c.env).HasKey(key) {
		log.Debug("No key found = ", key)
		return defaultValue
	}

	value, err := c.file.Section(c.env).Key(key).Int()
	if err != nil {
		log.Error(err)
		return defaultValue
	}

	return value
}

func (c *Config) GetString(key, defaultValue string) string {
	if c.file == nil || !c.file.Section(c.env).HasKey(key) {
		log.Debug("No key found = ", key)
		return defaultValue
	}

	return c.file.Section(c.env).Key(key).String()
}

func (c *Config) GetBool(key string, defaultValue bool) bool {
	if c.file == nil || !c.file.Section(c.env).HasKey(key) {
		log.Debug("No key found = ", key)
		return defaultValue
	}

	value, err := c.file.Section(c.env).Key(key).Bool()
	if err != nil {
		log.Error(err)
		return defaultValue
	}

	return value
}
