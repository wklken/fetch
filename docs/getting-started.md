---
layout: page
title: Getting Started
nav_order: 3
permalink: /getting-started/
---

# Getting Started

just create a file: `http_get.yaml`

```yaml
title: example
description: a simple example
request:
  method: get
  url: 'http://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

then run the `httptest`

```bash
$ ./httptest run http_get.yaml
```

![](assets/images/getting-started.jpg)

the `title` and `description` are not required

```yaml
request:
  method: get
  url: 'http://httpbin.org/get'
assert:
  status: ok
  statusCode: 200
  contentLength_gt: 180
  contentType: application/json
```

also, you can add more than one case in the yaml(note: separated by three dashes `---`)

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

If you want to know more about how to use httptest, see [Usage](/httptest/usage/)
