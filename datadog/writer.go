package datadog

import (
	"bytes"
	"context"
	"github.com/cultureamp/glamplify/env"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cultureamp/glamplify/log"
)

const (
	dataDogAPIKeyHeader = "DD-API-KEY"
	contentTypeHeader   = "Content-Type"
	applicationJSONType = "application/json"
)

// DDWriter interface represents a log writer for data dog
type DDWriter interface {
	WriteFields(sev string, system log.Fields, fields ...log.Fields) string
	IsEnabled(sev string) bool
	WaitAll()
}

// DDFieldWriter sends logging output to Data Dog
type DDFieldWriter struct {
	// PUBLIC

	// APIKey for Data Dog API key.
	//
	// https://app.datadoghq.com/account/settings#api
	APIKey string

	// Endpoint URL: https://http-intake.logs.datadoghq.com/v1/input  (default)
	Endpoint string

	// Timeout on HTTP requests
	Timeout time.Duration

	// OmitEmpty will remove empty fields before sending
	OmitEmpty bool

	// Level we are logging, DEBUG, INFO, etc.
	Level string

	// PRIVATE

	// Allows us to WaitAll if clients want to make sure all pending writes have been sent
	waitGroup *sync.WaitGroup
	// converts to and from string <-> int
	leveller *log.Leveller
}

// NewDataDogWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying  writer.
func NewDataDogWriter(configure ...func(*DDFieldWriter)) DDWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	writer := &DDFieldWriter{
		APIKey:    os.Getenv(env.DatadogAPIKey),
		Endpoint:  env.GetString(env.DatadogLogEndpoint, "https://http-intake.logs.datadoghq.com/v1/input"),
		Timeout:   time.Second * time.Duration(env.GetInt(env.DatadogTimeout, 5)),
		OmitEmpty: env.GetBool(env.LogOmitEmpty, false),
		Level:     env.GetString(env.LogLevel, log.DebugSev),
		waitGroup: &sync.WaitGroup{},
		leveller:  log.NewLevelMap(),
	}

	for _, config := range configure {
		config(writer)
	}

	return writer
}

//WriteFields - writes fields to the Data Dog log endpoint
func (writer *DDFieldWriter) WriteFields(sev string, system log.Fields, fields ...log.Fields) string {
	merged := log.Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[log.Properties] = properties
	}

	json := system.ToSnakeCase().ToJSON(writer.OmitEmpty)
	if writer.IsEnabled(sev) {
		writer.waitGroup.Add(1)
		go post(writer, json)
	}
	return json
}

// IsEnabled returns true if the sev is enabled, false otherwise
func (writer DDFieldWriter) IsEnabled(sev string) bool {
	return writer.leveller.ShouldLogSeverity(writer.Level, sev)
}

// WaitAll waits until all the writers have finished
func (writer *DDFieldWriter) WaitAll() {
	writer.waitGroup.Wait()
}

func post(writer *DDFieldWriter, jsonStr string) {
	defer writer.waitGroup.Done()

	jsonBytes := []byte(jsonStr)

	// golint-noctx requires ctx for http requests
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", writer.Endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		panic(err)
	}

	req.Header.Set(contentTypeHeader, applicationJSONType)
	req.Header.Set(dataDogAPIKeyHeader, writer.APIKey)

	var client = &http.Client{
		Timeout: writer.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
