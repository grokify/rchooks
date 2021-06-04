// main.go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/grokify/rchooks"
	"github.com/pkg/errors"
)

func checkAndFixSubscription() (string, error) {
	appCfg := rchooks.NewRcHooksConfigEnv(
		"RINGCENTRAL_TOKEN", "RINGCENTRAL_SERVER_URL", "RINGCENTRAL_WEBHOOK_DEFINITION_JSON")

	ctx := context.Background()
	if rchooksUtil, err := appCfg.InitilizeRcHooks(ctx); err != nil {
		return "", errors.Wrap(err, "InitilizeRcHooks")
	} else if _, err := rchooksUtil.CheckAndFixSubscription(ctx, appCfg.WebhookDefinition); err != nil {
		return "", errors.Wrap(err, "CheckAndFixSubscription")
	} else {
		return "", nil
	}
}

// main Lambda Function should be called with CloudWatch Event Rule
func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(checkAndFixSubscription)
}
