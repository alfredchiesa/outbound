# outbound
[![Build Status](https://travis-ci.org/alfredchiesa/outbound.svg?branch=master)](https://travis-ci.org/alfredchiesa/outbound)  
outbound is a simple abstraction of a couple of standard golang packages. It will allow for dirt simple network requests, similar to that of the [Requests](https://github.com/kennethreitz/requests/) library for python.

There are many other network clients available for Go, but none really seemed to fit my work flow. I'm starting this to build something, selfishly. To help my life easier. If you can use it to help make your life easier, please use and PR.

## Road map
Starting off with nothing at the moment and going to see if we can add a ton of stuff. I want to make sure that I get a basic client working first, with a basic structure of how the client feels. Then work on other, more advanced stuff.

Feature | Status
------- |:------:
define client method structure | :blue_book:
GET method | :blue_book:
PUT method | :ledger:
POST method | :ledger:
PATCH method | :ledger:
DELETE method | :ledger:
OPTIONS method | :ledger:
implement logging | :ledger:
documentation | :ledger:
dynamically add headers | :ledger:
keep-alive | :ledger:
basic-auth | :ledger:
json reflection helper | :ledger:
string reflection helper | :ledger:
cookie methods | :ledger:
gzip | :ledger:
timeouts | :ledger:
UDP method | :ledger:
connection pooling | :ledger:

  >**key**  
  >:green_book: = done  
  >:blue_book: = doing  
  >:ledger: = not started  


