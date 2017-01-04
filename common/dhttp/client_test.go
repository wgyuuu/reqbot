package dhttp

import (
	"testing"

	"fmt"
	"net/http"

	"net/url"

	"encoding/json"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type HttpSuite struct {
	url     string
	port    int
	handler func(http.ResponseWriter, *http.Request)
}

func init() {
	suite := check.Suite(&HttpSuite{
		url:  "127.0.0.1",
		port: 8001,
		handler: func(resp http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			body := make([]byte, 512)
			n, _ := req.Body.Read(body)
			body = body[:n]

			header, _ := json.Marshal(req.Header)

			res := struct {
				Host   string
				Method string
				Body   string
				Header string
			}{
				req.Host, req.Method, string(body), string(header),
			}

			resJson, _ := json.Marshal(res)
			resp.Write(resJson)
		},
	}).(*HttpSuite)

	var mux = http.NewServeMux()
	mux.HandleFunc("/", suite.handler)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", suite.url, suite.port), mux)
}

func (h *HttpSuite) TestPost(c *check.C) {
	header := make(Header)
	header["AA"] = "BB"
	var values = url.Values{}
	values.Set("key", "1")

	data, err := Post(fmt.Sprintf("http://%s:%d", h.url, h.port), header, values, "{1}")
	c.Log("res:", data)
	c.Assert(err, check.IsNil)
}

func (h *HttpSuite) TestGet(c *check.C) {
	header := make(Header)
	header["AA"] = "BB"
	var values = url.Values{}
	values.Set("key", "1")

	data, err := Get(fmt.Sprintf("http://%s:%d", h.url, h.port), header, values)
	c.Log("res:", data)
	c.Assert(err, check.IsNil)
}
