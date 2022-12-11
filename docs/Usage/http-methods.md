---
layout: default
title: HTTP Methods
parent: Usage
permalink: /usage/http-methods/
nav_order: 2
---

# HTTP Methods
{: .no_toc }

The HTTP methods supported
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## TLDR

Support http methods:
- get
- post
- put
- delete
- patch
- head
- options

## Examples

### get


```yaml
title: http method get
description: http method get
request:
  method: get
  url: 'http://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

### post

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

### put

```yaml
title: http method put
description: http method put
request:
  method: put
  url: 'http://httpbin.org/put'
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


### delete

```yaml
title: http method delete
description: http method delete
request:
  method: delete
  url: 'http://httpbin.org/delete'
  header:
    Content-Type: application/json
    Accept: application/json
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

### patch

```yaml
title: http method patch
description: http method patch
request:
  method: patch
  url: 'http://httpbin.org/patch'
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

### head

```yaml
title: http method head
description: http method head
request:
  method: head
  url: 'http://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

### options

```yaml
title: http method options
description: http method options
request:
  method: options
  url: 'http://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
  contentLength: 0
  contentType: text/html
```