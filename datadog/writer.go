package datadog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/log"
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

	// ApiKey for Data Dog API key.
	//
	// https://app.datadoghq.com/account/settings#api
	ApiKey string

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
	waitGroup sync.WaitGroup
	// converts to and from string <-> int
	leveller *log.Leveller

}

// NewDataDogWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying  writer.
func NewDataDogWriter(configure ...func(*DDFieldWriter)) DDWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	writer := &DDFieldWriter{
		ApiKey:    os.Getenv(DDApiKey),
		Endpoint:  helper.GetEnvString(DDLogEndpoint, "https://http-intake.logs.datadoghq.com/v1/input"),
		Timeout:   time.Second * time.Duration(helper.GetEnvInt(DDTimeout, 5)),
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

//WriteFields - writes fields to the Data Dog log endpoint
func (writer *DDFieldWriter) WriteFields(sev string, system log.Fields, fields ...log.Fields) string {

	merged := log.Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[log.Properties] = properties
	}

	json := system.ToSnakeCase().ToJson(writer.OmitEmpty)
	if writer.IsEnabled(sev) {
		writer.waitGroup.Add(1)
		go post(writer, json)
	}
	return json
}

// IsEnabled returns true if the sev is enabled, false otherwise
func (writer DDFieldWriter) IsEnabled(sev string) bool {
	if writer.leveller.ShouldLogSeverity(writer.Level, sev) {
		return true
	}
	return false
}

// WaitAll waits until all the writers have finished
func (writer *DDFieldWriter) WaitAll() {
	writer.waitGroup.Wait()
}

func post(writer *DDFieldWriter, jsonStr string) error {
	defer writer.waitGroup.Done()

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", writer.Endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", writer.ApiKey)

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
	return fmt.Errorf("bad server response: %d. body: %v", resp.StatusCode, str)
}
