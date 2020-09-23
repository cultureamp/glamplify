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

type DDWriter interface {
	WriteFields(sev int, system log.Fields, fields ...log.Fields)
	WaitAll()
}

// DDFieldWriter sends logging output to Data Dog
type DDFieldWriter struct {
	// ApiKey for Data Dog API key.
	//
	// https://app.datadoghq.com/account/settings#api
	ApiKey string

	// Endpoint URL: https://http-intake.logs.datadoghq.com/v1/input  (default)
	Endpoint string

	// Timeout on HTTP requests
	Timeout time.Duration

	// Omitempty will remove empty fields before sending
	Omitempty bool

	// WaitGroup can optionally be to a valid wait group, and the writer will signal when it sends and completes
	// so clients can
	WaitGroup *sync.WaitGroup
}

// NewDataDogWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying  writer.
func NewDataDogWriter(configure ...func(*DDFieldWriter)) DDWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	writer := &DDFieldWriter{
		ApiKey:    os.Getenv("DD_CLIENT_API_KEY"),
		Endpoint:  helper.GetEnvString("DATA_DOG_LOG_ENDPOINT", "https://http-intake.logs.datadoghq.com/v1/input"),
		Timeout:   time.Second * time.Duration(helper.GetEnvInt("DATA_DOG_TIMEOUT", 5)),
		Omitempty: helper.GetEnvBool(log.OmitEmpty, false),
		WaitGroup: nil,
	}

	for _, config := range configure {
		config(writer)
	}

	return writer
}

//WriteFields - writes fields to the Data Dog log endpoint
func (writer *DDFieldWriter) WriteFields(sev int, system log.Fields, fields ...log.Fields) {

	merged := log.Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[log.Properties] = properties
	}

	json := system.ToSnakeCase().ToJson(writer.Omitempty)

	if writer.WaitGroup != nil {
		writer.WaitGroup.Add(1)
	}
	go post(writer, json)
}

func (writer *DDFieldWriter) WaitAll() {
	if writer.WaitGroup != nil {
		writer.WaitGroup.Wait()
	} else {
		time.Sleep(writer.Timeout)
	}
}

func post(writer *DDFieldWriter, jsonStr string) error {
	if writer.WaitGroup != nil {
		defer writer.WaitGroup.Done()
	}

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
