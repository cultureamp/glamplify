package log

import (
	"fmt"
	"github.com/cultureamp/glamplify/env"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gookit/color"
)

// WriterConfig for setting initial values for Logger
type WriterConfig struct {
	Output     io.Writer
	OmitEmpty  bool
	UseColours bool
	Level      string
}

// FieldWriter wraps the standard library writer and add structured types as quoted key value pairs
type FieldWriter struct {
	mutex    *sync.Mutex
	levelMap *Leveller

	output    io.Writer
	omitempty bool
	useColors bool
	level     int
}

// Writer defines an interface for writing log messages
type Writer interface {
	WriteFields(sev string, system Fields, fields ...Fields) string
	IsEnabled(sev string) bool
}

// NewWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying standard writer.
// Useful for CLI apps that want to direct logging to a file or stderr
// eg. SetOutput
func NewWriter(configure ...func(*WriterConfig)) *FieldWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	writer := &FieldWriter{}
	conf := WriterConfig{
		Output:     os.Stdout,
		OmitEmpty:  env.GetBool(env.LogOmitEmpty, false),
		UseColours: env.GetBool(env.LogUseColours, false),
		Level:      env.GetString(env.LogLevel, DebugSev),
	}
	for _, config := range configure {
		config(&conf)
	}

	writer.mutex = &sync.Mutex{}
	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	writer.levelMap = NewLevelMap()
	writer.output = conf.Output
	writer.omitempty = conf.OmitEmpty
	writer.useColors = conf.UseColours
	writer.level = writer.levelMap.StringToLevel(conf.Level)

	return writer
}

// WriteFields returns a json string for the given severity and system and user Fields
func (writer *FieldWriter) WriteFields(sev string, system Fields, fields ...Fields) string {
	merged := Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[Properties] = properties
	}
	json := system.ToSnakeCase().ToJSON(writer.omitempty)

	if writer.IsEnabled(sev) {
		writer.write(sev, json)
	}
	return json
}

// IsEnabled returns true if the sev is enabled, false otherwise
func (writer FieldWriter) IsEnabled(sev string) bool {
	level := writer.levelMap.StringToLevel(sev)

	return writer.levelMap.ShouldLogLevel(writer.level, level)
}

func (writer *FieldWriter) write(sev string, json string) {
	// This can return an error, but we just swallow it here as what can we or a client really do? Try and log it? :)
	json = writer.addNewLineIfMissing(json)

	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	if writer.useColors {
		// Helpful for humans, but SLOWS down the output writing, so don't recommend this for production
		// Also we purposely print with double NewLines (1 in the string and an extra one when printing)
		// to make it easy to separate different log lines...
		color.SetOutput(writer.output)
		level := writer.levelMap.StringToLevel(sev)
		switch level {
		case DebugLevel:
			color.Debug.Println(json)
		case InfoLevel:
			color.Info.Println(json)
		case WarnLevel:
			color.Warn.Println(json)
		case ErrorLevel:
			color.Error.Println(json)
		case FatalLevel:
			color.Danger.Println(json)
		case AuditLevel:
			color.Notice.Println(json)
		default:
			color.Print(json)
		}
	} else {
		// Note: Making this faster is a good thing (while we are a sync writer - async writer is a different story)
		// So we don't use the stdlib writer.Print(), but rather have our own optimized version
		// Which does less, but is 3-10x faster
		_, err := writer.output.Write([]byte(json))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (writer *FieldWriter) addNewLineIfMissing(str string) string {
	var b strings.Builder
	b.WriteString(str)
	l := len(str)
	if str[l-1] != '\n' {
		b.WriteString("\n")
	}

	return b.String()
}
