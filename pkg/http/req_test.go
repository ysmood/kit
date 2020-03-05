package http_test

import (
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/ysmood/kit"
)

type RequestSuite struct {
	suite.Suite
	router   *gin.Engine
	listener net.Listener
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}

func (s *RequestSuite) path() (path, url string) {
	_, port, _ := net.SplitHostPort(s.listener.Addr().String())
	r := kit.RandString(5)
	path = "/" + r
	url = "http://127.0.0.1:" + port + path
	return path, url
}

func (s *RequestSuite) SetupSuite() {
	server := kit.MustServer(":0")
	s.listener = server.Listener
	s.router = server.Engine

	go server.MustDo()
}

func (s *RequestSuite) TestTimeout() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		time.Sleep(time.Hour)
	})

	err := kit.Req(url).Timeout(time.Nanosecond).Do()
	s.Error(err)
}

func (s *RequestSuite) TestTimeoutDone() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {})

	kit.Req(url).Timeout(time.Hour).MustDo()
}

func (s *RequestSuite) TestGetMustString() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		c.String(200, "ok")
	})

	c := kit.Req(url)
	req, err := c.Request()
	kit.E(err)

	s.Equal("ok", c.MustString())
	s.Equal(url, req.URL.String())
}

func (s *RequestSuite) TestGetMustResponse() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		c.String(200, "ok")
	})

	s.Equal(200, kit.Req(url).MustResponse().StatusCode)
}

func (s *RequestSuite) TestGetString() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		c.String(200, "ok")
	})

	c := kit.Req(url)
	req, err := c.Request()
	kit.E(err)

	s.Equal("ok", kit.E(c.String())[0].(string))
	s.Equal(url, req.URL.String())
}

func (s *RequestSuite) TestGetStringErr() {
	_, err := kit.Req("").String()
	s.EqualError(err, "Get \"\": unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestSetClient() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		c.String(200, "ok")
	})

	c := kit.Req(url).Client(&http.Client{})

	s.Equal("ok", c.MustString())
}

func (s *RequestSuite) TestMethodErr() {
	err := kit.Req("").Method("あ").Do()
	s.EqualError(err, "net/http: invalid method \"あ\"")
}

func (s *RequestSuite) TestURLErr() {
	s.EqualError(kit.ErrArg(kit.Req("").Do()), "Get \"\": unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestRequestErr() {
	err := kit.Req("").Do()
	s.EqualError(err, "Get \"\": unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestGetStringWithQuery() {
	path, url := s.path()
	var v string

	s.router.GET(path, func(c kit.GinContext) {
		v, _ = c.GetQuery("a")
	})

	kit.Req(url).Query("a", "1").MustDo()

	s.Equal("1", v)
}

func (s *RequestSuite) TestQueryWrongKVAmount() {
	s.Panics(func() {
		kit.Req("").Query("a")
	})
}

func (s *RequestSuite) TestGetJSON() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		c.String(200, `{ "deep": { "path": 1 } }`)
	})

	c := kit.Req(url)

	s.Equal(int64(1), c.MustJSON().Get("deep.path").Int())
}

func (s *RequestSuite) TestGetJSONErr() {
	s.EqualError(kit.ErrArg(kit.Req("").JSON()), "Get \"\": unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestPostForm() {
	path, url := s.path()

	s.router.POST(path, func(c kit.GinContext) {
		out, _ := c.GetPostForm("a")
		c.String(200, out)
	})

	c := kit.Req(url).Post().Form("a", "val")
	s.Equal("val", c.MustString())
}

func (s *RequestSuite) TestPostBytes() {
	path, url := s.path()

	s.router.POST(path, func(c kit.GinContext) {
		out, _ := c.GetRawData()
		c.Data(200, "", out)
	})

	c := kit.Req(url).Post().Body(strings.NewReader("raw"))
	s.Equal("raw", c.MustString())
}

func (s *RequestSuite) TestPutString() {
	path, url := s.path()

	s.router.PUT(path, func(c kit.GinContext) {
		out, _ := c.GetRawData()
		c.Data(200, "", out)
	})

	c := kit.Req(url).Put().StringBody("raw")
	s.Equal("raw", c.MustString())
}

func (s *RequestSuite) TestDelete() {
	path, url := s.path()

	s.router.DELETE(path, func(c kit.GinContext) {
		c.String(200, "ok")
	})

	c := kit.Req(url).Delete()
	s.Equal("ok", c.MustString())
}

func (s *RequestSuite) TestPostJSON() {
	path, url := s.path()

	s.router.POST(path, func(c kit.GinContext) {
		data, _ := c.GetRawData()
		c.String(200, kit.JSON(data).Get("A").String())
	})

	c := kit.Req(url).Post().JSONBody(struct{ A string }{"ok"})
	s.Equal("ok", c.MustString())
}

func (s *RequestSuite) TestJSONBodyError() {
	v := make(chan kit.Nil)
	err := kit.Req("").JSONBody(v).Do()
	s.EqualError(err, "json: unsupported type: chan utils.Nil")
}

func (s *RequestSuite) TestHeader() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		h := c.Request.Header["Test"]
		s.Equal("ok_ok", h[0]+h[1])
		s.Equal("test", c.Request.Host)
	})

	kit.Req(url).Host("test").Header("test", "ok", "test", "_ok").MustDo()
}

func (s *RequestSuite) TestReuseCookie() {
	path, url := s.path()

	var cookieA string
	var header string

	s.router.GET(path, func(c kit.GinContext) {
		cookieA, _ = c.Cookie("t")
		header = c.GetHeader("a")
		c.SetCookie("t", "val", 3600, "", "", false, true)
	})

	c := kit.Req(url)
	c.MustDo()
	c.URL(url).Header("a", "b").MustDo()

	s.Equal("val", cookieA)
	s.Equal("b", header)
}

func (s *RequestSuite) TestMustCurl() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		c.JSON(200, map[string]string{"a": "b"})
	})

	c := kit.Req(url).JSONBody([]int{10})

	res, err := c.Response()
	kit.E(err)

	expected := kit.S(`curl -X GET {{.url}} \
  -H 'Content-Type: application/json; charset=utf-8' \
  -d '[10]'

HTTP/1.1 200 OK
Content-Length: 10
Content-Type: application/json; charset=utf-8
Date: {{.date}}

{
  "a": "b"
}`, "url", url, "date", res.Header.Get("Date"))

	s.Equal(expected, c.MustCurl())
}

func (s *RequestSuite) TestMustCurlEmptyBody() {
	path, url := s.path()

	s.router.GET(path, func(c kit.GinContext) {
		kit.E(c.Writer.Write([]byte{0xff, 0xfe, 0xfd}))
	})

	c := kit.Req(url)

	res, err := c.Response()
	kit.E(err)

	expected := kit.S(`curl -X GET {{.url}}

HTTP/1.1 200 OK
Content-Length: 3
Content-Type: text/plain; charset=utf-8
Date: {{.date}}

//79`, "url", url, "date", res.Header.Get("Date"))

	s.Equal(expected, c.MustCurl())
}
