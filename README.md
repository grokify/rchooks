# RCHooks - RingCentral Webhook CLI Management App

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

## Usage

```
$ rchooks --list
$ rchooks --create=https://example.com/webhook
$ rchooks --recreate=https://example.com/webhook
$ rchooks --recreate=11112222-3333-4444-5555-66667777888
$ rchooks --delete=https://example.com/webhook
$ rchooks --delete=11112222-3333-4444-5555-66667777888
$ rchooks --env=~/.env --list
```

 [build-status-svg]: https://api.travis-ci.org/grokify/rchooks.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/rchooks
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/rchooks
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/rchooks
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/rchooks
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/rchooks/blob/master/LICENSE.md
