package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		debug          bool
		expectedConfig *Config
	}{
		{
			name:  "test1",
			debug: true,
			expectedConfig: &Config{
				Debug:          true,
				FileSpacecmd:   fileSpacecmd,
				FileUyuni:      fileUyuni,
				LogCompression: logCompression,
				LogFile:        filepath.Join(logDir, "test1.log"),
				LogMaxAge:      logMaxAge,
				LogMaxSize:     logMaxSize,
				Name:           "test1",
				RetryCount:     retrycount,
			},
		},
		{
			name:  "test2",
			debug: false,
			expectedConfig: &Config{
				Debug:          false,
				FileSpacecmd:   fileSpacecmd,
				FileUyuni:      fileUyuni,
				LogCompression: logCompression,
				LogFile:        filepath.Join(logDir, "test2.log"),
				LogMaxAge:      logMaxAge,
				LogMaxSize:     logMaxSize,
				Name:           "test2",
				RetryCount:     retrycount,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := New(tt.name, tt.debug)
			assert.Equal(t, tt.expectedConfig, config)
		})
	}
}
