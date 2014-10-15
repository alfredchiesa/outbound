# outbound
[![Build Status](https://travis-ci.org/alfredchiesa/outbound.svg?branch=master)](https://travis-ci.org/alfredchiesa/outbound)  
outbound is a simple abstraction of a couple of standard golang packages. It will allow for dirt simple network requests, similar to that of the [Requests](https://github.com/kennethreitz/requests/) library for python.

There are many other network clients available for Go, but none really seemed to fit my work flow. I'm starting this to build something, selfishly. To help my life easier. If you can use it to help make your life easier, please use and PR.

##Dependencies
There is only one dependency at the moment, **websocket**. I hear rumor that it will make it's way into stdlib though. Although, I wouldn't hold your breath for that.

```bash
go get code.google.com/p/go.net/websocket
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
connection pooling | :ledger:

  >**key**  
  >:green_book: = done  
  >:blue_book: = doing  
  >:ledger: = not started  


