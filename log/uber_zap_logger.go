package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(options ...Option) Logger {
	zapConf := zap.NewDevelopmentConfig()
	zapConf.DisableStacktrace = true
	zapConf.DisableCaller = true
	zapConf.EncoderConfig.EncodeName = func(name string, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(fmt.Sprintf("[%s]", name))
	}

	l, _ := zapConf.Build()
	logger := l.Sugar()

	// apply options
	c := &config{}
	for _, option := range options {
		option(c)
	}
	logger = logger.Named(c.namespace).With(c.fields...)

	return &zapLogger{logger: logger, minLevel: c.minLevel}
}

var _ Logger = &zapLogger{}

type zapLogger struct {
	logger   *zap.SugaredLogger
	minLevel Level
}

func (z *zapLogger) Debug(msg string, fields ...interface{}) {
	if DEBUG >= z.minLevel {
		z.logger.Debugw(z.buildMessage(msg, nil), fields...)
	}
}

func (z *zapLogger) Info(msg string, fields ...interface{}) {
	if INFO >= z.minLevel {
		z.logger.Infow(z.buildMessage(msg, nil), fields...)
	}
}

func (z *zapLogger) Error(err error, msg string, fields ...interface{}) {
	if ERROR >= z.minLevel {
		z.logger.Errorw(z.buildMessage(msg, err), fields...)
	}
}

func (z *zapLogger) Derive(options ...Option) Logger {
	// apply options
	c := &config{}
	for _, option := range options {
		option(c)
	}

	logger := z.logger.Named(c.namespace).With(c.fields...)
	minLevel := z.minLevel
	if c.minLevel != unset {
		minLevel = c.minLevel
	}

	return &zapLogger{logger: logger, minLevel: minLevel}
}

func (z *zapLogger) buildMessage(msg string, err error) string {
	if err != nil {
		msg = msg + ": " + err.Error()
	}
	return msg
}
