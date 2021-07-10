# httptest

A command line http test tool. Maintain the cases via git and pure text


## target

- all in text(.toml/.yaml/.json)
- easy to create/modify/copy and delete
- maintained by git
- run fast

## examples

```
[request]
method = "get"
url = "http://httpbin.org"

[assert]
status = "OK"
statusCode = 200
```

```toml
[request]
method = "get"
url = "http://httpbin.org/response-headers?freeform="
[request.header]
hello = "world"

[assert]
status = "ok"
statusCode = 200
statusCode_in = [400, 500]
statusCode_lt = 100
statusCode_lte = 100
statusCode_gt = 500
statusCode_gte = 500
contentLength = 18
contentLength_lt = 1
contentLength_lte = 1
contentLength_gt = 180
contentLength_gte = 180
body = "HTTPBIN is awesome"
body_contains = "awesome2"
body_not_contains = "awesome"
body_startswith = "A"
body_endswith = "a"
contentType = "abc"
latency_lt = 0
latency_lte = 0
latency_gt = 100
latency_gte = 100
```

## assertions

```toml
# status
status = "ok"
statusCode = 200
statusCode_in = [400, 500]
statusCode_lt = 100
statusCode_lte = 100
statusCode_gt = 500
statusCode_gte = 500

# content-length
contentLength = 18
contentLength_lt = 1
contentLength_lte = 1
contentLength_gt = 180
contentLength_gte = 180

# body
body = "HTTPBIN is awesome"
body_contains = "awesome2"
body_not_contains = "awesome"
body_startswith = "A"
body_endswith = "a"
```

## inspired by

- testify/assert https://github.com/stretchr/testify/tree/master/assert (use this module, and copied some un-exported codes from it, follow the license)
- postman & newman https://www.npmjs.com/package/newman

----------------------------------------

## packages

- http client: https://golang.org/src/net/http/request.go
- config file and cases: toml? https://github.com/pelletier/go-toml
- assert: use testify data compare? https://github.com/stretchr/testify/blob/master/assert/assertions.go

- config file and cases: toml? https://github.com/pelletier/go-toml => use viper to support most config file types

## how it works


## TODO


High: DO POST, the parse the json, do assert
- https://github.com/tidwall/gjson
- https://github.com/oliveagle/jsonpath

High:
- timeout: 5

High:
- latency assert

- [x] init project
- [x] the case name? where to put that?
- [x] assert http status/content-length/content-type
- [x] assert latency

- [ ] `-h/--help`
- [ ] `-v` verbose, simple => debug the request and response
- [ ] `-vv` verbose, detail. file/case? title/description/assert lint/why fail
- [ ] assert json
- [ ] invalid assert or not used assert
- [ ] post json in toml
- [ ] should support all request method
- [ ] should support all request body, json/form/msgpack/zip.....
- [ ] `bootstrap` create the raw template, like `a.hp`
- [ ] `generate x` generate a case
- [ ] `run` run all cases
- [ ] `run` specific file / dir
- [ ] support config file, like `prod.yaml`/`test.yaml`/`dev.yaml`, `-e prod.yaml`
- [ ] support environment vars, like `host/basic auth`,
- [ ] render environment vars in everywhere, like `path/request section/assert section`? which template to use?
- [ ] how to control the execute order?
- [ ] multiple cases in one file, like ginkgo?
- [ ] how to: long-live / file download / static file
- [ ] support retry
- [ ] support latency assertion, less than/greater than, or between
- [ ] support repeat, like run 5 times
- [ ] support assert redirect
- [ ] how to share the cookie between cases? claim? or default same dir
- [ ] run in parallel
- [ ] output stats
- [ ] dns / connection reset/timeout and so on
- [ ] case set some data, next case read it
- [ ] support ssl
- [ ] case with line number
- [ ] keep alive
- [ ] websocket


## go mod

command line:
- github.com/spf13/cobra
- github.com/spf13/viper
- github.com/fatih/color

unittest:
- github.com/onsi/ginkgo
- github.com/onsi/gomega
- github.com/stretchr/testify

- gin? https://github.com/gin-gonic/gin/blob/master/context.go

