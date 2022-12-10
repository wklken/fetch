---
layout: page
title: Features
nav_order: 2
permalink: /features/
---



- **Define:**
    - define the case via [toml](https://toml.io/en/) / yaml / json / properties / ini
    - support http methods: get/post/put/delete/patch/head/options
    - support post form, [examples/form.toml](./examples/form.toml)
    - sent request body via external file `body = "@./post_body_file"` [examples/post_with_body_file.toml](./examples/post_with_body_file.toml)
    - support [go template](https://golang.org/pkg/text/template/) render in all string value, the envs in config file, example: `./httptest run examples/use_template.toml -c examples/config/dev.toml -v`
    - support send cookie [examples/cookies.toml](./examples/cookies.toml)
    - support basic auth [examples/basic_auth.toml](./examples/basic_auth.toml)
    - support share cookie [examples/share_cookies_save.toml](./examples/share_cookies_save.toml) and [examples/share_cookies_use.toml](./examples/share_cookies_use.toml)

- **Assert:**
    - assert status/statusCode/contentLength/contentType/body
    - assert latency
    - assert numeric support `_in/_not_in/_lt/_lte/_gt/_gte`
    - assert string support `_contains/_not_contains/_startswith/_endswith`
    - assert response json body, the path syntax is [jmespath](https://jmespath.org/tutorial.html) [examples/json.toml](./examples/json.toml)
    - assert response headers [examples/header.toml](./examples/header.toml)
    - assert response has redirected [examples/redirect.toml](./examples/redirect.toml)
    - assert response proto/protoMajor/protoMinor

- **Cli and Config:**
    - `exit code != 0` if got any fail assertions
    - verbose mode: `./httptest run -v examples/get.toml` or `export HTTPTEST_DEBUG = true`
    - quiet mode: `-q/--quiet` to silent the print, for check `$?` only
    - show run result with stats
    - configfile: `./httptest run -c examples/config/dev.toml examples/get.toml`
    - configfile: `debug=true/false` to trigger debug print
    - configfile: `render=true/false` to use  go_template render the `env`
    - configfile: `failFast=true/false`, will exit if got one fail case while running
    - configfile: `timeout=1000`, will set request timeout to 1000ms, fail if exceed
    - support run cases in order, `./httptest run -c examples/config/order.toml`