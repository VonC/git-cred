package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = 99
	INFO  LogLevel = 50
	ERROR LogLevel = 0
)

var (
	loglevels map[string]LogLevel = map[string]LogLevel{
		"debug": DEBUG,
		"info":  INFO,
		"error": ERROR,
	}
)

type LoggerConfig struct {
	Loglevel LogLevel
	Logpath  string
	Prefix   string
	Flag     int
}

type Logger struct {
	Config LoggerConfig
	logger *log.Logger
}

func New() *Logger {
	l := &Logger{}
	l.Config.Loglevel = INFO
	l.Config.Flag = log.Ldate | log.Ltime

	return l
}

func (l *Logger) Logger() *log.Logger {
	return l.logger
}

var once sync.Once
var gl *Logger

func init() {
	once.Do(func() {
		gl = New()
		if err := gl.Init(); err != nil {
			log.Fatal(err)
		}
	})
}

func (l *Logger) SetLevel(level string) {
	if lvl, ok := loglevels[level]; !ok {
		l.Config.Loglevel = INFO
	} else {
		l.Config.Loglevel = lvl
	}
}

func (l *Logger) SetLogpath(path string) {
	l.Config.Logpath = path
}

func (l *Logger) SetPrefix(prefix string) {
	l.Config.Prefix = prefix
}

func (l *Logger) SetDefault(prefix string) {
	l.Config.Loglevel = INFO
	l.Config.Flag = log.Ldate | log.Ltime
	l.Config.Prefix = prefix
}

func (l *Logger) Init() error {
	var out io.Writer
	if l.Config.Logpath == "" {
		out = os.Stdout
	} else {
		f, err := os.OpenFile(l.Config.Logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		out = io.MultiWriter(os.Stdout, f)
	}

	l.logger = log.New(out, l.Config.Prefix, l.Config.Flag)
	return nil
}

func (l *Logger) Debug(msg string, v ...any) {
	if l.Config.Loglevel < DEBUG {
		return
	}

	l.logger.Printf(msg, v...)
}

func (l *Logger) Info(msg string, v ...any) {
	if l.Config.Loglevel < INFO {
		return
	}

	l.logger.Printf(msg, v...)
}

func (l *Logger) Error(theerr error, msg string, v ...any) {
	if l.Config.Loglevel < ERROR {
		return
	}

	l.logger.Printf(msg, v...)
	l.logger.Printf("Error %v", theerr)
}

func Debug(msg string, v ...any) {
	gl.Debug(msg, v...)
}

func SetLevel(level LogLevel) {
	gl.Config.Loglevel = level
}
