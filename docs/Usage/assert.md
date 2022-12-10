---
layout: default
title: Assert
parent: Usage
permalink: /usage/assert/
nav_order: 4
---


## assert header

```yaml
title: http method get headers
description: http method get headers
request:
  method: get
  url: 'http://httpbin.org/response-headers?freeform=hello'
assert:
  status: ok
  statusCode: 200
  contentType: application/json
  header:
    server: gunicorn/19.9.0
    freeform: hello
```

## assert cookie


```yaml
title: cookie assert
description: cookie assert
request:
  method: get
  url: 'http://httpbin.org/cookies/set?name1=value1&name2=value2'
  disable_redirect: true
assert:
  statusCode: 302
  contentType: text/html
  cookie:
    - name: name1
      value: value1
      path: /
    - name: name2
      value: value3
      path: /
```

## assert error

the fail error contains specific text

```yaml
title: http method get
description: http method get
request:
  method: get
  url: 'http://not.exists.local:8113/ping/'
config:
  timeout: 1
assert:
  error_contains: context deadline exceeded
```

## assert json

```yaml
title: json assert
description: json assert
request:
  method: post
  url: 'http://httpbin.org/post'
  body: |
    {
        "hello": "world",
        "array": [1, 2, 3, 4]
    }
  header:
    Content-Type: application/json
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
  json:
    - path: headers.Host
      value: httpbin.org
    - path: headers."Accept-Encoding"
      value: gzip
    - path: 'json.array[0]'
      value: 1
    - path: json.hello
      value: world
    - path: '*.hello'
      value:
        - world
    - path: length(json.array)
      value: 4
    - path: json.abcdefg
      value: 1
```

## assert html

```yaml
title: html assert
description: html assert
request:
  method: get
  url: 'http://httpbin.org/html'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: text/html
  html:
    - path: /html/body/h1
      value: Herman Melville - Moby-Dick
    - path: /slideshow/@author
      value: Overview
```

## assert xml

```yaml
title: xml assert
description: xml assert
request:
  method: get
  url: 'http://httpbin.org/xml'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/xml
  xml:
    - path: '/slideshow/slide[2]/title'
      value: Overview
    - path: /slideshow/@author
      value: Overview
```

## assert yaml or toml

use jmespath

```yaml
title: yaml assert
description: yaml assert
request:
  method: get
  url: 'http://0.0.0.0:8080/someYAML'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 1
  contentType: application/x-yaml
  yaml:
    - path: message
      value: hey
    - path: foo.bar
      value: hello2
```

```yaml
title: toml assert
description: toml assert
request:
  method: get
  url: 'http://0.0.0.0:8080/someTOML'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 1
  contentType: application/toml
  toml:
    - path: message
      value: hey
    - path: foo.bar
      value: hello2
```

## assert redirect

```yaml
title: http redirect
description: http redirect
request:
  method: get
  url: 'http://www.httpbin.org/get'
assert:
  hasRedirect: false
```