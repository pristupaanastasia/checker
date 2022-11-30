package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Log interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	With(args ...interface{}) Log
	Delete()
}

type Logger struct {
	Log  *zap.SugaredLogger
	file *os.File
}

func Info(log Log, v ...interface{}) {

	log.Info(v...)
}

func Error(log Log, v ...interface{}) {
	log.Error(v...)
}
func Fatal(log Log, v ...interface{}) {
	log.Error(v...)
}
func Delete(log Log) {
	log.Delete()
}
func With(log Log, v ...interface{}) Log {
	return log.With(v...)
}

func New(level string, file string) (Log, error) {

	var err error
	if file == "" {
		file = "info.txt"
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)

	pe.EncodeTime = zapcore.ISO8601TimeEncoder // The encoder can be customized for each output
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	var zapLevel zapcore.Level
	switch level {
	case "DEBUG", "ALL":
		zapLevel = zap.DebugLevel
	case "ERROR":
		zapLevel = zap.ErrorLevel
	case "WARN":
		zapLevel = zap.FatalLevel
	default:
		zapLevel = zap.InfoLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zapLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLevel),
	)

	l := zap.New(core) // Creating the logger
	log := Logger{Log: l.Sugar(), file: f}
	return log, nil
}

func (log Logger) With(v ...interface{}) Log {
	return Logger{Log: log.Log.With(v...), file: log.file}
}
func (log Logger) Info(v ...interface{}) {

	log.Log.Info(v...)
}
func (log Logger) Error(v ...interface{}) {
	log.Log.Error(v...)
}
func (log Logger) Fatal(v ...interface{}) {
	log.Log.Fatal(v...)
}
func (log Logger) Delete() {
	log.Log.Sync()
	log.file.Close()
}
