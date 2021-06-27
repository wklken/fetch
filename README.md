# httptest
A command line http test tool. Maintain the case via git and pure text


## how it works

## rule

- all in text
- easy to create/change/run
- maintained by git
- run quick
- 

## example

simplest

```
[request]
method: get
url: /
[assert]
status: 200
```

```
# a.case
[request]
method: post
url: /api/v1
timeout: 5
body: ??? josn or form

will got response, then do the assert
[assert]
status: 200
type: json
body.code = 0
body.data.length = 10
latency: < 20ms
```

## TODO

- [x] init project
- [ ] `-h/--help`
- [ ] `bootstrap` create the raw template, like `a.hp`
- [ ] `generate x` generate a case
- [ ] `run` run all cases
- [ ] `run` specific file / dir
- [ ] `-v` verbose, simple
- [ ] support config file, like `prod.yaml`/`test.yaml`/`dev.yaml`, `-e prod.yaml`
- [ ] support environment vars, like `host/basic auth`, 
- [ ] render environment vars in everywhere, like `path/request section/assert section`? which template to use?
- [ ] `-vv` verbose, detail. file/case? title/description/assert lint/why fail
- [ ] the case name? where to put that?
- [ ] how to control the execute order?
- [ ] multiple cases in one file, like ginkgo?
- [ ] should support all request method
- [ ] should support all request body, json/form/msgpack/zip.....
- [ ] how to: long-live / file download / static file
- [ ] support retry
- [ ] support repeat, like run 5 times
- [ ] support latency assertion, less than/greater than, or between
- [ ] support assert redirect
- [ ] how to share the cookie between cases? claim? or default same dir
- [ ] run in parallel
- [ ] output stats
- [ ] dns / connection reset/timeout and so on
- [ ] case set some data, next case read it
