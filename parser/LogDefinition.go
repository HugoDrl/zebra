package parser

import (
	"maps"
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
	switch inputLevel := Level(input); inputLevel {
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
	extra    map[string]string
}

func (l Log) Equal(log Log) bool {
	return (
		l.Time.Equal(log.Time) &&
		l.Level == log.Level &&
		l.Duration == log.Duration &&
		l.Message == log.Message &&
		l.Service == log.Service &&
		maps.Equal(l.extra, log.extra))
}

type ParseSettings struct {
	Files []string
	StartDate time.Time
	EndDate time.Time
	Level Level
	Service string
}
