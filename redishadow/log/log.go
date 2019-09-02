package log

/*
Author: Luncert
*/

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func infoF(format string, a ...interface{}) {
	fmt.Printf("[INFO] %s", fmt.Sprintf(format, a...))
}

func errorF(format string, a ...interface{}) {
	fmt.Printf("[ERROR] %s", fmt.Sprintf(format, a...))
}

func fatalF(format string, a ...interface{}) {
	fmt.Printf("[FATAL] %s", fmt.Sprintf(format, a...))
	os.Exit(1)
}

// log level
const (
	debugLevel = iota
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
)

type logAppender interface {
	Write(data []byte) error
	Close() error
}

type logger struct {
	level     int
	formatter *logFormatter
	appender  logAppender
}

func (l *logger) log(level int, v interface{}) {
	if level >= l.level {
		if err := l.appender.Write(l.formatter.format(v)); err != nil {
			errorF("Failed to write log: %v", err)
		}
	}
}

var log logger

func init() {
	data, err := ioutil.ReadFile("log.yml")
	if err != nil {
		fatalF("Failed to read log.yml: %v", err)
	}

	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fatalF("Failed ot unmarshal configuration: %v", err)
		return
	}

	log = logger{}

	if level, ok := config["level"]; !ok {
		fatalF("Config missing: level")
	} else {
		switch strings.ToLower(level.(string)) {
		case "debug":
			log.level = debugLevel
		case "info":
			log.level = infoLevel
		case "warn":
			log.level = warnLevel
		case "error":
			log.level = errorLevel
		case "fatal":
			log.level = fatalLevel
		default:
			fatalF("Unknown log level: %s", level)
		}
	}

	if format, ok := config["format"]; !ok {
		fatalF("Config missing: format")
	} else {
		log.formatter = newFormatter(format.(string))
	}

	appenderType, ok := config["appender"]
	if !ok {
		infoF("No appender defined, using default one: stdout appender")
	}
	switch strings.ToLower(appenderType.(string)) {
	case "tcp":
		log.appender = nil
	case "udp":
		log.appender = nil
	case "file":
		log.appender = nil
	case "stdout":
		fallthrough
	default:
		log.appender = &stdoutAppender{}
	}
}

func Debug(v interface{}) {
	log.log(debugLevel, v)
}

func Info(v interface{}) {
	log.log(infoLevel, v)
}

func Warn(v interface{}) {
	log.log(warnLevel, v)
}

func Error(v interface{}) {
	log.log(errorLevel, v)
}

func Fatal(v interface{}) {
	log.log(fatalLevel, v)
}
