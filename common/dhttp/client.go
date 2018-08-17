package dhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Header map[string]string

var cli = &http.Client{Timeout: 5 * time.Second}

func Get(uri string, header Header, params url.Values) ([]byte, error) {
	finalUrl := uri
	if len(params) > 0 {
		finalUrl = fmt.Sprintf("%s?%s", finalUrl, params.Encode())
	}

	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return []byte{}, err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return data, errors.New(fmt.Sprintf("status code %d", resp.StatusCode))
	}

	return data, nil
}

func Post(uri string, header Header, params url.Values, body interface{}) ([]byte, error) {
	finalUrl := uri
	if len(params) > 0 {
		finalUrl = fmt.Sprintf("%s?%s", finalUrl, params.Encode())
	}

	bodyValue := reflect.ValueOf(body)
	if bodyValue.Kind() == reflect.Ptr {
		bodyValue = bodyValue.Elem()
	}

	var bodyStr string
	switch bodyValue.Kind() {
	case reflect.String:
		bodyStr = body.(string)
	case reflect.Map, reflect.Struct:
		bytes, err := json.Marshal(body)
		if err != nil {
			return []byte{}, err
		}
		bodyStr = string(bytes)
	}

	req, err := http.NewRequest("POST", finalUrl, strings.NewReader(bodyStr))
	if err != nil {
		return []byte{}, err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return data, errors.New(fmt.Sprintf("status code %d", resp.StatusCode))
	}

	return data, nil
}
