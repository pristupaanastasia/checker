package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"time"
)

type Log interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	With(args ...interface{}) Log
	Delete()
}

type Logger struct {
	Log   *zap.SugaredLogger
	file  *os.File
	day   int
	level zapcore.Level
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
	y, m, d := time.Now().Date()
	if file == "" {
		file = "/log/info" + strconv.Itoa(y) + "_" + m.String() + "_" + strconv.Itoa(d) + ".txt"
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
	log := &Logger{Log: l.Sugar(), file: f, day: d, level: zapLevel}
	return log, nil
}
func (log *Logger) checkFile() {
	if time.Now().Day() != log.day {
		y, m, d := time.Now().Date()
		file := "/log/info" + strconv.Itoa(y) + "_" + m.String() + "_" + strconv.Itoa(d) + ".txt"

		f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Log.Fatal(err)
		}
		pe := zap.NewProductionEncoderConfig()

		fileEncoder := zapcore.NewJSONEncoder(pe)

		pe.EncodeTime = zapcore.ISO8601TimeEncoder // The encoder can be customized for each output
		consoleEncoder := zapcore.NewConsoleEncoder(pe)
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(f), log.level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), log.level),
		)

		l := zap.New(core) // Creating the logger
		log.Delete()

		log.file = f
		log.Log = l.Sugar()
	}
}
func (log *Logger) With(v ...interface{}) Log {
	log.checkFile()
	return &Logger{Log: log.Log.With(v...), file: log.file}
}
func (log *Logger) Info(v ...interface{}) {
	log.checkFile()
	log.Log.Info(v...)
}
func (log *Logger) Error(v ...interface{}) {
	log.checkFile()
	log.Log.Error(v...)
}
func (log *Logger) Fatal(v ...interface{}) {
	log.Log.Fatal(v...)
}

func (log *Logger) Delete() {
	log.Log.Sync()
	log.file.Close()
}
