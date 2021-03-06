// Package outbound is a multi-function http client for Golang. Yet another, if you will. I
// took a lot of best practices from various other libraries and combined their paradigms
// into a easy to use, but full featured, outbound http client.
//
// The package is a work in progress, and is considered stable enough for use. I see no
// major structure changes coming to the Request type. However, a new UDP type is in the
// works and will be inside a separate branch

package main

import (
	// "code.google.com/p/go.net/websocket"

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

// beginnings of a an outbound websocket client using channels
// type WebSocket struct {
// 	id     int
// 	ws     *websocket.Conn
// 	server *Server
// 	ch     chan *Message
// 	doneCh chan bool
// }

// Request is the core struct and entry point for creating connection objects
type Request struct {
	Host          string
	UserAgent     string
	SSL           bool
	MaxRedirects  int
	BasicAuthUser string
	BasicAuthPass string
	Method        string
	ContentType   string
	Accept        string
	URI           string
	Proxy         string
	Body          interface{}
	QueryString   interface{}
	Compression   *compression
	headers       []httpHeader
	Timeout       time.Duration
}

// Response is a container type that is used to ease of marshaling in the response
type Response struct {
	Body          *Body
	StatusCode    int
	Header        http.Header
	ContentLength int64
}

// all compression stuffs: https://github.com/rcoh/modata/blob/704c7b25baa2173d6a62403d4dc88b7a7f506727/src/diskv/compression.go
type compression struct {
	reader   func(buffer io.Reader) (io.ReadCloser, error)
	writer   func(buffer io.Writer) (io.WriteCloser, error)
	Encoding string
}

type httpHeader struct {
	name  string
	value string
}

type Body struct {
	reader     io.ReadCloser
	compressed io.ReadCloser
}

type Error struct {
	timeout bool
	Err     error
}

// func NewWebSocketClient(ws *websocket.Conn, server *Server) *WebSocket {

// 	if ws == nil {
// 		panic("ws cannot be nil")
// 	}

// 	if server == nil {
// 		panic("server cannot be nil")
// 	}

// 	maxId++
// 	ch := make(chan *Message, channelBufSize)
// 	doneCh := make(chan bool)

// 	return &WebSocket{maxId, ws, server, ch, doneCh}
// }

// func (c *WebSocket) Conn() *websocket.Conn {
// 	return c.ws
// }

// func (c *WebSocket) Write(msg *Message) {
// 	select {
// 	case c.ch <- msg:
// 	default:
// 		c.server.Del(c)
// 		err := fmt.Errorf("client %d is disconnected.", c.id)
// 		c.server.Err(err)
// 	}
// }

// func (c *WebSocket) Close() {
// 	c.doneCh <- true
// }

// func (c *WebSocket) Listen() {
// 	go c.listenWrite()
// 	c.listenRead()
// }

// func (c *WebSocket) listenWrite() {
// 	log.Println("Listening write to client")
// 	for {
// 		select {

// 		case msg := <-c.ch:
// 			log.Println("Send:", msg)
// 			websocket.JSON.Send(c.ws, msg)

// 		case <-c.doneCh:
// 			c.server.Del(c)
// 			c.doneCh <- true
// 			return
// 		}
// 	}
// }

// func (c *WebSocket) listenRead() {
// 	log.Println("Listening read from client")
// 	for {
// 		select {

// 		case <-c.doneCh:
// 			c.server.Del(c)
// 			c.doneCh <- true
// 			return

// 		default:
// 			var msg Message
// 			err := websocket.JSON.Receive(c.ws, &msg)
// 			if err == io.EOF {
// 				c.doneCh <- true
// 			} else if err != nil {
// 				c.server.Err(err)
// 			} else {
// 				c.server.SendAll(&msg)
// 			}
// 		}
// 	}
// }

// helper to cover all the redirect types
func isActualRedirect(status int) bool {
	switch status {
	case http.StatusMultipleChoices:
		return true
	case http.StatusMovedPermanently:
		return true
	case http.StatusFound:
		return true
	case http.StatusNotModified:
		return true
	case http.StatusUseProxy:
		return true
	case http.StatusSeeOther:
		return true
	case http.StatusTemporaryRedirect:
		return true
	default:
		return false
	}
}

func (e *Error) Timeout() bool {
	return e.timeout
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (b *Body) Read(p []byte) (int, error) {
	if b.compressed != nil {
		return b.compressed.Read(p)
	}

	return b.reader.Read(p)
}

func (b *Body) Close() error {
	err := b.reader.Close()

	if b.compressed != nil {
		return b.compressed.Close()
	}

	return err
}

func (b *Body) JsonToStruct(o interface{}) error {
	if body, err := ioutil.ReadAll(b); err != nil {
		return err
	} else if err := json.Unmarshal(body, o); err != nil {
		return err
	}

	return nil
}

func (b *Body) ToString() (string, error) {
	body, err := ioutil.ReadAll(b)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// https://github.com/shelakel/go-middleware/blob/e692288c2317fe9256547a7d28007fb1afc7db03/compression/gzip.go
func Gzip() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return gzip.NewReader(buffer)
	}

	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return gzip.NewWriter(buffer), nil
	}

	return &compression{
		writer:   writer,
		reader:   reader,
		Encoding: "gzip",
	}
}

