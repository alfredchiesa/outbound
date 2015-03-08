![outbound github banner](http://i.imgur.com/Dl7vPVe.png)
outbound is a multi-function http client for Golang. Yet another, if you will. It borrows a lot of *best* practices from various other packages and clients. It combines their paradigms into a easy to use, but full featured, outbound http client.

Much of the inspiration (*and slogan*) is derived from the Python [Requests](https://github.com/kennethreitz/requests/) package. It abides by parts of the 20th Python Enchancement Proposal ([PEP 20](https://www.python.org/dev/peps/pep-0020/)), which is also known as the Zen of Python:
```go
1. Beautiful is better than ugly.
2. Explicit is better than implicit.
3. Simple is better than complex.
4. Complex is better than complicated.
5. Readability counts.
```

This package is a continual *work in progress*, and is considered stable enough for use. I see no major API changes coming to the **Request** type any time soon. However, the WebSocket and UDP types are not completely fleshed out as of `March 7, 2015`.

**Table of Contents**

- [outbound](#outbound)
  - [Status](#status)
  - [Dependencies](#dependencies)
  - [Installation](#installation)
  - [Usage](#usage)
    - [REST Verbs](#rest-verbs)
      - [GET](#get)
      - [PUT](#put)
      - [POST](#post)
      - [PATCH](#patch)
      - [DELETE](#delete)
      - [OPTIONS](#options)
    - [Basic Auth](#basic-auth)
    - [Reflection Helper Methods](#reflection-helper-methods)
      - [Response Body JSON to Struct](#response-body-json-to-struct)
      - [Response Body to String](#response-body-to-string)
    - [Custom Headers](#custom-headers)
    - [Websocket](#websocket)
  - [Road map](#road-map)

##Status
[![Build Status](https://travis-ci.org/alfredchiesa/outbound.svg?branch=master)](https://travis-ci.org/alfredchiesa/outbound)  
outbound is currently built on the free continuous integration stack/site known as [Travis](https://travis-ci.org). To view the build history you can either click the button above or [this link](https://travis-ci.org/alfredchiesa/outbound).

##Dependencies
At the moment, there is only one dependency; [websocket](https://godoc.org/golang.org/x/net/websocket). I hear rumor that it will make it's way into stdlib. To install *websocket*, you can use your current deps manager or just runt the following:

    $ go get golang.org/x/net/websocket

##Installation
    $ go get github.com/alfredchiesa/outbound

##Usage
###REST Verbs
####GET
If you just want to fire off a simple GET request:
```go
res, err := outbound.Request{URI: "http://api.someserver.com/user"}.Send()
```

####PUT
In this example you can see how to do an incremental PUT update:
```go
type User struct {
    Age int
}

user := User{Age: 62}

res, err := outbound.Request{
    Method:      "PUT",
    Accept:      "application/json",
    ContentType: "application/json",
    Uri:         "http://api.someserver.com/user/carl/age",
    Body:        user,
}.Send()
```

####POST
An example adding an entirely new resource with POST:
```go
type User struct {
    Name string
    City string
    Age  int
}

user := User{Name: "Carl Sagan", City: "Seattle", Age: 62}

res, err := outbound.Request{
    Method:      "POST",
    Accept:      "application/json",
    ContentType: "application/json",
    Uri:         "http://api.someserver.com/user",
    Body:        user,
}.Send()
```

####PATCH
```go
type Patches struct {
    Operations []Patch
}

type Patch struct {
    Op    string `json:"op"`
    Path  string `json:"path"`
    Value string `json:"value"`
}

func (ps *Patches) AddOp(pt Patch) {
    ps.Operations = append(ps.Operations, pt)
    return
}

func main() {
    operations := []Patch{}
    payload := Patches{operations}

    payload.AddOp(Patch{
        Op:    "test",
        Path:  "/user/carl/age",
        Value: "88",
    })

    payload.AddOp(Patch{
        Op:    "replace",
        Path:  "/user/tina/age",
        Value: "73",
    })

    payload.AddOp(Patch{
        Op:    "delete",
        Path:  "/user",
        Value: "carl",
    })

    res, err := outbound.Request{
        Method:      "PATCH",
        Accept:      "application/json",
        ContentType: "application/json",
        Uri:         "http://api.someserver.com/user",
        Body:        payload,
    }.Send()
}
```

####DELETE
```go
res, err := outbound.Request{
    Method: "DELETE",
    URI:    "http://api.someserver.com/user/sally",
}.Send()
```

####OPTIONS
```go
res, err := outbound.Request{
    Method: "OPTIONS",
    Accept: "httpd/unix-directory",
    URI:    "http://api.someserver.com/user",
}.Send()
```

###Basic Auth
```go
client := outbound.Request{
    Method:        "GET",
    Uri:           "http://api.someserver.com/protected/resource",
    BasicAuthUser: "carl",
    BasicAuthPass: "sagan",
}

res, err := client.Send()
```

###Reflection Helper Methods
####Response Body JSON to Struct
```go
type User struct {
    Name string
    City string
    Age  int
}

res, err := outbound.Request{URI: "http://api.someserver.com/user/carl"}.Send()

res.Body.JsonToStruct(&User)
```

####Response Body to String
```go
res, err := outbound.Request{URI: "http://api.someserver.com/user/carl"}.Send()

fmt.Println(res.Body.ToString())
```

###Custom Headers
```go
client := outbound.Request{
    Method:        "GET",
    Uri:           "http://api.someserver.com/protected/resource",
}

client.AddHeader("X-CUSTOM-HEADER", "billions of billions")
client.AddHeader("X-ANOTHER-HEADER", "trillions of trillions")

res, err := client.Send()
```

###Unix Socket
This is on my todo list. Seems like using Unix sockets, if you are on a system that can, is the better and more reliable approach here. I plan to implement them and offer the option. More to come soon.

###Websocket
The current implementation is littered with bugs that I haven't fixed yet, so it's currently *Out of Order*. I will come back to fix this soon, just not now.
```go
ws, err := outbound.WebSocket{
    Server: "127.0.0.1:443",
}.Conn()

ws.Write("stuff")

ws.Close()
```

## Road map
Starting off with nothing at the moment and going to see if we can add a ton of stuff. I want to make sure that I get a basic client working first, with a basic structure of how the client feels. Then work on other, more advanced stuff.

Feature | Status
------- |:------:
define client method structure | :green_book:
GET method | :green_book:
PUT method | :green_book:
POST method | :green_book:
PATCH method | :blue_book:
DELETE method | :green_book:
OPTIONS method | :blue_book:
implement logging | :blue_book:
documentation | :ledger:
dynamically add headers | :green_book:
keep-alive | :ledger:
basic-auth | :green_book:
json reflection helper | :green_book:
string reflection helper | :green_book:
cookie methods | :blue_book:
gzip | :green_book:
timeouts | :green_book:
UDP method | :ledger:
Web Socket | :blue_book:
connection pooling | :ledger:

  >**key**  
  >:green_book: = done  
  >:blue_book: = doing  
  >:ledger: = not started  
