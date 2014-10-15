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

	req.Host = r.Host
	req.Header.Add("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", r.ContentType)
	req.Header.Add("Accept", r.Accept)

	if r.Compression != nil {
		req.Header.Add("Content-Encoding", r.Compression.ContentEncoding)
		req.Header.Add("Accept-Encoding", r.Compression.ContentEncoding)
	}

	if r.headers != nil {
		for _, header := range r.headers {
			req.Header.Add(header.name, header.value)
		}
	}

	if r.BasicAuthUsername != "" && r.BasicAuthPassword != "" {
		req.SetBasicAuth(r.BasicAuthUsername, r.BasicAuthPassword)
	}

	timeout := false
	var timer *time.Timer
	if r.Timeout > 0 {
		timer = time.AfterFunc(r.Timeout, func() {
			transport.CancelRequest(req)
			timeout = true
		})
	}

	res, err := client.Do(req)
	if timer != nil {
		timer.Stop()
	}

	if err != nil {
		if !timeout {
			switch err := err.(type) {
			case *net.OpError:
				timeout = err.Timeout()
			case *url.Error:
				if op, ok := err.Err.(*net.OpError); ok {
					timeout = op.Timeout()
				}
			}
		}

		return nil, &Error{timeout: timeout, Err: err}
	}

	if redirectSanitizer(res.StatusCode) && r.MaxRedirects > 0 {
		loc, _ := res.Location()
		r.MaxRedirects--
		r.Uri = loc.String()
		return r.Send()
	}

	if r.Compression != nil && strings.Contains(res.Header.Get("Content-Encoding"), r.Compression.ContentEncoding) {
		compressedReader, err := r.Compression.reader(res.Body)
		if err != nil {
			return nil, &Error{Err: err}
		}
		return &Response{StatusCode: res.StatusCode, ContentLength: res.ContentLength, Header: res.Header, Body: &Body{reader: res.Body, compressedReader: compressedReader}}, nil
	} else {
		return &Response{StatusCode: res.StatusCode, ContentLength: res.ContentLength, Header: res.Header, Body: &Body{reader: res.Body}}, nil
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
