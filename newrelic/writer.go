package newrelic

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/log"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

//  NRWriter defines an interface for writing log messages to newrelic
type NRWriter interface {
	WriteFields(sev string, system log.Fields, fields ...log.Fields) string
	IsEnabled(sev string) bool
	WaitAll()
}

// NRFieldWriter sends logging output to NR as per https://docs.newrelic.com/docs/logs/new-relic-logs/log-api/introduction-log-api
type NRFieldWriter struct {
	// PUBLIC

	// license is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key
	License string

	// URL
	// US: https://log-api.newrelic.com/log/v1  (default)
	// EU: https://log-api.eu.newrelic.com/log/v1
	Endpoint string

	// Timeout on HTTP requests
	Timeout time.Duration

	// OmitEmpty will remove empty fields before sending
	OmitEmpty bool

	// Level we are logging, DEBUG, INFO, etc.
	Level string

	// PRIVATE

	// Allows us to WaitAll if clients want to make sure all pending writes have been sent
	waitGroup sync.WaitGroup
	// converts to and from string <-> int
	leveller *log.Leveller
}

// NewNRWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying standard writer.
// Useful for CLI apps that want to direct logging to a file or stderr
// eg. SetOutput
func NewNRWriter(configure ...func(*NRFieldWriter)) NRWriter {
	// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	writer := &NRFieldWriter{
		License:   os.Getenv("NEW_RELIC_LICENSE_KEY"),
		Endpoint:  helper.GetEnvString("NEW_RELIC_LOG_ENDPOINT", "https://log-api.newrelic.com/log/v1"),
		Timeout:   time.Second * time.Duration(helper.GetEnvInt("NEW_RELIC_TIMEOUT", 5)),
		OmitEmpty: helper.GetEnvBool(log.OmitEmpty, false),
		Level:     helper.GetEnvString(log.Level, log.DebugSev),
		waitGroup: sync.WaitGroup{},
		leveller:  log.NewLevelMap(),
	}

	for _, config := range configure {
		config(writer)
	}

	return writer
}

// WriteFields returns json representing the system and user Fields if the severity is set
func (writer *NRFieldWriter) WriteFields(sev string, system log.Fields, fields ...log.Fields) string {
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

// IsEnabled returns true is the severity is set, false otherwise
func (writer NRFieldWriter) IsEnabled(sev string) bool {
	if writer.leveller.ShouldLogSeverity(writer.Level, sev) {
		return true
	}
	return false
}

// WaitAll waits for all the writers to finish before returning
func (writer *NRFieldWriter) WaitAll() {
	writer.waitGroup.Wait()
}

func post(writer *NRFieldWriter, jsonStr string) error {
	defer writer.waitGroup.Done()

	// https://docs.newrelic.com/docs/logs/new-relic-logs/log-api/introduction-log-api
	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", writer.Endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-license-Key", writer.License)

	var client = &http.Client{
		Timeout: writer.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	str := string(body)
	return errors.New(fmt.Sprintf("bad server response: %d. body: %v", resp.StatusCode, str))
}
