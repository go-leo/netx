package httpx

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/go-leo/iox"
)

type ResponseHelper struct {
	resp *http.Response
	err  error
}

func NewResponseHelper(resp *http.Response, err error) *ResponseHelper {
	return &ResponseHelper{resp: resp, err: err}
}

func (helper *ResponseHelper) Err() error {
	return helper.err
}

func (helper *ResponseHelper) StatusCode() int {
	return helper.resp.StatusCode
}

func (helper *ResponseHelper) Headers() http.Header {
	return helper.resp.Header
}

func (helper *ResponseHelper) LastModified() string {
	return helper.resp.Header.Get("Last-Modified")
}

func (helper *ResponseHelper) Etag() string {
	return helper.resp.Header.Get("Etag")
}

func (helper *ResponseHelper) CacheControl() string {
	return helper.resp.Header.Get("Cache-Control")
}

func (helper *ResponseHelper) Trailer() http.Header {
	return helper.resp.Trailer
}

func (helper *ResponseHelper) Cookies() []*http.Cookie {
	return helper.resp.Cookies()
}

func (helper *ResponseHelper) Body() io.ReadCloser {
	return helper.resp.Body
}

func (helper *ResponseHelper) BytesBody() ([]byte, error) {
	defer iox.QuiteClose(helper.resp.Body)
	b := new(bytes.Buffer)
	_, err := io.Copy(b, helper.resp.Body)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (helper *ResponseHelper) TextBody() (string, error) {
	defer iox.QuiteClose(helper.resp.Body)
	b := new(strings.Builder)
	_, err := io.Copy(b, helper.resp.Body)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (helper *ResponseHelper) ObjectBody(body any, unmarshal func([]byte, any) error) error {
	data, err := helper.BytesBody()
	if err != nil {
		return err
	}
	return unmarshal(data, body)
}

func (helper *ResponseHelper) JSONBody(body any) error {
	return helper.ObjectBody(body, json.Unmarshal)
}

func (helper *ResponseHelper) XMLBody(body any) error {
	return helper.ObjectBody(body, xml.Unmarshal)
}

func (helper *ResponseHelper) ProtobufBody(body proto.Message) error {
	unmarshal := func(data []byte, v any) error {
		m := v.(proto.Message)
		return proto.Unmarshal(data, m)
	}
	return helper.ObjectBody(body, unmarshal)
}

func (helper *ResponseHelper) GobBody(body proto.Message) error {
	unmarshal := func(data []byte, v any) error {
		return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
	}
	return helper.ObjectBody(body, unmarshal)
}

func (helper *ResponseHelper) FileBody(file io.Writer) error {
	defer iox.QuiteClose(helper.resp.Body)
	return iox.Copy(file, helper.Body())
}
