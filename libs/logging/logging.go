/*
 * Author: fasion
 * Created time: 2023-06-29 09:46:29
 * Last Modified by: fasion
 * Last Modified time: 2023-06-29 10:55:11
 */

package logging

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/fasionchan/goutils/stl"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "name",
	MessageKey:     "message",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

type DynamicEncoder struct {
	zapcore.Encoder
}

func NewDynamicEncoder() *DynamicEncoder {
	return &DynamicEncoder{
		Encoder: zapcore.NewJSONEncoder(defaultEncoderConfig),
	}
}

func (dynamic *DynamicEncoder) Dup() *DynamicEncoder {
	return stl.Dup(dynamic)
}

type CloserFunc func() error

type DynamicWriteSyncer struct {
	zapcore.WriteSyncer
	closer CloserFunc
	mutex  sync.Mutex
}

func NewDynamicWriteSyncer() *DynamicWriteSyncer {
	return &DynamicWriteSyncer{
		WriteSyncer: zapcore.AddSync(os.Stdout),
	}
}

func (dynamic *DynamicWriteSyncer) Dup() *DynamicWriteSyncer {
	return stl.Dup(dynamic)
}

func (dynamic *DynamicWriteSyncer) UseCustom(writeSyncer zapcore.WriteSyncer, closer CloserFunc) error {
	dynamic.mutex.Lock()
	defer dynamic.mutex.Unlock()

	if dynamic.closer != nil {
		if err := dynamic.closer(); err != nil {
			return err
		}
	}

	dynamic.WriteSyncer = writeSyncer
	dynamic.closer = closer

	return nil
}

func (dynamic *DynamicWriteSyncer) UseCustomWriter(w io.Writer) error {
	return dynamic.UseCustom(zapcore.AddSync(w), nil)
}

func (dynamic *DynamicWriteSyncer) UseCustomWriteCloser(wc io.WriteCloser) error {
	return dynamic.UseCustom(zapcore.AddSync(wc), wc.Close)
}

func (dynamic *DynamicWriteSyncer) UseCustomFileWriteSyncer(path string, maxSize, maxAge, maxBackups int, localTime, compress bool) error {
	// create parent directories
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	return dynamic.UseCustomWriteCloser(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  localTime,
		Compress:   compress,
	})
}

func (dynamic *DynamicWriteSyncer) UseFileWriteSyncer(path string) error {
	return dynamic.UseCustomFileWriteSyncer(
		// path
		path,
		// max size in megabytes
		20,
		// max age in days
		7,
		// max backups
		5,
		// use local time
		true,
		// compress
		true,
	)
}

type LoggerCreator struct {
	*DynamicWriteSyncer
	*DynamicEncoder
	zap.AtomicLevel
}

func NewLoggerCreator() *LoggerCreator {
	return &LoggerCreator{
		DynamicWriteSyncer: NewDynamicWriteSyncer(),
		DynamicEncoder:     NewDynamicEncoder(),
		AtomicLevel:        zap.NewAtomicLevelAt(zapcore.InfoLevel),
	}
}

func (creator *LoggerCreator) Dup() *LoggerCreator {
	return stl.Dup(creator)
}

func (creator *LoggerCreator) WithLevel(level zapcore.Level) *LoggerCreator {
	creator.SetLevel(level)
	return creator
}

func (creator *LoggerCreator) NewLoggerContainer() *LoggerContainer {
	return NewLoggerContainer(creator.DynamicEncoder.Dup(), creator.DynamicWriteSyncer.Dup(), creator.Level())
}

func (creator *LoggerCreator) NewLogger() *zap.Logger {
	return creator.NewLoggerContainer().GetLogger()
}

func (creator *LoggerCreator) NewStaticLogger() *zap.Logger {
	return zap.New(zapcore.NewCore(creator.DynamicEncoder.Encoder, creator.DynamicWriteSyncer.WriteSyncer, zap.NewAtomicLevelAt(creator.Level())))
}

type LoggerContainer struct {
	*DynamicEncoder
	*DynamicWriteSyncer
	zap.AtomicLevel
	logger *zap.Logger
}

func NewLoggerContainer(encoder *DynamicEncoder, writeSyncer *DynamicWriteSyncer, level zapcore.Level) (container *LoggerContainer) {
	levelEnabler := zap.NewAtomicLevelAt(level)
	return &LoggerContainer{
		DynamicEncoder:     encoder,
		DynamicWriteSyncer: writeSyncer,
		AtomicLevel:        levelEnabler,
		logger:             zap.New(zapcore.NewCore(encoder, writeSyncer, levelEnabler)),
	}
}

func (container *LoggerContainer) GetLogger() *zap.Logger {
	return container.logger
}

var defaultLoggerCreator = NewLoggerCreator()
var defaultLoggerContainer = defaultLoggerCreator.NewLoggerContainer()

var GetLogger = defaultLoggerContainer.GetLogger
var SetLoggerLevel = defaultLoggerContainer.SetLevel
