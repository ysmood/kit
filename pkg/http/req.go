package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/tidwall/gjson"
)

type ReqContext struct {
	Error error

	client   *http.Client
	request  *http.Request
	response *http.Response

	method string
	url    string
	header [][]string
	body   io.Reader
}

// Req send http request
func Req(url string) *ReqContext {
	return &ReqContext{
		url: url,
	}
}

func (ctx *ReqContext) Method(m string) *ReqContext {
	ctx.method = m
	return ctx
}

func (ctx *ReqContext) Req(url string) *ReqContext {
	ctx.url = url
	return ctx.Do()
}

func (ctx *ReqContext) Post() *ReqContext {
	return ctx.Method(http.MethodPost)
}

func (ctx *ReqContext) Put() *ReqContext {
	return ctx.Method(http.MethodPut)
}

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

func (ctx *ReqContext) Form(params ...interface{}) *ReqContext {
	query, _ := qs.Marshal(paramsToForm(params))
	ctx.header = append(ctx.header, []string{"Content-Type", "application/x-www-form-urlencoded; charset=utf-8"})
	ctx.body = strings.NewReader(query)
	return ctx
}

func (ctx *ReqContext) Body(b io.Reader) *ReqContext {
	ctx.body = b
	return ctx
}

func (ctx *ReqContext) JSONBody(data interface{}) *ReqContext {
	b, err := json.Marshal(data)
	if err != nil {
		ctx.Error = err
		return ctx
	}
	ctx.header = append(ctx.header, []string{"Content-Type", "application/json; charset=utf-8"})
	ctx.body = bytes.NewReader(b)

	return ctx
}

func (ctx *ReqContext) StringBody(s string) *ReqContext {
	ctx.body = strings.NewReader(string(s))
	return ctx
}

func (ctx *ReqContext) Do() *ReqContext {
	if ctx.client == nil {
		cookie, _ := cookiejar.New(nil)
		ctx.client = &http.Client{
			Jar: cookie,
		}
	}

	req, err := http.NewRequest(ctx.method, ctx.url, ctx.body)
	if err != nil {
		ctx.Error = err
		return ctx
	}

	ctx.request = req

	for _, h := range ctx.header {
		req.Header.Add(h[0], h[1])
	}

	res, err := ctx.client.Do(req)
	if err != nil {
		ctx.Error = err
		return ctx
	}
	ctx.response = res

	return ctx
}

// Request get request
func (ctx *ReqContext) Request() *http.Request {
	return ctx.request
}

// Response get response
func (ctx *ReqContext) Response() *http.Response {
	return ctx.Do().response
}

// Bytes get response body as bytes
func (ctx *ReqContext) Bytes() []byte {
	body, err := readBody(ctx.Response().Body)
	ctx.Error = err
	return body
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

// String get string response
func (ctx *ReqContext) String() string {
	return string(ctx.Bytes())
}

// JSON unmarshal json response to v
func (ctx *ReqContext) JSON(v interface{}) error {
	return json.Unmarshal(ctx.Bytes(), &v)
}

// GJSON parse body as json and provide searching for json strings
func (ctx *ReqContext) GJSON(path string) gjson.Result {
	r := gjson.ParseBytes(ctx.Bytes())
	return r.Get(path)
}

func paramsToForm(params []interface{}) map[string]interface{} {
	f := map[string]interface{}{}
	l := len(params) - 1
	for i := 0; i < l; i += 2 {
		f[params[i].(string)] = params[i+1]
	}
	return f
}
