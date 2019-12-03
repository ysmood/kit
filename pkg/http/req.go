package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/derekstavis/go-qs"
	"github.com/tidwall/gjson"
	"github.com/ysmood/kit/pkg/utils"
	"gopkg.in/alessio/shellescape.v1"
)

// ReqContext the request context
type ReqContext struct {
	client   *http.Client
	request  *http.Request
	response *http.Response

	method     string
	url        string
	header     [][]string
	jsonBody   interface{}
	stringBody string
	body       io.Reader
	resBytes   []byte
}

// JSONResult shortcut for gjson.Result
type JSONResult = gjson.Result

// Req send http request
func Req(url string) *ReqContext {
	return &ReqContext{
		url: url,
	}
}

// Method request method
func (ctx *ReqContext) Method(m string) *ReqContext {
	ctx.method = m
	return ctx
}

// URL the url path for request
func (ctx *ReqContext) URL(url string) *ReqContext {
	ctx.url = url
	return ctx
}

// Post use POST as the method
func (ctx *ReqContext) Post() *ReqContext {
	return ctx.Method(http.MethodPost)
}

// Put use PUT as the method
func (ctx *ReqContext) Put() *ReqContext {
	return ctx.Method(http.MethodPut)
}

// Delete use DELETE as the method
func (ctx *ReqContext) Delete() *ReqContext {
	return ctx.Method(http.MethodDelete)
}

// Query Query(k, v, k, v ...)
func (ctx *ReqContext) Query(params ...interface{}) *ReqContext {
	query, _ := qs.Marshal(paramsToForm(params))
	ctx.url += "?" + query
	return ctx
}

// Header Header(k, v, k, v ...)
func (ctx *ReqContext) Header(params ...string) *ReqContext {
	for i := 0; i < len(params); i += 2 {
		ctx.header = append(ctx.header, []string{params[i], params[i+1]})
	}

	return ctx
}

// Client set http client
func (ctx *ReqContext) Client(c *http.Client) *ReqContext {
	ctx.client = c
	return ctx
}

// Form the params is a key-value pair list, such as `Form(k, v, k, v)`
func (ctx *ReqContext) Form(params ...interface{}) *ReqContext {
	query, _ := qs.Marshal(paramsToForm(params))
	ctx.header = append(ctx.header, []string{"Content-Type", "application/x-www-form-urlencoded; charset=utf-8"})
	ctx.body = strings.NewReader(query)
	return ctx
}

// Body the request body to sent
func (ctx *ReqContext) Body(b io.Reader) *ReqContext {
	ctx.body = b
	return ctx
}

// JSONBody set request body as json
func (ctx *ReqContext) JSONBody(data interface{}) *ReqContext {
	ctx.header = append(ctx.header, []string{"Content-Type", "application/json; charset=utf-8"})
	ctx.jsonBody = data

	return ctx
}

// StringBody set request body as a string
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

// Do send the request
func (ctx *ReqContext) Do() error {
	if ctx.client == nil {
		cookie, _ := cookiejar.New(nil)
		ctx.client = &http.Client{
			Jar: cookie,
		}
	}

	body, err := ctx.getBody()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(ctx.method, ctx.url, body)
	if err != nil {
		return err
	}

	ctx.request = req

	for _, h := range ctx.header {
		req.Header.Add(h[0], h[1])
	}

	res, err := ctx.client.Do(req)
	if err != nil {
		return err
	}
	ctx.response = res

	return nil
}

// MustDo send request, panic if request fails
func (ctx *ReqContext) MustDo() {
	utils.E(ctx.Do())
}

// Request get native request struct, useful for debugging
func (ctx *ReqContext) Request() (*http.Request, error) {
	if ctx.request != nil {
		return ctx.request, nil
	}
	err := ctx.Do()
	if err != nil {
		return nil, err
	}
	return ctx.request, nil
}

// Response send request, get response
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

// Bytes send request, read response body as bytes
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

// String send request, read response as string
func (ctx *ReqContext) String() (string, error) {
	s, err := ctx.Bytes()
	return string(s), err
}

// MustString panic version of String()
func (ctx *ReqContext) MustString() string {
	return string(ctx.MustBytes())
}

// JSON send request, get response and parse body as json and provide searching for json strings
func (ctx *ReqContext) JSON() (*JSONResult, error) {
	b, err := ctx.Bytes()
	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(b)
	return &r, nil
}

// MustJSON panic version of JSON()
func (ctx *ReqContext) MustJSON() *JSONResult {
	return utils.E(ctx.JSON())[0].(*gjson.Result)
}

func paramsToForm(params []interface{}) map[string]interface{} {
	f := map[string]interface{}{}

	for i := 0; i < len(params); i += 2 {
		f[params[i].(string)] = params[i+1]
	}
	return f
}

// MustCurl generate request and response details in curl style.
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
