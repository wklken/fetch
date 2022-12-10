---
layout: default
title: Examples
parent: Usage
permalink: /usage/examples/

---


- `title`/`description`: basic info of this case
- `request`: the http request definition
    - `method`: http method like `get/post/put/patch/delete/head`
    - `url`: the target url
- `assert`: the assertions for the response
    - `status`: should be `ok`
    - `statusCode`: should be `200`
    - `contentLength_gt`: content-length should be greater than `180`
    - `contentType`: should be `application/json`


json assertions

```toml
[request]
method = "post"
url = "http://httpbin.org/post"
body = """
{
    "hello": "world",
    "array": [1, 2, 3, 4]
}
"""
[request.header]
Content-Type = "application/json"


[assert]
status = "ok"
statusCode = 200
contentLength_gt = 180
contentType = "application/json"

[[assert.json]]
path = "headers.Host"
value =  "httpbin.org"

[[assert.json]]
# path = "headers.\"Accept-Encoding\""
path = 'headers."Accept-Encoding"'
value =  "gzip"

[[assert.json]]
path = 'json.array[0]'
value =  1

[[assert.json]]
path = 'json.hello'
value =  "world"

[[assert.json]]
path = '*.hello'
value =  ["world"]

[[assert.json]]
path = "length(json.array)"
value = 4
```

with template render, `./httptest run examples/use_template.toml -c examples/config/dev.toml -v`

```toml
title = "http method post, use template"
description = "http method post"

[request]
method = "post"
url = "{{.host}}/post"
body = """
{
    "hello": "{{.name}}",
    "world": {{if .debug}}"in debug mode"{{else}}"not debug mode"{{end}},
    "array": "{{range $i, $a := .array}} {{$i}}{{$a}} {{end}}"
}
"""
[request.header]
Content-Type = "{{.content_type}}"


[assert]
status = "ok"
statusCode = 200
contentLength_gt = 180
contentType = "{{.content_type}}"
```
