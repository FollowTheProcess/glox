# glox

[![License](https://img.shields.io/github/license/FollowTheProcess/glox)](https://github.com/FollowTheProcess/glox)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/glox)](https://goreportcard.com/report/github.com/FollowTheProcess/glox)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/glox?logo=github&sort=semver)](https://github.com/FollowTheProcess/glox)
[![CI](https://github.com/FollowTheProcess/glox/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/glox/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/glox/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/glox)

An implementation of the Lox language from [Crafting Interpreters], written in Go

> [!WARNING]
> **glox is in early development and is not yet ready for use**

![caution](./img/caution.png)

## Project Description

`glox` is a Go implementation of the Lox programming language. I'm translating the Java (and then C) code to idiomatic Go on the fly so it might look a bit different to the book!

On top of that I'll likely make some fun tweaks to the functionality too, I'll try and keep track of that here

## Installation

```shell
brew install FollowTheProcess/tap/glox
```

## Quickstart

```shell
# Run a file
glox run file.lox

# Launch an interactive REPL
glox repl
```

[Crafting Interpreters]: https://craftinginterpreters.com
