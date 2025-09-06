package infoblox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type Client struct {
	baseURL  string
	username string
	password string
	version  string
	client   *http.Client
	log      *zap.Logger
}

type ProviderArgs struct {
	Type     string
	Host     string
	Port     string
	Version  string
	Username string
	Password string
	log      *zap.Logger
}

// Create
//
// param: context
// return:
func Create(context ProviderArgs) *Client {
	baseUrl := "https://infoblox01v.dtlabs.de"
	username := context.Username
	password := context.Password
	version := context.Version
	log := context.log

	tr := &http.Transport{
		/* #nosec */
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return &Client{baseUrl, username, password, version, client, log}
}

// get
//
// param: reqURL
// param: response
func (c *Client) get(reqURL string, request interface{}, response interface{}) error {
	slogger := c.log.Sugar()
	reqURL, _ = url.JoinPath(c.baseURL, "wapi", c.version, reqURL)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	httpRequest, err := http.NewRequest(http.MethodGet, reqURL, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.SetBasicAuth(c.username, c.password)
	logRequest(httpRequest, slogger)
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	logResponse(httpResponse, slogger)
	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		if httpResponse.StatusCode == 204 || response == nil {
			return err
		}
		err := json.NewDecoder(httpResponse.Body).Decode(response)
		if err != nil {
			slogger.Errorf("response deserialization error - %v", err)
			return err
		}
		return nil
	}
	return fmt.Errorf("infoblox responded with an error")
}
func (c *Client) post(reqURL string, request interface{}, response interface{}) error {
	slogger := c.log.Sugar()
	reqURL, _ = url.JoinPath(c.baseURL, "wapi", c.version, reqURL)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	httpRequest, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.SetBasicAuth(c.username, c.password)
	logRequest(httpRequest, slogger)
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		slogger.Errorf("unable to make HTTP request - %v", err)
		return err
	}
	defer httpResponse.Body.Close()
	logResponse(httpResponse, slogger)
	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		if httpResponse.StatusCode == 204 || response == nil {
			return err
		}
		err := json.NewDecoder(httpResponse.Body).Decode(response)
		if err != nil {
			slogger.Errorf("response deserialization error - %v", err)
			return err
		}
		return nil
	}
	return fmt.Errorf("infoblox responded with an error")
}

func (c *Client) delete(reqURL string) error {
	slogger := c.log.Sugar()
	reqURL, _ = url.JoinPath(c.baseURL, "wapi", c.version, reqURL)
	httpRequest, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}
	httpRequest.SetBasicAuth(c.username, c.password)
	logRequest(httpRequest, slogger)
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	logResponse(httpResponse, slogger)
	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("infoblox responded with an error")
}

func logRequest(httpRequest *http.Request, slogger *zap.SugaredLogger) {
	slogger.Debugf("HTTP request - method: %v url: %v", httpRequest.Method, httpRequest.URL)

	// var bodyBytes []byte
	if httpRequest.Body != nil {
		bodyBytes, _ := io.ReadAll(httpRequest.Body)
		// Reset the body so it can be read again
		httpRequest.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		slogger.Debugf("HTTP request - payload: %s", bodyBytes)
	}
}

func logResponse(httpResponse *http.Response, slogger *zap.SugaredLogger) {
	slogger.Debugf("HTTP response - code: %v", httpResponse.StatusCode)

	// var bodyBytes []byte
	if httpResponse.Body != nil {
		bodyBytes, _ := io.ReadAll(httpResponse.Body)
		// Reset the body so it can be read again
		httpResponse.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		slogger.Debugf("HTTP response - payload: %s", bodyBytes)
	}
}
