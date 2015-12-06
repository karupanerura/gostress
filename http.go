package gostress

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type RequestEncoder interface {
	GetContentType() string
	Encode(interface{}) (io.Reader, error)
}

type ResponseDecoder interface {
	SupportedContentType(string) bool
	Decode(io.Reader) (interface{}, error)
}

type ServerConfig struct {
	Hostname string
	Secure   bool
}

type HttpClientConfig struct {
	Server              ServerConfig
	Headers             map[string]string
	UserAgent           string
	MaxIdleConnsPerHost int
	RequestEncoder      RequestEncoder
	ResponseDecoder     ResponseDecoder
}

type HttpResponse struct {
	StatusCode int
	Header     http.Header
	Content    interface{}
}

type HttpClient struct {
	Client http.Client
	Config HttpClientConfig
}

func NewHttpClient(config HttpClientConfig) *HttpClient {
	return &HttpClient{
		Client: http.Client{
			Timeout: 0,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			},
		},
		Config: config,
	}
}

func (c *ServerConfig) MakeUrl(path string) string {
	if c.Secure {
		return fmt.Sprintf("https://%s%s", c.Hostname, path)
	}
	return fmt.Sprintf("http://%s%s", c.Hostname, path)
}

func (c *HttpClient) Request(method, path string, headers map[string]string, content interface{}) (*HttpResponse, error) {
	req, err := c.makeRequest(method, path, headers, content)
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return c.parseResponse(res)
}

func (c *HttpClient) makeRequest(method, path string, headers map[string]string, content interface{}) (*http.Request, error) {
	uri, body, err := c.preformContent(method, path, content)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	if content != nil {
		contentType := c.Config.RequestEncoder.GetContentType()
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("User-Agent", c.Config.UserAgent)
	for k, v := range c.Config.Headers {
		req.Header.Set(k, v)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	return req, nil
}

func (c *HttpClient) preformContent(method, path string, content interface{}) (uri string, body io.Reader, err error) {
	uri = c.Config.Server.MakeUrl(path)

	if content != nil {
		if isContentMethod(method) {
			body, err = c.Config.RequestEncoder.Encode(content)
			if err != nil {
				return
			}
		} else {
			query := url.Values{}
			switch content := content.(type) {
			case map[string]string:
				for k, v := range content {
					query.Set(k, v)
				}
			case map[fmt.Stringer]string:
				for k, v := range content {
					query.Set(k.String(), v)
				}
			case map[string]fmt.Stringer:
				for k, v := range content {
					query.Set(k, v.String())
				}
			case map[fmt.Stringer]fmt.Stringer:
				for k, v := range content {
					query.Set(k.String(), v.String())
				}
			}
			uri = uri + "?" + query.Encode()
		}
	}

	return
}

func (c *HttpClient) parseResponse(res *http.Response) (*HttpResponse, error) {
	defer res.Body.Close()
	if decoder := c.Config.ResponseDecoder; res.ContentLength > 0 && decoder != nil {
		if contentType := res.Header.Get("Content-Type"); decoder.SupportedContentType(contentType) {
			content, err := decoder.Decode(res.Body)
			if err != nil {
				return nil, err
			}
			res := &HttpResponse{
				StatusCode: res.StatusCode,
				Header:     res.Header,
				Content:    content,
			}
			return res, err
		} else {
			buf := make([]byte, res.ContentLength)
			io.ReadFull(res.Body, buf)
			log.Print("Content-Type: " + contentType)
			log.Print(string(buf))
		}
	}

	return &HttpResponse{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Content:    map[string]string{},
	}, nil
}

func isContentMethod(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}
