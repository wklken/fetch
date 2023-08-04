# httptest

httptest is a **lightweight** and **powerful** API testing tool.

You can use git and configuration files (e.g. yaml/toml) to maintain all test cases.

See [https://wklken.me/httptest](https://wklken.me/httptest) for more information.

## Installation

### Binary releases

See the available binaries for different operating systems/architectures from the [releases page](https://github.com/wklken/httptest/releases).

### go install

```bash
go install github.com/wklken/httptest@latest
```

### hombrew

```bash
brew tap wklken/httptest
brew install httptest
```

### docker

```bash
docker run --rm --net=host wklken/httptest help

# apple m1*/m2*
docker pull --platform=linux/amd64 wklken/httptest
docker run --rm --net=host wklken/httptest help
```

### Build from source

- dependencies: go1.19

```bash
git clone https://github.com/wklken/httptest.git
cd httptest
make build
```

## Getting Started

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

run

```bash
$ ./httptest run http_get.yaml
```
![](./docs/assets/images/getting-started.jpg)

See [examples](https://github.com/wklken/httptest/tree/master/examples) for a variety of examples.

## Features

- config file types: yaml/toml/json/properties/ini
- http methods: get/post/put/delete/patch/head/options
- build request: file/template([go template](https://golang.org/pkg/text/template/))/set cookie/basic auth/share cookie/msgpack
- assert response:
  - status/statusCode/contentLength/contentType/body
  - latency
  - numeric support `_in/_not_in/_lt/_lte/_gt/_gte`
  - string support `_contains/_not_contains/_startswith/_endswith`
  - proto/protoMajor/protoMinor
  - redirect
  - header/cookie/error/json([jmespath](https://jmespath.org/tutorial.html))/html(xpath)/xml(xpath)/yaml/toml/redirect
- cli:
  - progress bar
  - show stats
  - verbose mode: `-v` or set `export HTTPTEST_DEBUG = true`
  - quiet mode: `-q/--quiet` to silent the print, for check `$?` only
  - run in order / run in parallel
  - set failFast/timeout

## Feedback

If you have any feedback, please create an [issue](https://github.com/wklken/httptest/issues)

## License

Copyright (c) 2021-present [wklken](https://github.com/wklken)

Licensed under [MIT License](https://github.com/wklken/httptest/blob/master/LICENSE)
