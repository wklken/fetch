---
layout: page
title: Installation
nav_order: 4
permalink: /installation/
---

# Installation
{: .no_toc }

Some examples show how to build reuqest.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---


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

- dependencies: go1.21

```bash
git clone https://github.com/wklken/fetch.git
cd fetch
make build
```
