package logx

import (
	"errors"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func makeZapLogLevel(lvl Level) zapcore.Level {
	switch lvl {
	case LvlDebug:
		return zapcore.DebugLevel
	case LvlInfo:
		return zapcore.InfoLevel
	case LvlWarn:
		return zapcore.WarnLevel
	case LvlError:
		return zapcore.ErrorLevel
	case LvlDPanic:
		return zapcore.DPanicLevel
	case LvlPanic:
		return zapcore.PanicLevel
	case LvlFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InvalidLevel
	}
}

type zapLogger struct {
	leveledLogger  sync.Map
	leveledSLogger sync.Map
	level          Level
}

func InitZapLogger(config *Config, opts ...zap.Option) (*zapLogger, error) {
	if config == nil {
		return nil, errors.New("no configuration to initialize zap logger")
	}

	syncers, err := config.BuildZapWriteSyncer()
	if err != nil {
		return nil, err
	}

	//level debug
	core := zapcore.NewCore(
		config.BuildZapEncoder(),
		zap.CombineWriteSyncers(syncers...),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	if opts != nil {
		opts = append(opts, config.BuildZapOptions()...)
	}

	debugLogger := zap.New(core, opts...)

	levels := []Level{
		LvlDebug,
		LvlInfo,
		LvlWarn,
		LvlError,
		LvlDPanic,
		LvlPanic,
		LvlFatal,
	}

	lg := zapLogger{level: config.Lvl}
	for _, level := range levels {
		promotedLogger := debugLogger.WithOptions(zap.IncreaseLevel(makeZapLogLevel(level)))
		lg.leveledLogger.Store(level, promotedLogger)
		lg.leveledSLogger.Store(level, promotedLogger.Sugar())
	}

	return &lg, nil
}

func (z *zapLogger) Logger() *zap.Logger {
	lg, _ := z.leveledLogger.Load(z.level)
	return lg.(*zap.Logger)
}

func (z *zapLogger) DebugLogger() *zap.Logger {
	lg, _ := z.leveledLogger.Load(LvlDebug)
	return lg.(*zap.Logger)
}

func (z *zapLogger) InfoLogger() *zap.Logger {
	lg, _ := z.leveledLogger.Load(LvlInfo)
	return lg.(*zap.Logger)
}

func (z *zapLogger) WarnLogger() *zap.Logger {
	lg, _ := z.leveledLogger.Load(LvlWarn)
	return lg.(*zap.Logger)
}

func (z *zapLogger) ErrorLogger() *zap.Logger {
	lg, _ := z.leveledLogger.Load(LvlError)
	return lg.(*zap.Logger)
}

func (z *zapLogger) FatalLogger() *zap.Logger {
	lg, _ := z.leveledLogger.Load(LvlFatal)
	return lg.(*zap.Logger)
}

func (z *zapLogger) SLogger() *zap.SugaredLogger {
	lg, _ := z.leveledSLogger.Load(z.level)
	return lg.(*zap.SugaredLogger)
}

func (z *zapLogger) DebugSLogger() *zap.SugaredLogger {
	lg, _ := z.leveledSLogger.Load(LvlDebug)
	return lg.(*zap.SugaredLogger)
}

func (z *zapLogger) InfoSLogger() *zap.SugaredLogger {
	lg, _ := z.leveledSLogger.Load(LvlInfo)
	return lg.(*zap.SugaredLogger)
}

func (z *zapLogger) WarnSLogger() *zap.SugaredLogger {
	lg, _ := z.leveledSLogger.Load(LvlWarn)
	return lg.(*zap.SugaredLogger)
}

func (z *zapLogger) ErrorSLogger() *zap.SugaredLogger {
	lg, _ := z.leveledSLogger.Load(LvlError)
	return lg.(*zap.SugaredLogger)
}

func (z *zapLogger) FatalSLogger() *zap.SugaredLogger {
	lg, _ := z.leveledSLogger.Load(LvlFatal)
	return lg.(*zap.SugaredLogger)
}

func (z *zapLogger) GetLevel() Level {
	return z.level
}

func (z *zapLogger) Sync() error {
	var retErr error
	z.leveledLogger.Range(func(key, value interface{}) bool {
		lg := value.(*zap.Logger)
		if err := lg.Sync(); err != nil {
			retErr = err
			return false
		}

		return true
	})

	if retErr != nil {
		return retErr
	}

	z.leveledSLogger.Range(func(key, value interface{}) bool {
		lg := value.(*zap.SugaredLogger)
		if err := lg.Sync(); err != nil {
			retErr = err
			return false
		}

		return true
	})

	return retErr
}

func NewZapLogger(configs ...configurer) (*zapLogger, error) {
	initConfig := &Config{}
	if configs == nil {
		initConfig = NewDefaultConfig()
	} else {
		for _, config := range configs {
			config(initConfig)
		}
	}

	return InitZapLogger(initConfig)
}

func NewLogger(configs ...configurer) (*zap.Logger, error) {
	zapLg, err := NewZapLogger(configs...)
	if err != nil {
		return nil, err
	}

	return zapLg.Logger(), nil
}

func NewSugaredLogger(configs ...configurer) (*zap.SugaredLogger, error) {
	zapLg, err := NewZapLogger(configs...)
	if err != nil {
		return nil, err
	}

	return zapLg.SLogger(), nil
}
