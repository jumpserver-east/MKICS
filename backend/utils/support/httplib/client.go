package httplib

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const miniTimeout = time.Second * 30

func NewClient(baseUrl string, timeout time.Duration) (*Client, error) {
	_, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	jar := &customCookieJar{
		data: map[string]string{},
	}
	if timeout < miniTimeout {
		timeout = miniTimeout
	}
	client := http.Client{
		Timeout:   timeout,
		Jar:       jar,
		Transport: NewTransport(true),
	}
	return &Client{
		baseUrl: baseUrl,
		Headers: headers,
		http:    client,
	}, nil
}

type AuthSign interface {
	Sign(req *http.Request) error
}

type Client struct {
	Headers   map[string]string
	Auth      AuthSign
	baseUrl   string
	http      http.Client
	basicAuth []string
}

func NewTransport(insecure bool) http.RoundTripper {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return transport
}

func (c *Client) SetAuthSign(auth AuthSign) {
	c.Auth = auth
}

func (c *Client) SetHeader(k, v string) {
	c.Headers[k] = v
}

func (c *Client) SetBasicAuth(username, password string) {
	c.basicAuth = append(c.basicAuth, username)
	c.basicAuth = append(c.basicAuth, password)
}

func (c *Client) marshalData(data interface{}) (reader io.Reader, error error) {
	dataRaw, err := json.Marshal(data)
	if err != nil {
		return
	}
	reader = bytes.NewReader(dataRaw)
	return
}

func (c *Client) parseQueryUrl(reqUrl string, params []map[string]string) string {
	if len(params) < 1 {
		return reqUrl
	}
	query := url.Values{}
	for _, item := range params {
		for k, v := range item {
			query.Add(k, v)
		}
	}
	if strings.Contains(reqUrl, "?") {
		reqUrl += "&" + query.Encode()
	} else {
		reqUrl += "?" + query.Encode()
	}
	return reqUrl
}

func (c *Client) parseUrl(reqUrl string, params []map[string]string) string {
	reqUrl = c.parseQueryUrl(reqUrl, params)
	if c.baseUrl != "" {
		reqUrl = strings.TrimSuffix(c.baseUrl, "/") + reqUrl
	}
	return reqUrl
}
func (c *Client) setReqAuthHeader(r *http.Request) error {
	if len(c.basicAuth) == 2 {
		r.SetBasicAuth(c.basicAuth[0], c.basicAuth[1])
		return nil
	}
	if c.Auth != nil {
		return c.Auth.Sign(r)
	}
	return nil

}

func (c *Client) SetReqHeaders(req *http.Request) error {
	if len(c.Headers) != 0 {
		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.setReqAuthHeader(req)
}

func (c *Client) NewRequest(method, url string, body interface{}, params []map[string]string) (req *http.Request, err error) {
	url = c.parseUrl(url, params)
	reader, err := c.marshalData(body)
	if err != nil {
		return
	}
	req, err = http.NewRequest(method, url, reader)
	if err != nil {
		return req, err
	}
	err = c.SetReqHeaders(req)
	return req, err
}

// Do wrapper http.Client Do() for using auth and error handle
// params:
//  1. query string if set {"name": "ibuler"}
func (c *Client) Do(method, url string, data, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	req, err := c.NewRequest(method, url, data, params)
	if err != nil {
		return
	}
	resp, err = c.http.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("%s %s failed, get code: %d, %s", req.Method, req.URL, resp.StatusCode, body)
		err = errors.New(msg)
		return
	}

	// If is buffer return the raw response body
	if buf, ok := res.(*bytes.Buffer); ok {
		buf.Write(body)
		return
	}
	// Unmarshal response body to result struct
	if res != nil {
		switch {
		case strings.Contains(resp.Header.Get("Content-Type"), "application/json"):
			err = json.Unmarshal(body, res)
			if err != nil {
				msg := fmt.Sprintf("%s %s failed, unmarshal '%s' response failed: %s", req.Method, req.URL, body[:12], err)
				err = errors.New(msg)
				return
			}
		}
	}
	return
}

func (c *Client) Get(url string, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("GET", url, nil, res, params...)
}

func (c *Client) Post(url string, data interface{}, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("POST", url, data, res, params...)
}

func (c *Client) Delete(url string, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("DELETE", url, nil, res, params...)
}

func (c *Client) Put(url string, data interface{}, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("PUT", url, data, res, params...)
}

func (c *Client) Patch(url string, data interface{}, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("PATCH", url, data, res, params...)
}
