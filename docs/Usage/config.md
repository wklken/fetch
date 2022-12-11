---
layout: default
title: Config
parent: Usage
permalink: /usage/config/
nav_order: 5
---


# Config

{: .no_toc }

The config file explaination
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## overview

the config file example

```yaml
debug: false

failFast: true
timeout: 2000

order:
  - pattern: examples/post.toml
  - pattern: examples/get.toml
  - pattern: examples/asserts.*

render: true
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

- `debug`: equals to `./httptest run -v`, will show request/response detail. (default `false`)
- `timeout`: `ms`, the timeout for each case (default `0`, means `no timeout`)
- `failFast`: if true, will exit if got fail(default `false`)
- `order`: run cases in order

