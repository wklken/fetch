---
layout: default
title: Build Request
parent: Usage
permalink: /usage/build-request/
nav_order: 3
---

# Build Request
{: .no_toc }

Some examples show how to build reuqest.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## set headers

```yaml
title: set header
description: set header
request:
  method: get
  url: 'http://httpbin.org/get'
  header:
    User-Agent: "test-client"

assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```


## post form

```yaml
title: http method post form
description: http method post form
request:
  method: post
  url: 'http://httpbin.org/post'
  body: 'param1=value1&param2=value2'
  header:
    Content-Type: application/x-www-form-urlencoded
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

## post with body in file

normally, we can put the body in config file

```yaml
title: http method post
description: http method post
request:
  method: post
  url: 'http://httpbin.org/post'
  body: |
    {
        "hello": "world"
    }
  header:
    Content-Type: application/json
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

or, we can post with a file by `@`

```yaml
title: http method post with body file
description: http method post with body file
request:
  method: post
  url: 'http://httpbin.org/post'
  body: '@./post_body_file'
  header:
    Content-Type: application/json
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```


## set cookie

```yaml
title: http method set cookies
description: http method set cookies
request:
  method: get
  url: 'http://httpbin.org/cookies'
  cookie: uid=123
  header:
    Content-Type: application/json
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 0
  contentType: application/json
  json:
    - path: cookies.uid
      value: '123'
```

## save a cookie then use it in another case

use `hook.save_cookie` to save cookie into file

```yaml
title: http method share cookies
description: http method share cookies
request:
  method: get
  url: 'http://httpbin.org/cookies/set?name1=value1&name2=value2'
  header:
    Content-Type: application/json
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 0
  contentType: application/json
  json:
    - path: cookies.name1
      value: value1
hook:
  save_cookie: share_cookies.txt
```

then use it via `@` in another case

```yaml
title: http method share cookies
description: http method share cookies
request:
  method: get
  url: 'http://httpbin.org/cookies'
  cookie: '@share_cookies.txt'
  header:
    Content-Type: application/json
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 0
  contentType: application/json
  json:
    - path: cookies.name1
      value: value1
```

## basic auth

```yaml
title: http method basic auth
description: http method basic auth
request:
  method: get
  url: 'http://httpbin.org/basic-auth/hello/world'
  header:
    Content-Type: application/json
    # equals to
    # Authorization: Basic aGVsbG86d29ybGQ=
  basic_auth:
    username: hello
    password: world
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 0
  contentType: application/json
  json:
    - path: authenticated
      value: true
    - path: user
      value: hello
```



## send msgpack

will auto encode the json with msgpack then send

```yaml
title: http method post msgpack
description: http method post msgpack
request:
  method: post
  url: 'http://127.0.0.1:8080/anything'
  body: |
    {
        "hello": "world",
        "a": "1",
        "b": "2",
        "c": "1",
        "foo": "bar"
    }
  header:
    Content-Type: application/msgpack
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

## use tempalte

you can create an `config.yaml` with a lot of envs (`examples/config.yaml`)

```yaml
env:
  hello: world
  name: tom
  host: 'http://httpbin.org'
  content_type: application/json
  array:
    - a
    - b
    - c
```

then use [go template](https://golang.org/pkg/text/template/) in the case (`examples/request_use_template.yaml`)

```yaml
title: http method post, use template
description: http method post use template
request:
  method: post
  url: '{{.host}}/post'
  body: |
    {
        "hello": "{{.name}}",
        "world": {{if .debug}}"in debug mode"{{else}}"not debug mode"{{end}},
        "array": "{{range $i, $a := .array}} {{$i}}{{$a}} {{end}}"
    }
  header:
    Content-Type: '{{.content_type}}'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: '{{.content_type}}'
```

run with command

```bash
$ ./fetch run -c examples/config.yaml examples/request_use_template.yaml

# will send
DEBUG request:
> POST /post HTTP/1.1
> Host: httpbin.org
> User-Agent: fetch/1.0.0
> Content-Length: 83
> Content-Type: application/json
> Accept-Encoding: gzip
>
> {
>     "hello": "tom",
>     "world": "not debug mode",
>     "array": " 0a  1b  2c "
> }
>
```

also, you can add `env` in case, which's priority is higher than `env` in config file(`./fetch -c config.yaml`) (`examples/request_use_template_local.yaml`)

```yaml
title: 'http method post, use template local'
description: http method post, use template local
request:
  method: post
  url: '{{.host}}/post'
  body: |
    {
        "hello": "{{.name}}",
        "world": {{if .debug}}"in debug mode"{{else}}"not debug mode"{{end}},
        "array": "{{range $i, $a := .array}} {{$i}}{{$a}} {{end}}"
    }
  header:
    Content-Type: '{{.content_type}}'
env:
  name: jerry
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: '{{.content_type}}'
```

run with command

```bash
$ ./fetch run -c examples/config.yaml examples/request_use_template_local.yaml

# will send
DEBUG request:
> POST /post HTTP/1.1
> Host: httpbin.org
> User-Agent: fetch/1.0.0
> Content-Length: 85
> Content-Type: application/json
> Accept-Encoding: gzip
>
> {
>     "hello": "jerry",
>     "world": "not debug mode",
>     "array": " 0a  1b  2c "
> }
>
```

## chain the requests

you can chain the requests, the response of the previous request will be parsed and saved in the context, then you can use it in the next request

```yaml

```yaml
request:
  method: get
  url: 'https://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
parse:
  - key: origin
    source: body
    jmespath: "origin"
  - key: length
    source: header
    header: Content-Length
---
request:
  method: get
  url: 'https://httpbin.org/get'
  header:
    X-GOT-ORIGIN: '{{.origin}}'
assert:
  statusCode: 302
```

## config: set a timeout for case

```yaml
title: http timeout
description: http timeout
request:
  method: get
  url: 'http://httpbin.org/get'
config:
  timeout: 1
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 10
  contentType: application/json
```

the priority of `timeout` in case is higher than `timeout` in config file(`./fetch -c config.yaml`)

## config: set repeat for case

```yaml
title: http method get repeat
description: http method get repeat
request:
  method: get
  url: 'http://httpbin.org/get'
config:
  repeat: 3
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 10
  contentType: application/json
```

## config: set retry for case

if statusCode not equals to 200, will retry for 3 times, and the interval is 1000ms

if equals to 200, will do assertions, or use the last response to do assertions.

```yaml
title: http method get retry
description: http method get retry
request:
  method: get
  url: 'http://httpbin.org/status/500'
config:
  retry:
    enable: true
    count: 3
    interval: 1000
    statusCodes:
      - 200
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 10
  contentType: application/json
```
