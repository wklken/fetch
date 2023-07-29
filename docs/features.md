---
layout: page
title: Features
nav_order: 2
permalink: /features/
---

# Features

## Request

- http methods: `get/post/put/delete/patch/head/options`
- cases config file types: `toml/yaml/json/ini/properties`
- build request:
    - post form
    - send cookie
    - send body via file `body = "@./post_body_file"`
    - basic auth
    - share cookie
    - send and parse `application/msgpack`
    - disable_redirect

## Assert

- proto/protoMajor/protoMinor
- status/statusCode/contentLength/contentType/body
- headers
- cookies
- error: `assert.error_contains` for send fail
- operators:
  - numeric: `_in/_not_in/_lt/_lte/_gt/_gte`
  - string: `_contains/_not_contains/_startswith/_endswith`
- latency
- has redirected or not
- json: use [jmespath](https://jmespath.org/tutorial.html) to get value then assert
- xml/html: use xpath to get value then assert
- yaml/toml: convert to json then use jmespath to get value then assert

## Config file

- use a config file: `./httptest run -c examples/config/dev.toml examples/get.toml`
- config: `debug=true/false` to trigger debug print
- config: `render=true/false` to use  go_template render the `env`
- config: `failFast=true/false`, will exit if got one fail case while running
- config: `timeout=1000`, will set request timeout to 1000ms, fail if exceed

## command line

- show progress bar
- show run result with stats
- `exit code != 0` if got any fail assertions
- show the fail assertion line number in file
- verbose mode: `./httptest run -v examples/get.toml` or `export HTTPTEST_DEBUG = true`
- quiet mode: `-q/--quiet` to silent the print, for check `$?` only
- support run cases in order, `./httptest run -c examples/config/order.toml`
- support run cases in parallel, `./httptest run -c *.yaml -p 10`
- support case retry/repeat




