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

type NRWriter interface {
	WriteFields(sev int, system log.Fields, fields ...log.Fields) string
	WaitAll()
}

// NRFieldWriter sends logging output to NR as per https://docs.newrelic.com/docs/logs/new-relic-logs/log-api/introduction-log-api
type NRFieldWriter struct {
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

	// Omitempty will remove empty fields before sending
	Omitempty bool

	// Allows us to WaitAll if clients want to make sure all pending writes have been sent
	waitGroup sync.WaitGroup
}

// newWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying standard writer.
// Useful for CLI apps that want to direct logging to a file or stderr
// eg. SetOutput
func newWriter(configure ...func(*NRFieldWriter)) NRWriter {
	// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	writer := &NRFieldWriter{
		License:  os.Getenv("NEW_RELIC_LICENSE_KEY"),
		Endpoint: helper.GetEnvString("NEW_RELIC_LOG_ENDPOINT", "https://log-api.newrelic.com/log/v1"),
		Timeout:  time.Second * time.Duration(helper.GetEnvInt("NEW_RELIC_TIMEOUT", 5)),
		Omitempty: helper.GetEnvBool(log.OmitEmpty, false),
		waitGroup: sync.WaitGroup{},
	}

	for _, config := range configure {
		config(writer)
	}

	return writer
}

func (writer *NRFieldWriter) WriteFields(sev int, system log.Fields, fields ...log.Fields) string {
	merged := log.Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[log.Properties] = properties
	}

	json := system.ToSnakeCase().ToJson(writer.Omitempty)

	writer.waitGroup.Add(1)
	go post(writer, json)
	return json
}

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
