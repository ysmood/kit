package gokit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/kataras/iris/core/errors"
	"github.com/tidwall/gjson"
)

// HTTPClient ...
type HTTPClient struct {
	Client   *http.Client
	Response *http.Response
}

// Req send http request
func Req(params ...interface{}) (*HTTPClient, error) {
	method, url, cookie, header, reqBody, err := parseReqParams(params)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: cookie,
	}

	req, err := http.NewRequest(string(method), url, reqBody)
	if err != nil {
		return nil, err
	}

	for _, h := range header {
		req.Header.Add(h[0], h[1])
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		Client:   client,
		Response: res,
	}, nil
}

// Req reuse the cookie
func (req *HTTPClient) Req(params ...interface{}) (*HTTPClient, error) {
	method, url, _, header, reqBody, err := parseReqParams(params)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(string(method), url, reqBody)
	if err != nil {
		return nil, err
	}

	for _, h := range header {
		r.Header.Add(h[0], h[1])
	}

	res, err := req.Client.Do(r)
	if err != nil {
		return nil, err
	}
	req.Response = res

	return req, nil
}

// Bytes ...
func (req *HTTPClient) Bytes() []byte {
	body, err := ioutil.ReadAll(req.Response.Body)
	if err != nil {
		return []byte(err.Error())
	}

	err = req.Response.Body.Close()
	if err != nil {
		return []byte(err.Error())
	}

	return body
}

func (req *HTTPClient) String() string {
	return string(req.Bytes())
}

// JSON ...
func (req *HTTPClient) JSON(v interface{}) error {
	return json.Unmarshal(req.Bytes(), &v)
}

// GJSON parse body as json and provide searching for json strings
func (req *HTTPClient) GJSON() *gjson.Result {
	r := gjson.ParseBytes(req.Bytes())
	return &r
}

// Method ...
type Method string

// QueryParams ...
type QueryParams map[string]interface{}

// Header ...
type Header map[string]string

// FormParams ...
type FormParams map[string]interface{}

// JSONBody ...
type JSONBody interface{}

// StringBody ...
type StringBody string

// Infer the params from their type, the order doesn't matter
func parseReqParams(params []interface{}) (
	method Method,
	reqURL string,
	cookie *cookiejar.Jar,
	header [][]string,
	body io.Reader,
	err error,
) {
	for _, param := range params {
		switch v := param.(type) {
		case Method:
			method = v
		case string:
			reqURL = v

		case QueryParams:
			var query string
			query, err = qs.Marshal(v)
			reqURL += "?" + query

		case *cookiejar.Jar:
			cookie = v

		case Header:
			for key, val := range v {
				header = append(header, []string{key, val})
			}

		case FormParams:
			var query string
			query, err = qs.Marshal(v)
			header = append(header, []string{"Content-Type", "application/x-www-form-urlencoded; charset=utf-8"})
			body = strings.NewReader(query)

		case io.Reader:
			body = v

		case JSONBody:
			var b []byte
			b, err = json.Marshal(v)
			if err != nil {
				continue
			}
			header = append(header, []string{"Content-Type", "application/json; charset=utf-8"})
			body = bytes.NewReader(b)

		case StringBody:
			body = strings.NewReader(string(v))

		default:
			err = errors.New(fmt.Sprintf("params type not supported: %T", v))
			break
		}
	}

	if method == "" {
		method = http.MethodGet
	}

	if cookie == nil {
		cookie, err = cookiejar.New(nil)
	}

	return method, reqURL, cookie, header, body, err
}
