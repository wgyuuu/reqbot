package controller

import (
	"net/url"
	"strings"
	"time"

	"github.com/wgyuuu/reqbot/common/dhttp"
	"github.com/wgyuuu/reqbot/common/dlog"
)

func SendRequest(uri, method string, mapHeaders, mapParams map[string]string, body string) {
	var err error
	var result []byte

	start := time.Now()
	defer func() {
		if err == nil {
			return
		}

		fields := map[string]interface{}{
			"url":    uri,
			"method": method,
			"timer":  time.Since(start),
		}
		dlog.NewEntry(fields).Info2("params", mapParams, "body", body, "result", result, "error", err)
	}()

	values := make(url.Values)
	for key, value := range mapParams {
		values.Add(key, value)
	}

	switch strings.ToUpper(method) {
	case "POST":
		result, err = dhttp.Post(uri, mapHeaders, values, body)
	case "GET":
		result, err = dhttp.Get(uri, mapHeaders, values)
	}
}
