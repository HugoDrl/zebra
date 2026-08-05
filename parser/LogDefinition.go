package parser

import (
	"strings"
	"time"
)

type Level string

const (
	Debug   Level = "debug"
	Info    Level = "info"
	Warning Level = "warning"
	Error   Level = "error"
	Fatal   Level = "fatal"
)

func toLevel(input string) (Level, bool) {
	switch inputLevel := Level(strings.ToLower(input)); inputLevel {
	case Debug:
		return Debug, true
	case Info:
		return Info, true
	case Warning:
		return Warning, true
	case Error:
		return Error, true
	case Fatal:
		return Fatal, true
	default:
		return "", false
	}
}

type Log struct {
	Time     time.Time
	Level    Level
	Duration time.Duration
	Message  string
	Service  string
	Extra    map[string]string
}

type ParseSettings struct {
	Files     []string
	StartDate time.Time
	EndDate   time.Time
	Level     Level
	Service   string
}
