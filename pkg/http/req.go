package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alessio/shellescape"
	"github.com/derekstavis/go-qs"
	"github.com/tidwall/gjson"
	"github.com/ysmood/kit/pkg/utils"
)

// ReqContext the request context
type ReqContext struct {
	context  context.Context
	client   *http.Client
	request  *http.Request
	response *http.Response

	method     string
	url        string
	host       string
	header     http.Header
	jsonBody   interface{}
	stringBody string
	body       io.Reader
	resBytes   []byte

	timeout       time.Duration
	timeoutCancel func()
}

// Req creates http request instance
func Req(url string) *ReqContext {
	return &ReqContext{
		header: http.Header{},
		url:    url,
	}
}

// Context sets the context of the request
func (ctx *ReqContext) Context(c context.Context) *ReqContext {
	ctx.context = c
	return ctx
}

// Timeout sets the timeout of the request, it will inherit the Context
func (ctx *ReqContext) Timeout(d time.Duration) *ReqContext {
	ctx.timeout = d
	return ctx
}

// Method sets request method
func (ctx *ReqContext) Method(m string) *ReqContext {
	ctx.method = m
	return ctx
}

// URL sets the url for request
func (ctx *ReqContext) URL(url string) *ReqContext {
	ctx.url = url
	return ctx
}

// Post sets the request method to POST
func (ctx *ReqContext) Post() *ReqContext {
	return ctx.Method(http.MethodPost)
}

// Put sets the request method to PUT
func (ctx *ReqContext) Put() *ReqContext {
	return ctx.Method(http.MethodPut)
}

// Delete sets the request method to DELETE
func (ctx *ReqContext) Delete() *ReqContext {
	return ctx.Method(http.MethodDelete)
}

// Query sets the query string of the request, example Query(k, v, k, v ...)
func (ctx *ReqContext) Query(params ...interface{}) *ReqContext {
	query, _ := qs.Marshal(paramsToForm(params))
	ctx.url += "?" + query
	return ctx
}

// Host sets the host request header
func (ctx *ReqContext) Host(host string) *ReqContext {
	ctx.host = host
	return ctx
}

// Header appends the request header, example Header(k, v, k, v ...)
func (ctx *ReqContext) Header(params ...string) *ReqContext {
	for i := 0; i < len(params); i += 2 {
		k := params[i]
		v := params[i+1]
		if _, has := ctx.header[k]; has {
			ctx.header[k] = append(ctx.header[k], v)
		} else {
			ctx.header[k] = []string{v}
		}
	}

	return ctx
}

// Headers sets the request header
func (ctx *ReqContext) Headers(header http.Header) *ReqContext {
	ctx.header = header
	return ctx
}

// Client sets http client
func (ctx *ReqContext) Client(c *http.Client) *ReqContext {
	ctx.client = c
	return ctx
}

// Form sets the request body as form, example Form(k, v, k, v)
func (ctx *ReqContext) Form(params ...interface{}) *ReqContext {
	query, _ := qs.Marshal(paramsToForm(params))
	ctx.header["Content-Type"] = []string{"application/x-www-form-urlencoded; charset=utf-8"}
	ctx.body = strings.NewReader(query)
	return ctx
}

// Body sets the request body
func (ctx *ReqContext) Body(b io.Reader) *ReqContext {
	ctx.body = b
	return ctx
}

// JSONBody sets request body as json
func (ctx *ReqContext) JSONBody(data interface{}) *ReqContext {
	ctx.header["Content-Type"] = []string{"application/json; charset=utf-8"}
	ctx.jsonBody = data

	return ctx
}

// StringBody sets request body as string
func (ctx *ReqContext) StringBody(s string) *ReqContext {
	ctx.stringBody = s
	return ctx
}

func (ctx *ReqContext) getBody() (io.Reader, error) {
	if ctx.stringBody != "" {
		return strings.NewReader(ctx.stringBody), nil
	}

	if ctx.jsonBody != nil {
		body, err := json.Marshal(ctx.jsonBody)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(body), nil
	}

	return ctx.body, nil
}

// Do the request
func (ctx *ReqContext) Do() error {
	req, err := ctx.Request()
	if err != nil {
		return err
	}

	res, err := ctx.client.Do(req)
	if err != nil {
		return err
	}
	if ctx.timeout != 0 {
		ctx.timeoutCancel()
	}
	ctx.response = res

	return nil
}

// MustDo send request, panic if request fails
func (ctx *ReqContext) MustDo() {
	utils.E(ctx.Do())
}

