# outbound

outbound is a simple abstraction of a couple of standard golang packages. It will allow for dirt simple network requests, similar to that of the [Requests](https://github.com/kennethreitz/requests/) library for python.

There are many other network clients available for Go, but none really seemed to fit my work flow. I'm starting this to build something, selfishly. To help my life easier. If you can use it to help make your life easier, please use and PR.

## Road map
Starting off with nothing at the moment and going to see if we can add a ton of stuff. I want to make sure that I get a basic client working first, with a basic structure of how the client feels. Then work on other, more advanced stuff.

Version | Feature | Status
:------:| ------- |:------:
1 | define client method structure | :blue_book:
1 | GET method | :blue_book:
2 | PUT method | :ledger:
2 | POST method | :ledger:
2 | PATCH method | :ledger:
2 | DELETE method | :ledger:
2 | OPTIONS method | :ledger:
3 | implement logging | :ledger:
3 | documentation | :ledger:
3 | dynamically add headers | :ledger:
3 | keep-alive | :ledger:
3 | basic-auth | :ledger:
4 | json reflection helper | :ledger:
4 | string reflection helper | :ledger:
4 | cookie methods | :ledger:
5 | gzip | :ledger:
5 | timeouts | :ledger:
6 | UDP method | :ledger:
7 | connection pooling | :ledger:

  >**key**  
  >:green_book: = done  
  >:blue_book: = doing  
  >:ledger: = not started  


