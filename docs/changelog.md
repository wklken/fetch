---
layout: page
title: CHANGELOG
nav_order: 6
permalink: /changelog/
---

# CHANGELOG

## v1.0.2(2022-12-23)

### Add

- assert: header_exists/cookie_exists/body_icontains/body_matches/body_not_matches

### Documentation

- add more examples and docs

### CICD

- add goreleaser on github action


## v1.0.1(2022-12-18)

### Bug Fixes

- Json assertion type not match / repsonse body length=0 EOF
- Fix ordering wrong, duplicate files in the pattern
- Verbose output log print order

### Documentation

- Refactor examples
- Update docs
- Update readme
- Add more examples

### Features

- Add more examples
- Add http proxy
- Add default user-agent
- Support hasRedirect and redirectCount_{op}


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

