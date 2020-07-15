package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"
)

// GetContentType from filename
func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	ct := "application/octet-stream"

	switch ext {
	case ".js":
		ct = "text/javascript;charset=UTF-8"
	case ".json":
		ct = "application/json"
	case ".png":
		ct = "image/png"
	case ".html":
		ct = "text/html"
	}

	return ct
}

// GetHTTPBody extracts a body from http request
func GetHTTPBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("Request did not have a body")
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	var bodyData interface{}
	if err = json.Unmarshal(data, &bodyData); err != nil {
		return nil, err
	}

	return data, nil
}

// HTTPAsJSON writes json response
func HTTPAsJSON(w http.ResponseWriter, d interface{}) {
	buffer := []byte{}
	switch v := d.(type) {
	case []byte:
		buffer = v
	default:
		buffer, _ = json.MarshalIndent(d, "", "    ")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(buffer)
}

// HTTPErrorAsJSON writes a error json response
func HTTPErrorAsJSON(w http.ResponseWriter, s int, err interface{}) {
	errStr := "undefined error type"
	switch v := err.(type) {
	case string:
		errStr = v
	case error:
		errStr = v.Error()
	case []byte:
		errStr = string(v)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	w.Write([]byte(fmt.Sprintf(`{"status":"error","message":"%s"}`, errStr)))
}

// HeadRequest a site
func HeadRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	return client.Do(req)
}

// GetRequest a site
func GetRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	return client.Do(req)
}

// HasContentType checks for mime
func HasContentType(r *http.Response, mimetype string) bool {
	contentType := r.Header.Get("Content-type")

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
