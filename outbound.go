package main

import "net/http"

// Some basic structs for now
type Response struct {
	StatusCode int
	Status     string
	Body       string
}

type Request struct {
	Uri         string
	Method      string
	Body        string
	UserAgent   string
	ContentType string
	Host        string
	Headers     []string
}

// Struct method testing
//
// trying methods on the Request struct. This might be cleaner for when
// it comes time to start passing values. instead of doing kwargs of interfaces
func (req Request) Send(*Response, error) {

	client := &http.Client{}

	r, err := http.NewRequest(req.Method, req.Uri, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(r)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Public Func Prototypes
// these will be like outbound.GET, outbound.PUT, outbound.UDP
func GET(...interface{}) {
	//pass
}

func POST(...interface{}) {
	//pass
}

func PUT(...interface{}) {
	//pass
}

func DELETE(...interface{}) {
	//pass
}

func HEAD(...interface{}) {
	//pass
}

func OPTIONS(...interface{}) {
	//pass
}

func PATCH(...interface{}) {
	//pass
}

func UDP(...interface{}) {
	//pass
}

// main loop to test things. won't be here when live
func main() {
	res, err := Request{
		Method: "GET",
		Uri:    "http://www.villaroad.com",
	}.Send()
}
