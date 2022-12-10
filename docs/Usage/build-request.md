---
layout: default
title: Build Request
parent: Usage
permalink: /usage/build-request/
nav_order: 3
---


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