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

// writerConfig for setting initial values for Monitor Writer
type writerConfig struct {
	// license is your Data Dog API key.
	//
	// https://app.datadoghq.com/account/settings#api
	apiKey string

	// URL: https://http-intake.logs.datadoghq.com/v1/input  (default)
	endpoint string

	// timeout
	timeout time.Duration

	omitempty bool
}

// FieldWriter sends logging output to Data Dog
type FieldWriter struct {
	mutex  sync.Mutex
	config writerConfig
}

// NewDataDogWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying  writer.
func NewDataDogWriter(configure ...func(*writerConfig)) *FieldWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	conf := writerConfig{
		apiKey:    os.Getenv("DD_CLIENT_API_KEY"),
		endpoint:  helper.GetEnvString("DATA_DOG_LOG_ENDPOINT", "https://http-intake.logs.datadoghq.com/v1/input"),
		timeout:   time.Second * time.Duration(helper.GetEnvInt("DATA_DOG_TIMEOUT", 5)),
		omitempty: helper.GetEnvBool(log.OmitEmpty, false),
	}

	for _, config := range configure {
		config(&conf)
	}

	writer := &FieldWriter{
		config: conf,
	}

	return writer
}

//WriteFields - writes fields to the Data Dog log endpoint
func (writer *FieldWriter) WriteFields(sev int, system log.Fields, fields ...log.Fields) {

	merged := log.Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[log.Properties] = properties
	}

	json := system.ToSnakeCase().ToJson(writer.config.omitempty)

	go post(writer.config, json)
}

func post(config writerConfig, jsonStr string) error {

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", config.endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", config.apiKey)

	var client = &http.Client{
		Timeout: config.timeout,
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
