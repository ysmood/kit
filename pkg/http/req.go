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
	"github.com/ysmood/kit/pkg/utils"
)

// ReqContext ...
type ReqContext struct {
	client   *http.Client
	request  *http.Request
	response *http.Response

	method   string
	url      string
	header   [][]string
	jsonBody interface{}
	body     io.Reader
}

// JSONResult ...
type JSONResult = gjson.Result

// Req send http request
func Req(url string) *ReqContext {
	return &ReqContext{
		url: url,
	}
}

// Method ...
func (ctx *ReqContext) Method(m string) *ReqContext {
	ctx.method = m
	return ctx
}

// URL ...
func (ctx *ReqContext) URL(url string) *ReqContext {
	ctx.url = url
	return ctx
}

// Post ...
func (ctx *ReqContext) Post() *ReqContext {
	return ctx.Method(http.MethodPost)
}

// Put ...
func (ctx *ReqContext) Put() *ReqContext {
	return ctx.Method(http.MethodPut)
}

// Delete ...
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

// Body ...
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

// StringBody ...
func (ctx *ReqContext) StringBody(s string) *ReqContext {
	ctx.body = strings.NewReader(string(s))
	return ctx
}

// Do ...
func (ctx *ReqContext) Do() error {
	if ctx.client == nil {
		cookie, _ := cookiejar.New(nil)
		ctx.client = &http.Client{
			Jar: cookie,
		}
	}

	if ctx.jsonBody != nil {
		body, err := json.Marshal(ctx.jsonBody)
		if err != nil {
			return err
		}
		ctx.body = bytes.NewReader(body)
	}

	req, err := http.NewRequest(ctx.method, ctx.url, ctx.body)
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

// MustDo ...
func (ctx *ReqContext) MustDo() {
	utils.E(ctx.Do())
}

// Request get request
func (ctx *ReqContext) Request() *http.Request {
	return ctx.request
}

// Response get response
func (ctx *ReqContext) Response() (*http.Response, error) {
	err := ctx.Do()
	if err != nil {
		return nil, err
	}
	return ctx.response, nil
}

// MustResponse get response
func (ctx *ReqContext) MustResponse() *http.Response {
	return utils.E(ctx.Response())[0].(*http.Response)
}

// Bytes get response body as bytes
func (ctx *ReqContext) Bytes() ([]byte, error) {
	res, err := ctx.Response()
	if err != nil {
		return nil, err
	}
	return readBody(res.Body)
}

// MustBytes ...
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

// String get string response
func (ctx *ReqContext) String() (string, error) {
	s, err := ctx.Bytes()
	return string(s), err
}

// MustString ...
func (ctx *ReqContext) MustString() string {
	return string(ctx.MustBytes())
}

// JSON parse body as json and provide searching for json strings
func (ctx *ReqContext) JSON() (*JSONResult, error) {
	b, err := ctx.Bytes()
	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(b)
	return &r, nil
}

// MustJSON ...
func (ctx *ReqContext) MustJSON() *JSONResult {
	return utils.E(ctx.JSON())[0].(*gjson.Result)
}

func paramsToForm(params []interface{}) map[string]interface{} {
	f := map[string]interface{}{}
	l := len(params) - 1
	for i := 0; i < l; i += 2 {
		f[params[i].(string)] = params[i+1]
	}
	return f
}
