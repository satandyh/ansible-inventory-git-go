// got from gist. really cool. thanks bro
// https://gist.github.com/panta/2530672ca641d953ae452ecb5ef79d7d

package logging

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
What need
1. log levels
2. log to stdout with syslog compatible format
3. log to file with json format
4. rotation log files
*/

/* USAGE EXAMPLE

import (
	logging "github.com/satandyh/ansible-inventory-git-go/internal/logger"
)
var logConfig = logging.LogConfig{
	ConsoleLoggingEnabled: true,
	EncodeLogsAsJson:      true,
	FileLoggingEnabled:    true,
	Directory:             "./data",
	Filename:              "lovely_app.log",
	MaxSize:               10,
	MaxBackups:            7,
	MaxAge:                7,
	LogLevel:              6,
}
var logger = logging.Configure(logConfig)
func main() {
	logger.Info().
		Str("module", "main").
		Msg("Another instance already work.")
	logger.Error().
		Str("module", "main").
		Err(err).
		Msg("Cannot change lockfile mtime")
	log.Fatal().
		Err(err).
		Str("module", "config").
		Msg("")
}
*/

// Configuration for logging
type LogConfig struct {
	// Enable console logging
	ConsoleLoggingEnabled bool
	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
	// LogLevel as in https://en.wikipedia.org/wiki/Syslog
	LogLevel int
}

type Logger struct {
	*zerolog.Logger
}

func Configure(config LogConfig) *Logger {
	ex, ex_err := os.Executable()
	if ex_err != nil {
		panic(ex_err)
	}
	selfName := filepath.Base(ex)

	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	logger := zerolog.
		New(mw).
		With().
		Timestamp().
		Str("process", selfName).
		Logger()

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile(config LogConfig) io.Writer {
	if err := os.MkdirAll(config.Directory, 0750); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
