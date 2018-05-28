package rchooks

import (
	"context"
	"os"

	om "github.com/grokify/oauth2more"
	"github.com/pkg/errors"

	rc "github.com/grokify/go-ringcentral/client"
	ru "github.com/grokify/go-ringcentral/clientutil"
)

type RcHooksConfig struct {
	Token                 string `env:"RINGCENTRAL_TOKEN"`
	ServerUrl             string `env:"RINGCENTRAL_SERVER_URL"`
	WebhookDefinitionJson string `env:"RINGCENTRAL_WEBHOOK_DEFINITION_JSON"`
	WebhookDefinition     rc.CreateSubscriptionRequest
}

func NewRcHooksConfigEnv(envVarTokenOrJson, envVarServerUrl, envVarHookDef string) RcHooksConfig {
	return RcHooksConfig{
		Token:                 os.Getenv(envVarTokenOrJson),
		ServerUrl:             os.Getenv(envVarServerUrl),
		WebhookDefinitionJson: os.Getenv(envVarHookDef)}
}

func (appCfg *RcHooksConfig) InitilizeRcHooks(ctx context.Context) (RcHooks, error) {
	rchooksUtil := RcHooks{}

	if req, err := ParseCreateSubscriptionRequest([]byte(appCfg.WebhookDefinitionJson)); err != nil {
		return rchooksUtil, errors.Wrap(err, "Parse subscription definition")
	} else {
		appCfg.WebhookDefinition = req
	}

	if httpClient, err := om.NewClientBearerTokenSimpleOrJson(ctx, []byte(appCfg.Token)); err != nil {
		return rchooksUtil, errors.Wrap(err, "New client token")
	} else if apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, appCfg.ServerUrl); err != nil {
		return rchooksUtil, err
	} else {
		rchooksUtil.Client = apiClient
		return rchooksUtil, nil
	}
}
