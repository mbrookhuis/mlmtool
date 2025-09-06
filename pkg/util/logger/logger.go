package logger

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"mlmtool/pkg/config"
)

func getJSONEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeTime:     zapcore.RFC3339TimeEncoder, // zapcore.ISO8601TimeEncoder
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Undefined or empty keys will be hidden from output
		// TimeKey:        "T",
		CallerKey:      "C",
		LevelKey:       "L",
		MessageKey:     "M",
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func testLogFileCreation(cfg *config.Config) {

	logDir := filepath.Dir(cfg.LogFile)
	if err := os.MkdirAll(logDir, 0700); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	// Attempt to open the logfile to ensure it can be created
	file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("failed to create logfile: %v", err)
	}
	defer file.Close() // Ensure the file handle is closed

}

func getJSONCore(cfg *config.Config) zapcore.Core {
	// Severity
	jsonLevel := zap.InfoLevel
	if cfg.Debug {
		jsonLevel = zapcore.DebugLevel
	}

	testLogFileCreation(cfg)
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(getJSONEncoderConfig()),
		// logfile output and rotation
		zapcore.AddSync(&lumberjack.Logger{
			Filename: cfg.LogFile,
			MaxSize:  cfg.LogMaxSize,
			MaxAge:   cfg.LogMaxAge,
			Compress: cfg.LogCompression,
		}),
		jsonLevel,
	)
}

func getConsoleCore(cfg *config.Config) zapcore.Core {
	// Severity
	consoleLevel := zap.InfoLevel
	if cfg.Debug {
		consoleLevel = zapcore.DebugLevel
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(getConsoleEncoderConfig()),
		zapcore.AddSync(os.Stdout), // Console output to stdout
		consoleLevel,
	)
}

func getLogger(cfg *config.Config, core zapcore.Core) *zap.Logger {
	if debug {
		cfg.Debug = true
	}

	// Create a logger with specific core
	logger := zap.New(core)

	if cfg.Debug {
		// Ensure DPANIC levels do only panic in debug mode
		logger = logger.WithOptions(zap.Development(), zap.AddCaller())
	}

	return logger
}

var debug bool

func init() {
	// add via environment variable enables debug log switch for tests
	debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
}

func New(cfg *config.Config) *zap.Logger {

	// Combine cores
	core := zapcore.NewTee(getJSONCore(cfg), getConsoleCore(cfg))

	logger := getLogger(cfg, core)
	return logger

}

func NewTestingLogger(name string) *zap.Logger {

	cfg := config.New(name, debug)
	// console output only: Does not need a defer logger.Sync() call
	core := getConsoleCore(cfg)
	logger := getLogger(cfg, core)
	logger.Sugar().Debugf("Logger started: %s", cfg.Name)
	return logger
}
