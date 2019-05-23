package gokit_test

import (
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	g "github.com/ysmood/gokit"
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
	r, _ := g.GenerateRandomString(5)
	path = "/" + r
	url = "http://127.0.0.1:" + port + path
	return path, url
}

func (s *RequestSuite) SetupSuite() {
	wait := make(chan g.Nil)

	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()

		ln, _ := net.Listen("tcp", ":0")

		s.router = r
		s.listener = ln

		wait <- g.Nil{}

		http.Serve(ln, r)
	}()

	<-wait
}

func (s *RequestSuite) TearDownSuite() {
	s.listener.Close()
}

func (s *RequestSuite) TestGetString() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		c.String(200, "ok")
	})

	s.Equal("ok", g.Req(url).String())
}

func (s *RequestSuite) TestGetStringWithQuery() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		v, _ := c.GetQuery("a")
		s.Equal("1", v)
	})

	c := g.Req(url).Query("a", "1")

	s.Nil(c.Err())
}

func (s *RequestSuite) TestGetJSON() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		c.String(200, `{ "A": "ok" }`)
	})

	c := g.Req(url)

	type t struct {
		A string
	}

	var data t
	c.JSON(&data)

	s.Equal("ok", data.A)
}

func (s *RequestSuite) TestGetGJSON() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		c.String(200, `{ "deep": { "path": 1 } }`)
	})

	c := g.Req(url)

	s.Equal(int64(1), c.GJSON("deep.path").Int())
}

func (s *RequestSuite) TestPostForm() {
	path, url := s.path()

	s.router.POST(path, func(c *gin.Context) {
		out, _ := c.GetPostForm("a")
		c.String(200, out)
	})

	c := g.Req(url).Post().Form("a", "val")
	s.Equal("val", c.String())
}

func (s *RequestSuite) TestPostBytes() {
	path, url := s.path()

	s.router.POST(path, func(c *gin.Context) {
		out, _ := c.GetRawData()
		c.Data(200, "", out)
	})

	c := g.Req(url).Post().Body(strings.NewReader("raw"))
	s.Equal("raw", c.String())
}

func (s *RequestSuite) TestPutString() {
	path, url := s.path()

	s.router.PUT(path, func(c *gin.Context) {
		out, _ := c.GetRawData()
		c.Data(200, "", out)
	})

	c := g.Req(url).Put().StringBody("raw")
	s.Equal("raw", c.String())
}

func (s *RequestSuite) TestDelete() {
	path, url := s.path()

	s.router.DELETE(path, func(c *gin.Context) {
		c.String(200, "ok")
	})

	c := g.Req(url).Delete()
	s.Equal("ok", c.String())
}

func (s *RequestSuite) TestPostJSON() {
	path, url := s.path()

	s.router.POST(path, func(c *gin.Context) {
		data, _ := c.GetRawData()
		c.String(200, g.JSON(data).Get("A").String())
	})

	c := g.Req(url).Post().JSONBody(struct{ A string }{"ok"})
	s.Equal("ok", c.String())
}

func (s *RequestSuite) TestHeader() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		h := c.GetHeader("test")
		c.String(200, h)
	})

	c := g.Req(url).Header("test", "ok").Do()
	s.Equal("ok", c.String())
}

func (s *RequestSuite) TestReuseCookie() {
	path, url := s.path()

	var cookieA string
	var header string

	s.router.GET(path, func(c *gin.Context) {
		cookieA, _ = c.Cookie("t")
		header = c.GetHeader("a")
		c.SetCookie("t", "val", 3600, "", "", false, true)
	})

	c := g.Req(url).Do()
	c.Header("a", "b").Req(url)

	s.Equal("val", cookieA)
	s.Equal("b", header)
}
