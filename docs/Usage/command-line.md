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
$ ./httptest -h                           
A command lin http test tool. Complete documentation is available at https://github.com/wklken/httptest

Usage:
  httptest [flags]
  httptest [command]

Available Commands:
  bootstrap   A brief description of your command
  completion  Generate the autocompletion script for the specified shell
  generate    A brief description of your command
  help        Help about any command
  run         Run cases
  version     Print the version number

Flags:
  -h, --help   help for httptest

Use "httptest [command] --help" for more information about a command.
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
$ ./httptest run examples/http_get.yaml -v
DEBUG request:
> GET /get HTTP/1.1
> ......
```

### quiet `-q`

will no output at all

```bash
$ ./httptest run examples/http_get.yaml -q
$ echo $?
0
```

### parallel `-p`

will run all cases with `N` goroutines, much more faster than run cases one by one

```bash
$ ./httptest run examples/http_*.yaml -p 5
```

### config `-c`

run cases with configs

```bash
$ ./httptest run -c examples/config.yaml examples/request_use_template.yaml
```

the config.yaml file example

```yaml
debug: false
render: true
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
