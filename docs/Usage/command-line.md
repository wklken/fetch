---
layout: default
title: Command Line
parent: Usage
permalink: /usage/command-line/
nav_order: 7
---

# Command Line
{: .no_toc }

The cli options.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## overview

```bash
$ ./fetch -h
A command lin http test tool. Complete documentation is available at https://github.com/wklken/fetch

Usage:
  fetch [flags]
  fetch [command]

Available Commands:
  bootstrap   A brief description of your command
  completion  Generate the autocompletion script for the specified shell
  generate    A brief description of your command
  help        Help about any command
  run         Run cases
  version     Print the version number

Flags:
  -h, --help   help for fetch

Use "fetch [command] --help" for more information about a command.
```

avaliable commands:
- run
- version
- help
- completion

not ready commands:
- bootstrap
- generate

## run

### verbose `-v`

will show the request/response detail

```bash
$ ./fetch run examples/http_get.yaml -v
 100% |███████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (1/1, 2 it/s)
Run Case: examples/http_get.yaml | example | [GET http://httpbin.org/get] | 466ms
> GET /get HTTP/1.1
> Host: httpbin.org
> User-Agent: fetch/1.0.0
> Accept-Encoding: gzip
>
>
< HTTP/1.1 200 OK
< Content-Length: 272
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Connection: keep-alive
< Content-Type: application/json
< Date: Mon, 12 Dec 2022 15:22:23 GMT
< Server: gunicorn/19.9.0
<
< {
<   "args": {},
<   "headers": {
<     "Accept-Encoding": "gzip",
<     "Host": "httpbin.org",
<     "User-Agent": "fetch/1.0.0",
<     "X-Amzn-Trace-Id": "Root=1-6397472f-07cfba2c0e88067e1abead0c"
<   },
<   "origin": "1.1.1.1",
<   "url": "http://httpbin.org/get"
< }
<
assert.statuscode: Pass
assert.status: Pass
assert.contenttype: Pass
assert.contentlength_gt: Pass

┌─────────────────────────┬─────────────────┬─────────────────┬─────────────────┐
│                         │           total │          passed │          failed │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│                   cases │               1 │               1 │               0 │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│              assertions │               4 │               4 │               0 │
├─────────────────────────┴─────────────────┴─────────────────┴─────────────────┤
│                    Time :    467 ms                                           │
└───────────────────────────────────────────────────────────────────────────────┘
the execute result: 0
```

### quiet `-q`

will no output at all

```bash
$ ./fetch run examples/http_get.yaml -q
$ echo $?
0
```

### parallel `-p`

will run all cases with `N` goroutines, much more faster than run cases one by one

```bash
$ ./fetch run examples/http_*.yaml -p 5
```

### config `-c`

run cases with configs

```bash
$ ./fetch run -c examples/config.yaml examples/request_use_template.yaml
```

the config.yaml file example

```yaml
debug: false
failFast: true
timeout: 2000
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

please see [Usage/Config](/usage/config/) for more detail

### proxy `--proxy`

```
$ ./fetch run examples/http_get.yaml --proxy http://127.0.0.1:9999
```
