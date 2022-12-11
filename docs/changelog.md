---
layout: page
title: CHANGELOG
nav_order: 5
permalink: /changelog/
---

# CHANGELOG

## v1.0.0 (2022-12-09)

### config

- support retry
- support run in parallel
- support case timeout
- support send `application/msgpack`
- support request support disable_redirect

### assert

- assert cookie
- assert error: `assert.error_contains` for http fail
- assert xml: xpath
- assert html: xpath
- assert yaml: jmespath, convert to json and assert as json
- assert toml: jmespath, convert to json and assert as json
- assert msgpack: decode msgpack and assert as json

### output

- show progress bar
- show the line number if fail


## v0.0.1 (2021-08-02)

First release!

Implement all the basic api test assertions

Please read the README.md