// https://github.com/shelakel/go-middleware/blob/e692288c2317fe9256547a7d28007fb1afc7db03/compression/deflate.go
func Deflate() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return flate.NewReader(buffer), nil
	}

	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(buffer, -1)
	}

	return &compression{
		writer:   writer,
		reader:   reader,
		Encoding: "deflate",
	}
}

// https://github.com/rcoh/modata/blob/704c7b25baa2173d6a62403d4dc88b7a7f506727/src/diskv/compression.go
func Zlib() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return zlib.NewReader(buffer)
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return zlib.NewWriter(buffer), nil
	}
	return &compression{
		writer:   writer,
		reader:   reader,
		Encoding: "deflate",
	}
}

// https://github.com/gustavo-hms/trama/blob/b7b94a4d7a90475aa2e51db723605dfb70f52dc1/param_decoder.go
func retrieveParams(q interface{}) (string, error) {
	var values = &url.Values{}
	var str = reflect.ValueOf(q)
	var typ = reflect.TypeOf(q)

	switch q.(type) {
	case url.Values:
		return q.(url.Values).Encode(), nil
	default:
		for i := 0; i < str.NumField(); i++ {
			values.Add(strings.ToLower(typ.Field(i).Name), fmt.Sprintf("%v", str.Field(i).Interface()))
		}
		return values.Encode(), nil
	}
}

func sanitizeBody(b interface{}) (io.Reader, error) {
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

var outboundDialer = &net.Dialer{Timeout: 1000 * time.Millisecond}
var outboundTransport = &http.Transport{Dial: outboundDialer.Dial, Proxy: http.ProxyFromEnvironment}
var outboundClient = &http.Client{Transport: outboundTransport}
var proxyTransport *http.Transport
var proxyClient *http.Client

func SetOutboundTimeout(duration time.Duration) {
	outboundDialer.Timeout = duration
}

func (r *Request) AddHeader(name string, value string) {
	if r.headers == nil {
		r.headers = []httpHeader{}
	}
	r.headers = append(r.headers, httpHeader{name: name, value: value})
}

func (r Request) Send() (*Response, error) {
	var req *http.Request
	var er error
	var transport = outboundTransport
	var client = outboundClient

	if r.Proxy != "" {
		proxyUrl, err := url.Parse(r.Proxy)
		if err != nil {
			return nil, &Error{Err: err}
		}
		if proxyTransport == nil {
			proxyTransport = &http.Transport{
				Dial:  outboundDialer.Dial,
				Proxy: http.ProxyURL(proxyUrl),
			}
			proxyClient = &http.Client{
				Transport: proxyTransport,
			}
		} else {
			proxyTransport.Proxy = http.ProxyURL(proxyUrl)
		}
		transport = proxyTransport
		client = proxyClient
	}

	if !r.SSL {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else if transport.TLSClientConfig != nil {
		transport.TLSClientConfig.InsecureSkipVerify = false
	}

	b, e := sanitizeBody(r.Body)
	if e != nil {
		return nil, &Error{Err: e}
	}

	// https://github.com/gustavo-hms/trama/blob/b7b94a4d7a90475aa2e51db723605dfb70f52dc1/param_decoder.go
	if r.QueryString != nil {
		param, e := retrieveParams(r.QueryString)
		if e != nil {
			return nil, &Error{Err: e}
		}
		r.URI = r.URI + "?" + param
	}

	var bodyReader io.Reader

	if b != nil && r.Compression != nil {
		buffer := bytes.NewBuffer([]byte{})
		readBuffer := bufio.NewReader(b)
		writer, err := r.Compression.writer(buffer)
		if err != nil {
			return nil, &Error{Err: err}
		}
		_, e = readBuffer.WriteTo(writer)
		writer.Close()
		if e != nil {
			return nil, &Error{Err: e}
		}
		bodyReader = buffer
	} else {
		bodyReader = b
	}
	req, er = http.NewRequest(r.Method, r.URI, bodyReader)

	if er != nil {
		// we couldn't parse the URL.
		return nil, &Error{Err: er}
	}

	// add headers to the request
	req.Host = r.Host
	req.Header.Add("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", r.ContentType)
	req.Header.Add("Accept", r.Accept)

	if r.Compression != nil {
		req.Header.Add("Content-Encoding", r.Compression.Encoding)
		req.Header.Add("Accept-Encoding", r.Compression.Encoding)
	}

	if r.headers != nil {
		for _, header := range r.headers {
			req.Header.Add(header.name, header.value)
		}
	}

	if r.BasicAuthUser != "" && r.BasicAuthPass != "" {
		req.SetBasicAuth(r.BasicAuthUser, r.BasicAuthPass)
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

	// some fancy error catching
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

	if isActualRedirect(res.StatusCode) && r.MaxRedirects > 0 {
		loc, _ := res.Location()
		r.MaxRedirects--
		r.URI = loc.String()
		return r.Send()
	}

	if r.Compression != nil && strings.Contains(res.Header.Get("Content-Encoding"), r.Compression.Encoding) {
		compressed, err := r.Compression.reader(res.Body)
		if err != nil {
			return nil, &Error{Err: err}
		}
		return &Response{
			StatusCode:    res.StatusCode,
			ContentLength: res.ContentLength,
			Header:        res.Header,
			Body: &Body{
				reader:     res.Body,
				compressed: compressed,
			}}, nil
	} else {
		return &Response{
			StatusCode:    res.StatusCode,
			ContentLength: res.ContentLength,
			Header:        res.Header,
			Body: &Body{
				reader: res.Body,
			}}, nil
	}
}
