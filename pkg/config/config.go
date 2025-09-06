package config

import (
	"path/filepath"
)

type Config struct {
	Debug          bool
	FileSpacecmd   string
	FileUyuni      string
	LogCompression bool
	LogFile        string
	LogMaxAge      int
	LogMaxSize     int
	Name           string
	RetryCount     int
}

// Defaults. TODO: override them from environment variables
const (
	fileSpacecmd   = "/root/.spacecmd/config"
	fileUyuni      = "/opt/uyunihub/uyunihub.yaml"
	logCompression = true
	logDir         = "/var/log/mlmtool"
	logMaxAge      = 2 // in MB
	logMaxSize     = 2 // in MB
	retrycount     = 5
)

func New(name string, debug bool) *Config {
	return &Config{
		Debug:          debug,
		FileSpacecmd:   fileSpacecmd,
		FileUyuni:      fileUyuni,
		LogCompression: logCompression,
		LogFile:        filepath.Join(logDir, name+".log"),
		LogMaxAge:      logMaxAge,
		LogMaxSize:     logMaxSize,
		Name:           name,
		RetryCount:     retrycount,
	}
}
