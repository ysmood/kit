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

	client, err := g.Req(url)

	g.E(err)

	s.Equal("ok", client.String())
}

func (s *RequestSuite) TestGetStringWithQuery() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		v, _ := c.GetQuery("a")
		s.Equal("1", v)
	})

	_, err := g.Req(url, g.QueryParams{"a": "1"})

	g.E(err)
}

func (s *RequestSuite) TestGetJSON() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		c.String(200, `{ "A": "ok" }`)
	})

	client, err := g.Req(url)

	g.E(err)

	type t struct {
		A string
	}

	var data t
	client.JSON(&data)

	s.Equal("ok", data.A)
}

func (s *RequestSuite) TestGetGJSON() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		c.String(200, `{ "deep": { "path": 1 } }`)
	})

	client, err := g.Req(url)

	g.E(err)

	s.Equal(int64(1), client.GJSON().Get("deep.path").Int())
}

func (s *RequestSuite) TestPostForm() {
	path, url := s.path()

	s.router.POST(path, func(c *gin.Context) {
		out, _ := c.GetPostForm("a")
		s.Equal("val", out)
	})

	_, err := g.Req(
		g.Method(http.MethodPost),
		url,
		g.FormParams{"a": "val"},
	)
	g.E(err)
}

func (s *RequestSuite) TestPostBytes() {
	path, url := s.path()

	s.router.POST(path, func(c *gin.Context) {
		out, _ := c.GetRawData()
		s.Equal([]byte("raw"), out)
	})

	_, err := g.Req(
		g.Method(http.MethodPost),
		url,
		strings.NewReader("raw"),
	)
	g.E(err)
}

func (s *RequestSuite) TestPostJSON() {
	path, url := s.path()

	s.router.POST(path, func(c *gin.Context) {
		data, _ := c.GetRawData()
		s.Equal("ok", g.JSON(data).Get("A").String())
	})

	_, err := g.Req(
		g.Method(http.MethodPost),
		url,
		g.JSONBody(struct{ A string }{"ok"}),
	)
	g.E(err)
}

func (s *RequestSuite) TestHeader() {
	path, url := s.path()

	s.router.GET(path, func(c *gin.Context) {
		h := c.GetHeader("test")
		s.Equal("ok", h)
	})

	_, err := g.Req(
		url,
		g.Header{"test": "ok"},
	)
	g.E(err)
}

func (s *RequestSuite) TestReuseCookie() {
	path, url := s.path()

	var cookieVal string

	s.router.GET(path, func(c *gin.Context) {
		cookieVal, _ = c.Cookie("t")
		c.SetCookie("t", "val", 3600, "", "", false, true)
	})

	client, _ := g.Req(url)
	client.Req(url)

	s.Equal("val", cookieVal)
}
