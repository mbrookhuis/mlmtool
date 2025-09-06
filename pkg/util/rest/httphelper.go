// Package rest - rest api helper
package rest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	// INFOBLOX infoblox URL
	INFOBLOX        = "infoblox01v"
	statusCodeMin   = 300
	statusCodeMax   = 500
	intervalSeconds = 2
)

// HTTPHelperStruct - HTTP call
type HTTPHelperStruct struct {
	Body       []byte
	StatusCode int
	Cookies    []*http.Cookie
}

// IRestHelper - Rest Helper Interface
type IRestHelper interface {
	HTTPCaller(skipStatusRetry bool, requestBody []byte, method string, url string,
		headers map[string]string) (*HTTPHelperStruct, error)
}

// Helper - Rest API Helper Struct
type Helper struct {
	// retryCount
	retryCount int
	logger     *zap.Logger
}

// NewRestHelper - Create New Rest API Helper
// @return IRestHelper
func NewRestHelper(retryCount int, logger *zap.Logger) IRestHelper {
	return &Helper{
		retryCount: retryCount,
		logger:     logger,
	}
}

// HTTPCaller Caller
//
//	@receiver r *Helper
//	@param skipStatusRetry
//	@param requestBody []byte
//	@param method string
//	@param url string
//	@param headers map[string]string
//	@return *HTTPHelperStruct
//	@return error
func (r *Helper) HTTPCaller(skipStatusRetry bool, requestBody []byte, method string, url string,
	headers map[string]string) (*HTTPHelperStruct, error) {
	// Create Request
	r.logger.Info("HTTPCaller serving request")
	reqBody := bytes.NewBuffer(requestBody)

	client := &http.Client{}

	// This block is added as Infoblox certificates are not updated
	// Adding the check to skip the certificate check when accessing Infoblox
	if strings.Contains(url, INFOBLOX) {
		client = &http.Client{
			Transport: &http.Transport{
				/* #nosec */
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	//Send Request
	var res *http.Response

	// Retry
	retry := 1
	r.logger.Info(fmt.Sprintf("Making rest call to %s with retries %d", url, r.retryCount))
	for {
		res, err = client.Do(req)
		if err != nil {
			r.logger.Debug(fmt.Sprintf("error from server %s", err))
		} else if res.StatusCode >= statusCodeMin && res.StatusCode <= statusCodeMax {
			r.logger.Debug(fmt.Sprintf("error response from server %+v", res))
			if skipStatusRetry {
				break
			}
		} else {
			break
		}

		r.logger.Debug(fmt.Sprintf("Retry %d out of %d, for URL: %s", retry, r.retryCount, url))
		time.Sleep(time.Duration(retry*intervalSeconds) * time.Second)
		retry++
		if retry >= r.retryCount {
			if err != nil {
				return nil, err
			}
			break
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		r.logger.Error(fmt.Sprintf("error while reading server response %s", err))
		return nil, err
	}
	defer res.Body.Close() // nolint:all
	output := &HTTPHelperStruct{
		Body:       body,
		StatusCode: res.StatusCode,
		Cookies:    res.Cookies(),
	}
	r.logger.Debug(fmt.Sprintf("Status code from server %v", output.StatusCode))
	r.logger.Info("HTTPCaller served request with success")

	return output, nil
}

func (h *HTTPHelperStruct) String() string {
	return fmt.Sprintf("Body: %s, Status Code: %d, Cookies: %v", h.Body, h.StatusCode, h.Cookies)
}

// HTTPHelper - rest api helper
//
// param: log
// param: retrycount
// param: requsetBody
// param: method
// param: url
// param: insecure
// param: header
// return:
func HTTPHelper(log *zap.Logger, retrycount int, requsetBody []byte, method string, url string, insecure bool, header ...map[string]string) (*HTTPHelperStruct, error) {
	reqBody := bytes.NewBuffer(requsetBody)

	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	for k := range header {
		for i, j := range header[k] {
			req.Header.Add(i, j)
		}
	}
	var res *http.Response
	retry := 0
	log.Info("Connecting to: ", zap.Any("url", url))
	for ok := true; ok; ok = (err != nil || retry == retrycount) {
		res, err = client.Do(req)
		if err != nil {
			log.Info("Retrying connect to: ", zap.Any("retry count", retry))
			time.Sleep(time.Duration(retry*2) * time.Second)
			if retry == retrycount {
				log.Error("Failed connect to: ", zap.Any("url", url))
				return nil, err
			}
		}
		retry++
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	httpHelperStruct := HTTPHelperStruct{
		Body:       body,
		StatusCode: res.StatusCode,
		Cookies:    res.Cookies(),
	}
	return &httpHelperStruct, nil
}
