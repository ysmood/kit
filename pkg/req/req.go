package req

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

type Context struct {
	client   *http.Client
	request  *http.Request
	response *http.Response

	err    error
	method string
	url    string
	header [][]string
	body   io.Reader
}

// Req send http request
func Req(url string) *Context {
	return &Context{
		url: url,
	}
}

func (ctx *Context) Method(m string) *Context {
	ctx.method = m
	return ctx
}

func (ctx *Context) Req(url string) *Context {
	ctx.url = url
	return ctx.Do()
}

func (ctx *Context) Post() *Context {
	return ctx.Method(http.MethodPost)
}

func (ctx *Context) Put() *Context {
	return ctx.Method(http.MethodPut)
}

func (ctx *Context) Delete() *Context {
	return ctx.Method(http.MethodDelete)
}

// Query Query(k, v, k, v ...)
func (ctx *Context) Query(params ...interface{}) *Context {
	query, _ := qs.Marshal(paramsToForm(params))
	ctx.url += "?" + query
	return ctx
}

// Header Header(k, v, k, v ...)
func (ctx *Context) Header(params ...string) *Context {
	for i := 0; i < len(params); i += 2 {
		ctx.header = append(ctx.header, []string{params[i], params[i+1]})
	}

	return ctx
}

// Client set http client
func (ctx *Context) Client(c *http.Client) *Context {
	ctx.client = c
	return ctx
}

func (ctx *Context) Form(params ...interface{}) *Context {
	query, _ := qs.Marshal(paramsToForm(params))
	ctx.header = append(ctx.header, []string{"Content-Type", "application/x-www-form-urlencoded; charset=utf-8"})
	ctx.body = strings.NewReader(query)
	return ctx
}

func (ctx *Context) Body(b io.Reader) *Context {
	ctx.body = b
	return ctx
}

func (ctx *Context) JSONBody(data interface{}) *Context {
	b, err := json.Marshal(data)
	if err != nil {
		ctx.err = err
		return ctx
	}
	ctx.header = append(ctx.header, []string{"Content-Type", "application/json; charset=utf-8"})
	ctx.body = bytes.NewReader(b)

	return ctx
}

func (ctx *Context) StringBody(s string) *Context {
	ctx.body = strings.NewReader(string(s))
	return ctx
}

func (ctx *Context) Do() *Context {
	if ctx.client == nil {
		cookie, _ := cookiejar.New(nil)
		ctx.client = &http.Client{
			Jar: cookie,
		}
	}

	req, err := http.NewRequest(ctx.method, ctx.url, ctx.body)
	if err != nil {
		ctx.err = err
		return ctx
	}

	ctx.request = req

	for _, h := range ctx.header {
		req.Header.Add(h[0], h[1])
	}

	res, err := ctx.client.Do(req)
	if err != nil {
		ctx.err = err
		return ctx
	}
	ctx.response = res

	return ctx
}

// Err get the error
func (ctx *Context) Err() error {
	return ctx.err
}

// Request get request
func (ctx *Context) Request() *http.Request {
	return ctx.request
}

// Response get response
func (ctx *Context) Response() *http.Response {
	return ctx.Do().response
}

func (ctx *Context) Bytes() []byte {
	body, err := ioutil.ReadAll(ctx.Response().Body)
	if err != nil {
		ctx.err = err
		return nil
	}

	err = ctx.response.Body.Close()
	if err != nil {
		ctx.err = err
		return nil
	}

	return body
}

// String get string response
func (ctx *Context) String() string {
	return string(ctx.Bytes())
}

// JSON unmarshal json response to v
func (ctx *Context) JSON(v interface{}) error {
	return json.Unmarshal(ctx.Bytes(), &v)
}

// GJSON parse body as json and provide searching for json strings
func (ctx *Context) GJSON(path string) gjson.Result {
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
