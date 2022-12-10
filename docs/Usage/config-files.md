---
layout: default
title: Config files
parent: Usage
permalink: /usage/config-files/
nav_order: 1
---

## TLDR

You can define the cases in many file types
- YAML (recommend)
- TOML (recommend)
- JSON
- ini
- prop

and run with

```bash
$ ./httptest run a.yaml b.toml c.json d.ini e.prop
```

## Examples

### example: yaml


```yaml
title: asserts
description: all supported assertions
request:
  method: get
  url: 'http://httpbin.org/get'
  header:
    hello: world
assert:
  status: ok
  statusCode: 200
  statusCode_in:
    - 400
    - 500
  statusCode_not_in:
    - 200
    - 400
  statusCode_lt: 100
  statusCode_lte: 100
  statusCode_gt: 500
  statusCode_gte: 500
  contentLength: 18
  contentLength_lt: 1
  contentLength_lte: 1
  contentLength_gt: 180
  contentLength_gte: 180
  contentType: abc
  body: HTTPBIN is awesome
  body_contains: awesome2
  body_not_contains: awesome
  body_startswith: A
  body_endswith: a
  body_not_startswith: '{'
  body_not_endswith: '}'
  latency_lt: 0
  latency_lte: 0
  latency_gt: 100
  latency_gte: 100
```

### example: toml

```
title = "asserts"
description = "all supported assertions"

[request]
method = "get"
url = "http://httpbin.org/get"
[request.header]
hello = "world"


[assert]
# status
status = "ok"
status_in = ["bad request", "geteway timeout"]
status_not_in = ["ok", "geteway timeout"]

statusCode = 200
statusCode_in = [400, 500]
statusCode_not_in = [200, 400]
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

# content-type
contentType = "abc"
contentType_in = ["abc"]
contentType_not_in = ["abc"]

# body
body = "HTTPBIN is awesome"
body_contains = "awesome2"
body_not_contains = "awesome"
body_startswith = "A"
body_endswith = "a"
body_not_startswith = "{"
body_not_endswith = "}"

# latency
latency_lt = 0
latency_lte = 0
latency_gt = 100
latency_gte = 100

# proto
proto = "HTTP/2.0"
protoMajor = 2
protoMinor = 0
```

### example: json

```json
{
  "title": "asserts",
  "description": "all supported assertions",
  "request": {
    "method": "get",
    "url": "http://httpbin.org/get",
    "header": {
      "hello": "world"
    }
  },
  "assert": {
    "status": "ok",
    "statusCode": 200,
    "statusCode_in": [
      400,
      500
    ],
    "statusCode_not_in": [
      200,
      400
    ],
    "statusCode_lt": 100,
    "statusCode_lte": 100,
    "statusCode_gt": 500,
    "statusCode_gte": 500,
    "contentLength": 18,
    "contentLength_lt": 1,
    "contentLength_lte": 1,
    "contentLength_gt": 180,
    "contentLength_gte": 180,
    "contentType": "abc",
    "body": "HTTPBIN is awesome",
    "body_contains": "awesome2",
    "body_not_contains": "awesome",
    "body_startswith": "A",
    "body_endswith": "a",
    "body_not_startswith": "{",
    "body_not_endswith": "}",
    "latency_lt": 0,
    "latency_lte": 0,
    "latency_gt": 100,
    "latency_gte": 100
  }
}
```

### example: ini

```ini
title=asserts
description=all supported assertions

[request]
method=get
url=http://httpbin.org/get

[request.header]
hello=world

[assert]
status=ok
statusCode=200
statusCode_in=400,500
statusCode_not_in=200,400
statusCode_lt=100
statusCode_lte=100
statusCode_gt=500
statusCode_gte=500
contentLength=18
contentLength_lt=1
contentLength_lte=1
contentLength_gt=180
contentLength_gte=180
contentType=abc
body=HTTPBIN is awesome
body_contains=awesome2
body_not_contains=awesome
body_startswith=A
body_endswith=a
body_not_startswith={
body_not_endswith=}
latency_lt=0
latency_lte=0
latency_gt=100
latency_gte=100
```

### example: prop

```
title=asserts
description=all supported assertions
request.method=get
request.url=http://httpbin.org/get
request.header.hello=world
assert.status=ok
assert.statusCode=200
assert.statusCode_lt=100
assert.statusCode_lte=100
assert.statusCode_gt=500
assert.statusCode_gte=500
assert.statusCode_in=400,500
assert.statusCode_not_in=200,400
assert.contentLength=18
assert.contentLength_lt=1
assert.contentLength_lte=1
assert.contentLength_gt=180
assert.contentLength_gte=180
assert.contentType=abc
assert.body=HTTPBIN is awesome
assert.body_contains=awesome2
assert.body_not_contains=awesome
assert.body_startswith=A
assert.body_endswith=a
assert.body_not_startswith={
assert.body_not_endswith=}
assert.latency_lt=0
assert.latency_lte=0
assert.latency_gt=100
assert.latency_gte=100
```