package controller

import (
	"net/url"
	"time"

	"github.com/wgyuuu/reqbot/common/dhttp"
	"github.com/wgyuuu/reqbot/common/dlog"
)

func SendRequest(uri, method string, mapHeaders, mapParams map[string]string, body string) {
	var err error
	var result []byte

	start := time.Now()
	defer func() {
		fields := map[string]interface{}{
			"url":    uri,
			"method": method,
			"timer":  time.Since(start) / time.Millisecond,
		}
		dlog.NewEntry(fields).Debug2("params", mapParams, "body", body, "result", result)
	}()

	values := make(url.Values)
	for key, value := range mapParams {
		values.Add(key, value)
	}

	switch method {
	case "POST":
		result, err = dhttp.Post(uri, mapHeaders, values, body)
	case "GET":
		result, err = dhttp.Get(uri, mapHeaders, values)
	}
	if err != nil {
		dlog.Error("url", uri, "error", err)
	}
}
