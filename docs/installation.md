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

See the available binaries for different operating systems/architectures from the [releases page](https://github.com/wklken/httptest/releases).

### go install

```bash
go install github.com/wklken/httptest@latest
```

### Build from source

- dependencies: go1.19

```bash
git clone https://github.com/wklken/httptest.git
cd httptest
make build
```