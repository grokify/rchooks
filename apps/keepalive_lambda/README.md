# RCHooks - RingCentral Webhook Keepalive Lambda Function

## Keepalive Lambda Function

### Installation

Build the lambda function and then upload to AWS:

```
$ go get github.com/grokify/rchooks
$ cd $GOPATH/src/github.com/grokify/rchooks/apps/keepalive_lambda
$ sh build_lambda.sh
```

Set the following enviroment variables:

* `RINGCENTRAL_TOKEN_JSON`
* `RINGCENTRAL_SERVER_URL`
* `RINGCENTRAL_WEBHOOK_DEFINITION_JSON`

### Configuration

Add a CloudWatch Event Rule with this Lambda function as the target.

## Notes

### Blacklist Reasons

* `I/O operation is failed. Details: [Read timed out]`
* `Webhook response exceeds max size. Read bytes count: [1024]`
* `Webhook responses with code: [404], reason: [Not Found]`