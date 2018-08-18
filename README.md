# RCHooks - RingCentral Webhook Tools

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]
[![Video][video-svg]][video-link]

`rchooks` is a toolset for creating, listing, recreating, and deleting RingCentral webhooks. It is especially useful in development when working with ngrok and when webhooks are blocked when test servers are taken down and non-responsive. It includes the following:

* `rchooks` CLI app
* `keepalive_lambda` AWS Lambda function to check and rebuild webhook when blacklisted
* `rchooks` SDK package for utilities to build your own apps

YouTube Tutorial Video: https://youtu.be/DYrzzJe8OyI

## Apps

### CLI App

#### CLI Installation

```
$ go get github.com/grokify/rchooks/apps/rchooks
$ rchooks --env=/path/to/.env --list
```

#### CLI Usage

By default, `rchooks` looks for an environment file path specified by the `ENV_PATH` environment variable or a `.env` file in the current working directory. You can also explicitly specify a `.env` file with the `--env` path parameter.

```
$ rchooks --list
$ rchooks --create=https://example.com/webhook
$ rchooks --recreate=https://example.com/webhook
$ rchooks --recreate=11112222-3333-4444-5555-66667777888
$ rchooks --delete=https://example.com/webhook
$ rchooks --delete=11112222-3333-4444-5555-66667777888
$ rchooks --env=~/.env --list
```

Set the following enviroment variables:

* `RINGCENTRAL_TOKEN` - JSON string or simple access token string
* `RINGCENTRAL_SERVER_URL`
* `RINGCENTRAL_WEBHOOK_DEFINITION_JSON` - Create subscription JSON body

### Keepalive Lambda Function

#### Installation

Build the lambda function and then upload to AWS:

```
$ go get github.com/grokify/rchooks
$ cd $GOPATH/src/github.com/grokify/rchooks/apps/keepalive_lambda
$ sh build_lambda.sh
```

Set the following enviroment variables:

* `RINGCENTRAL_TOKEN` - JSON string or simple access token string
* `RINGCENTRAL_SERVER_URL`
* `RINGCENTRAL_WEBHOOK_DEFINITION_JSON` - Create subscription JSON body

## Notes

### Blacklist Reasons

* `I/O operation is failed. Details: [Read timed out]`
* `Webhook response exceeds max size. Read bytes count: [1024]`
* `Webhook responses with code: [404], reason: [Not Found]`

 [build-status-svg]: https://api.travis-ci.org/grokify/rchooks.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/rchooks
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/rchooks
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/rchooks
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/rchooks
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/rchooks/blob/master/LICENSE.md
 [video-svg]: https://img.shields.io/badge/tutorial-YouTube-blue.svg
 [video-link]: https://youtu.be/DYrzzJe8OyI
