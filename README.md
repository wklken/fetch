# fetch

fetch is a **lightweight** and **powerful** API testing tool.

You can use git and yaml files to maintain all test cases.

See [https://wklken.me/fetch](https://wklken.me/fetch) for more information.

## Installation

### Binary releases

See the available binaries for different operating systems/architectures from the [releases page](https://github.com/wklken/fetch/releases).

### go install

```bash
go install github.com/wklken/fetch@latest
```

### hombrew

```bash
brew tap wklken/fetch
brew install fetch
```

### docker

```bash
docker run --rm --net=host wklken/fetch help

# apple m1*/m2*
docker pull --platform=linux/amd64 wklken/fetch
docker run --rm --net=host wklken/fetch help
```

### Build from source

- dependencies: go1.19

```bash
git clone https://github.com/wklken/fetch.git
cd fetch
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
$ ./fetch run http_get.yaml
```
![](./docs/assets/images/getting-started.jpg)

See [examples](https://github.com/wklken/fetch/tree/master/examples) for a variety of examples.

## Features

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
  - verbose mode: `-v` or set `export FETCH_DEBUG = true`
  - quiet mode: `-q/--quiet` to silent the print, for check `$?` only
  - run in order / run in parallel
  - set failFast/timeout

## Feedback

If you have any feedback, please create an [issue](https://github.com/wklken/fetch/issues)

## License

Copyright (c) 2021-present [wklken](https://github.com/wklken)

Licensed under [MIT License](https://github.com/wklken/fetch/blob/master/LICENSE)