// Request gets native request struct, useful for debugging
func (ctx *ReqContext) Request() (*http.Request, error) {
	if ctx.request != nil {
		return ctx.request, nil
	}

	if ctx.context == nil {
		ctx.context = context.Background()
	}

	if ctx.timeout != 0 {
		c, cancel := context.WithTimeout(ctx.context, ctx.timeout)
		ctx.Context(c)
		ctx.timeoutCancel = cancel
	}

	if ctx.client == nil {
		cookie, _ := cookiejar.New(nil)
		ctx.client = &http.Client{
			Jar: cookie,
		}
	}

	body, err := ctx.getBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx.context, ctx.method, ctx.url, body)
	if err != nil {
		return nil, err
	}

	req.Header = ctx.header
	req.Host = ctx.host

	ctx.request = req

	return ctx.request, nil
}

// Response sends request, get response
func (ctx *ReqContext) Response() (*http.Response, error) {
	if ctx.response != nil {
		return ctx.response, nil
	}

	err := ctx.Do()
	if err != nil {
		return nil, err
	}
	return ctx.response, nil
}

// MustResponse panic version of Response
func (ctx *ReqContext) MustResponse() *http.Response {
	return utils.E(ctx.Response())[0].(*http.Response)
}

// Bytes sends request, read response body as bytes
func (ctx *ReqContext) Bytes() ([]byte, error) {
	res, err := ctx.Response()
	if err != nil {
		return nil, err
	}

	if ctx.resBytes == nil {
		ctx.resBytes, err = readBody(res.Body)
	}
	return ctx.resBytes, err
}

// MustBytes panic version of Bytes()
func (ctx *ReqContext) MustBytes() []byte {
	return utils.E(ctx.Bytes())[0].([]byte)
}

func readBody(b io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = b.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

// String sends request, read response as string
func (ctx *ReqContext) String() (string, error) {
	s, err := ctx.Bytes()
	return string(s), err
}

// MustString panic version of String()
func (ctx *ReqContext) MustString() string {
	return string(ctx.MustBytes())
}

// JSON sends request, get response and parse body as json and provide searching for json strings
func (ctx *ReqContext) JSON() (utils.JSONResult, error) {
	b, err := ctx.Bytes()
	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(b)
	return &r, nil
}

// MustJSON panic version of JSON()
func (ctx *ReqContext) MustJSON() utils.JSONResult {
	return utils.E(ctx.JSON())[0].(*gjson.Result)
}

func paramsToForm(params []interface{}) map[string]interface{} {
	f := map[string]interface{}{}

	for i := 0; i < len(params); i += 2 {
		f[params[i].(string)] = params[i+1]
	}
	return f
}

// MustCurl generates request and response details in curl style.
// Useful when reproduce request on other systems with minimum dependencies.
// For now gzip is not handled.
func (ctx *ReqContext) MustCurl() string {
	// get body
	body, err := ctx.getBody()
	utils.E(err)
	if body != nil {
		bodyData, err := ioutil.ReadAll(body)
		utils.E(err)
		ctx.stringBody = string(bodyData)
	}
	stringBody := ""
	if ctx.stringBody != "" {
		stringBody = " \\\n  -d " + shellescape.Quote(ctx.stringBody)
	}

	res, err := ctx.Response()
	utils.E(err)

	req, err := ctx.Request()
	utils.E(err)

	// request header
	reqHeaderStr := ""
	for _, h := range headerToArr(req.Header) {
		reqHeaderStr += " \\\n  -H " + shellescape.Quote(h[0]+": "+h[1])
	}

	resStr := res.Proto + " " + res.Status + "\n"

	for _, h := range headerToArr(res.Header) {
		resStr += h[0] + ": " + h[1] + "\n"
	}

	resBytes := ctx.MustBytes()
	var obj interface{}
	err = json.Unmarshal(resBytes, &obj)
	if err == nil {
		resBytes, _ = json.MarshalIndent(obj, "", "  ")
	} else if !utf8.Valid(resBytes) {
		resBytes = []byte(base64.StdEncoding.EncodeToString(resBytes))
	}

	resStr += "\n" + string(resBytes)

	return utils.S(
		"curl -X {{.method}} {{.url}}{{.header}}{{.data}}\n\n{{.resStr}}",
		"method", shellescape.Quote(req.Method),
		"url", shellescape.Quote(req.URL.String()),
		"header", reqHeaderStr,
		"data", stringBody,
		"resStr", strings.Trim(resStr, "\n"),
	)
}

func headerToArr(header http.Header) [][]string {
	headers := [][]string{}
	for k, vs := range header {
		for _, v := range vs {
			headers = append(headers, []string{k, v})
		}
	}
	sort.Slice(headers, func(a, b int) bool { return headers[a][0] < headers[b][0] })
	return headers
}
