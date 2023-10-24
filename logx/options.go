package logx

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// makeTimeEncoder makes timestamp with human-readable style.
func makeTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000 -07:00"))
}

type configurer func(*Config)

// FileConfig represents file to record log in toml/json.
type FileConfig struct {
	Path string `toml:"path" json:"path"`
	Name string `toml:"name" json:"name"`

	// Along with lumberjack:
	// https://github.com/natefinch/lumberjack/blob/v2.0/lumberjack.go
	MaxSize    int  `toml:"maxsize" json:"maxsize"`
	MaxAge     int  `toml:"maxage" json:"maxage"`
	MaxBackups int  `toml:"maxbackups" json:"maxbackups"`
	Compress   bool `toml:"compress" json:"compress"`
}

func NewFileConfig(path, name string, maxSize, maxAge, maxBackups int, compress bool) *FileConfig {
	return &FileConfig{
		Path:       path,
		Name:       name,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		Compress:   compress,
	}
}

func (f *FileConfig) BuildLogger() (*lumberjack.Logger, error) {
	filePath := filepath.Join(f.Path, f.Name)
	if stat, err := os.Stat(filePath); err == nil {
		if stat.IsDir() {
			return nil, errors.New("log file cannot be directory")
		}
	}

	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    f.MaxSize,
		MaxAge:     f.MaxAge,
		MaxBackups: f.MaxBackups,
		Compress:   f.Compress,
	}, nil
}

// Config represents log configurations in toml/json.
type Config struct {
	Lvl                 Level       `toml:"level" json:"level"`
	OsStdout            bool        `toml:"os-stdout" json:"os-stdout"`
	File                *FileConfig `toml:"file" json:"file"`
	Format              string      `toml:"format" json:"format"`
	DisableTimestamp    bool        `toml:"disable-timestamp" json:"disable-timestamp"`
	Development         bool        `toml:"development" json:"development"`
	DisableCaller       bool        `toml:"disable-caller" json:"disable-caller"`
	DisableStacktrace   bool        `toml:"disable-stacktrace" json:"disable-stacktrace"`
	DisableErrorVerbose bool        `toml:"disable-error-verbose" json:"disable-error-verbose"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Lvl:      LvlDebug,
		Format:   FormatJson,
		OsStdout: true,
	}
}

func (c *Config) BuildZapEncoderConfig() zapcore.EncoderConfig {
	encConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     makeTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	if c.DisableTimestamp {
		encConfig.TimeKey = ""
	}

	return encConfig
}

func (c *Config) BuildZapEncoder() zapcore.Encoder {
	config := c.BuildZapEncoderConfig()
	switch c.Format {
	case FormatJson:
		return zapcore.NewJSONEncoder(config)
	case FormatConsole:
		return zapcore.NewConsoleEncoder(config)
	case FormatText:
		// TODO:
		panic("it is still in todo list")
	default:
		panic("invalid format")
	}
}

func (c *Config) BuildZapWriteSyncer() ([]zapcore.WriteSyncer, error) {
	syncers := []zapcore.WriteSyncer{}
	if c.File != nil && len(c.File.Name) > 0 {
		// build lumberjack logger
		lg, err := c.File.BuildLogger()
		if err != nil {
			return nil, err
		}

		syncers = append(syncers, zapcore.AddSync(lg))
	}

	if c.OsStdout {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	return syncers, nil
}

func (c *Config) BuildZapOptions() []zap.Option {
	opts := []zap.Option{}

	stackLevel := zap.ErrorLevel
	if c.Development {
		opts = append(opts, zap.Development())
		stackLevel = zap.WarnLevel
	}

	if !c.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if !c.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	return opts
}

func WithLevel(lvl Level) configurer {
	return func(c *Config) {
		c.Lvl = lvl
	}
}

func WithFile(path, name string, maxSize, maxAge, maxBackups int, compress bool) configurer {
	return func(c *Config) {
		c.File = NewFileConfig(path, name, maxSize, maxAge, maxBackups, compress)
	}
}

func WithConsole() configurer {
	return func(c *Config) {
		c.OsStdout = true
	}
}

func WithFormat(format string) configurer {
	return func(c *Config) {
		c.Format = format
	}
}

func WithDisableTimestamp() configurer {
	return func(c *Config) {
		c.DisableTimestamp = true
	}
}

func WithDisableCaller() configurer {
	return func(c *Config) {
		c.DisableCaller = true
	}
}

func WithDisableStacktrace() configurer {
	return func(c *Config) {
		c.DisableStacktrace = true
	}
}

func WithDisableErrorVerbose() configurer {
	return func(c *Config) {
		c.DisableErrorVerbose = true
	}
}
