package http_test

import (
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	. "github.com/ysmood/gokit/pkg/http"
	. "github.com/ysmood/gokit/pkg/utils"
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
	r := GenerateRandomString(5)
	path = "/" + r
	url = "http://127.0.0.1:" + port + path
	return path, url
}

func (s *RequestSuite) SetupSuite() {
	server := MustServer(":0")
	s.listener = server.Listener
	s.router = server.Handler

	go server.MustDo()
}

func (s *RequestSuite) TestGetMustString() {
	path, url := s.path()

	s.router.GET(path, func(c GinContext) {
		c.String(200, "ok")
	})

	c := Req(url)

	s.Equal("ok", c.MustString())
	s.Equal(url, c.Request().URL.String())
}

func (s *RequestSuite) TestGetMustResponse() {
	path, url := s.path()

	s.router.GET(path, func(c GinContext) {
		c.String(200, "ok")
	})

	s.Equal(200, Req(url).MustResponse().StatusCode)
}

func (s *RequestSuite) TestGetString() {
	path, url := s.path()

	s.router.GET(path, func(c GinContext) {
		c.String(200, "ok")
	})

	c := Req(url)

	s.Equal("ok", E(c.String())[0].(string))
	s.Equal(url, c.Request().URL.String())
}

func (s *RequestSuite) TestGetStringErr() {
	_, err := Req("").String()
	s.EqualError(err, "Get : unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestSetClient() {
	path, url := s.path()

	s.router.GET(path, func(c GinContext) {
		c.String(200, "ok")
	})

	c := Req(url).Client(&http.Client{})

	s.Equal("ok", c.MustString())
}

func (s *RequestSuite) TestMethodErr() {
	err := Req("").Method("あ").Do()
	s.EqualError(err, "net/http: invalid method \"あ\"")
}

func (s *RequestSuite) TestURLErr() {
	s.EqualError(ErrArg(Req("").Do()), "Get : unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestGetStringWithQuery() {
	path, url := s.path()
	var v string

	s.router.GET(path, func(c GinContext) {
		v, _ = c.GetQuery("a")
	})

	Req(url).Query("a", "1").MustDo()

	s.Equal("1", v)
}

func (s *RequestSuite) TestGetJSON() {
	path, url := s.path()

	s.router.GET(path, func(c GinContext) {
		c.String(200, `{ "A": "ok" }`)
	})

	c := Req(url)

	type t struct {
		A string
	}

	var data t
	E(c.JSON(&data))

	s.Equal("ok", data.A)
}
func (s *RequestSuite) TestGetJSONErr() {
	var v interface{}
	s.EqualError(ErrArg(Req("").JSON(v)), "Get : unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestGetGJSON() {
	path, url := s.path()

	s.router.GET(path, func(c GinContext) {
		c.String(200, `{ "deep": { "path": 1 } }`)
	})

	c := Req(url)

	s.Equal(int64(1), c.MustGJSON("deep.path").Int())
}

func (s *RequestSuite) TestGetGJSONErr() {
	s.EqualError(ErrArg(Req("").GJSON("deep.path")), "Get : unsupported protocol scheme \"\"")
}

func (s *RequestSuite) TestPostForm() {
	path, url := s.path()

	s.router.POST(path, func(c GinContext) {
		out, _ := c.GetPostForm("a")
		c.String(200, out)
	})

	c := Req(url).Post().Form("a", "val")
	s.Equal("val", c.MustString())
}

func (s *RequestSuite) TestPostBytes() {
	path, url := s.path()

	s.router.POST(path, func(c GinContext) {
		out, _ := c.GetRawData()
		c.Data(200, "", out)
	})

	c := Req(url).Post().Body(strings.NewReader("raw"))
	s.Equal("raw", c.MustString())
}

func (s *RequestSuite) TestPutString() {
	path, url := s.path()

	s.router.PUT(path, func(c GinContext) {
		out, _ := c.GetRawData()
		c.Data(200, "", out)
	})

	c := Req(url).Put().StringBody("raw")
	s.Equal("raw", c.MustString())
}

func (s *RequestSuite) TestDelete() {
	path, url := s.path()

	s.router.DELETE(path, func(c GinContext) {
		c.String(200, "ok")
	})

	c := Req(url).Delete()
	s.Equal("ok", c.MustString())
}

func (s *RequestSuite) TestPostJSON() {
	path, url := s.path()

	s.router.POST(path, func(c GinContext) {
		data, _ := c.GetRawData()
		c.String(200, JSON(data).Get("A").String())
	})

	c := Req(url).Post().JSONBody(struct{ A string }{"ok"})
	s.Equal("ok", c.MustString())
}

func (s *RequestSuite) TestJSONBodyError() {
	v := make(chan Nil)
	err := Req("").JSONBody(v).Do()
	s.EqualError(err, "json: unsupported type: chan utils.Nil")
}

func (s *RequestSuite) TestHeader() {
	path, url := s.path()

	s.router.GET(path, func(c GinContext) {
		h := c.GetHeader("test")
		c.String(200, h)
	})

	c := Req(url).Header("test", "ok")
	s.Equal("ok", c.MustString())
}

func (s *RequestSuite) TestReuseCookie() {
	path, url := s.path()

	var cookieA string
	var header string

	s.router.GET(path, func(c GinContext) {
		cookieA, _ = c.Cookie("t")
		header = c.GetHeader("a")
		c.SetCookie("t", "val", 3600, "", "", false, true)
	})

	c := Req(url)
	c.MustDo()
	c.URL(url).Header("a", "b").MustDo()

	s.Equal("val", cookieA)
	s.Equal("b", header)
}
