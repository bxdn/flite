package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RequestConfig struct {
	Url     string
	Headers map[string]string
	Query   map[string]string
}

type ConfigWithBody struct {
	RequestConfig
	Body []byte
}

type fullConfig struct {
	ConfigWithBody
	Method string
}

func FromJson[T any](res *http.Response, bodyBytes []byte, e error) (*http.Response, T, error) {
	ptr := new(T)
	if e != nil {
		return res, *ptr, e
	}
	decoder := json.NewDecoder(bytes.NewBuffer(bodyBytes))
	e = decoder.Decode(ptr)
	return res, *ptr, e
}

func ToJson(config RequestConfig, object any, req func(ConfigWithBody) (*http.Response, []byte, error)) (*http.Response, []byte, error) {
	jsonBytes, e := json.Marshal(object)
	if e != nil {
		return nil, nil, fmt.Errorf("Error Marshalling body to JSON: %w", e)
	}
	return req(ConfigWithBody{RequestConfig: config, Body: jsonBytes})
}

func Get(config RequestConfig) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: ConfigWithBody{RequestConfig: config}, Method: "GET"})
}

func Delete(config RequestConfig) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: ConfigWithBody{RequestConfig: config}, Method: "DELETE"})
}

func Post(config ConfigWithBody) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: config, Method: "POST"})
}

func Put(config ConfigWithBody) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: config, Method: "PUT"})
}

func Patch(config ConfigWithBody) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: config, Method: "PATCH"})
}

func req(config fullConfig) (*http.Response, []byte, error) {

	u, err := url.Parse(config.Url)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing url: %w", err)
	}

	for k, v := range config.Query {
		u.Query().Set(k, v)
	}

	req, err := http.NewRequest(config.Method, u.String(), bytes.NewBuffer(config.Body))
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating request: %w", err)
	}

	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Error executing request: %w", err)
	}

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading response body: %w", err)
	}

	return resp, resBody, nil
}
