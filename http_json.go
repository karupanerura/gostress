package gostress

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type JsonRequestEncoder struct{}
type JsonResponseDecoder struct{}

func (c *JsonRequestEncoder) GetContentType() string {
	return "application/json"
}

func (c *JsonRequestEncoder) Encode(data interface{}) (io.Reader, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(json), nil
}

func (c *JsonResponseDecoder) SupportedContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "application/json")
}

func (c *JsonResponseDecoder) Decode(reader io.Reader) (interface{}, error) {
	decoder := json.NewDecoder(reader)

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
