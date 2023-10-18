package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

func New(level string, stage string) *Logger {
	var globalLevel zerolog.Level

	switch strings.ToLower(level) {
	case "debug":
		globalLevel = zerolog.DebugLevel
	case "info":
		globalLevel = zerolog.InfoLevel
	case "warn":
		globalLevel = zerolog.WarnLevel
	case "error":
		globalLevel = zerolog.ErrorLevel
	case "fatal":
		globalLevel = zerolog.FatalLevel
	default:
		globalLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(globalLevel)

	var logger zerolog.Logger
	if stage == "dev" {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		output.FormatLevel = func(i any) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}

		logger = zerolog.New(output).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.CallerMarshalFunc = lshortfile

		logger = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	zerolog.DefaultContextLogger = &logger

	return &Logger{&logger}
}

// lshortfile - final file name element and line number: d.go:23
//
// From https://github.com/rs/zerolog#add-file-and-line-number-to-log
func lshortfile(pc uintptr, file string, line int) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file + ":" + strconv.Itoa(line)
}
