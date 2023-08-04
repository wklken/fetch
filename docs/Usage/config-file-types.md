---
layout: default
title: Config File Types
parent: Usage
permalink: /usage/config-file-types/
nav_order: 1
---

# Config File Types
{: .no_toc }

You can use any of these file types for case
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

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

if you choose, you can add more than one case in the yaml(note: separated by three dashes `---`)

```yaml
request:
  method: get
  url: 'http://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
---
request:
  method: get
  url: 'http://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```