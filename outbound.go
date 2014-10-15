package outbound

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Request struct {
	headers           []headerTuple
	Method            string
	Uri               string
	Body              interface{}
	QueryString       interface{}
	Timeout           time.Duration
	ContentType       string
	Accept            string
	Host              string
	UserAgent         string
	Insecure          bool
	MaxRedirects      int
	Proxy             string
	Compression       *compression
	BasicAuthUsername string
	BasicAuthPassword string
}

type Response struct {
	StatusCode    int
	ContentLength int64
	Body          *Body
	Header        http.Header
}

type Body struct {
	reader           io.ReadCloser
	compressedReader io.ReadCloser
}

type Error struct {
	timeout bool
	Err     error
}

func (e *Error) Timeout() bool {
	return e.timeout
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (b *Body) Read(p []byte) (int, error) {
	if b.compressedReader != nil {
		return b.compressedReader.Read(p)
	}
	return b.reader.Read(p)
}

func (b *Body) Close() error {
	err := b.reader.Close()
	if b.compressedReader != nil {
		return b.compressedReader.Close()
	}
	return err
}

func paramParse(query interface{}) (string, error) {
	var (
		v = &url.Values{}
		s = reflect.ValueOf(query)
		t = reflect.TypeOf(query)
	)

	switch query.(type) {
	case url.Values:
		return query.(url.Values).Encode(), nil
	default:
		for i := 0; i < s.NumField(); i++ {
			v.Add(strings.ToLower(t.Field(i).Name), fmt.Sprintf("%v", s.Field(i).Interface()))
		}
		return v.Encode(), nil
	}
}

func prepareRequestBody(b interface{}) (io.Reader, error) {
	switch b.(type) {
	case string:
		return strings.NewReader(b.(string)), nil
	case io.Reader:
		return b.(io.Reader), nil
	case []byte:
		return bytes.NewReader(b.([]byte)), nil
	case nil:
		return nil, nil
	default:
		j, err := json.Marshal(b)
		if err == nil {
			return bytes.NewReader(j), nil
		}
		return nil, err
	}
}
func (r Request) Send() (*Response, error) {
	var req *http.Request
	var er error
	var transport = defaultTransport
	var client = defaultClient

	if r.QueryString != nil {
		param, e := paramParse(r.QueryString)
		if e != nil {
			return nil, &Error{Err: e}
		}
		r.Uri = r.Uri + "?" + param
	}

	var bodyReader io.Reader
	req, er = http.NewRequest(r.Method, r.Uri, bodyReader)

	if er != nil {
		return nil, &Error{Err: er}
	}

	res, err := client.Do(req)
	if timer != nil {
		timer.Stop()
	}

	if redirectSanitizer(res.StatusCode) && r.MaxRedirects > 0 {
		loc, _ := res.Location()
		r.MaxRedirects--
		r.Uri = loc.String()
		return r.Send()
	}

}

func redirectSanitizer(status int) bool {
	switch status {
	case http.StatusMovedPermanently:
		return true
	case http.StatusFound:
		return true
	case http.StatusSeeOther:
		return true
	case http.StatusTemporaryRedirect:
		return true
	default:
		return false
	}
}
