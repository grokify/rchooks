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
	TokenJson             string `env:"RINGCENTRAL_TOKEN_JSON"`
	ServerUrl             string `env:"RINGCENTRAL_SERVER_URL"`
	WebhookDefinitionJson string `env:"RINGCENTRAL_WEBHOOK_DEFINITION_JSON"`
	WebhookDefinition     rc.CreateSubscriptionRequest
}

func NewRcHooksConfigEnv(envVarTokenJson, envVarServerUrl, envVarHookDef string) RcHooksConfig {
	return RcHooksConfig{
		TokenJson:             os.Getenv(envVarTokenJson),
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

	if httpClient, err := om.NewClientTokenJSON(ctx, []byte(appCfg.TokenJson)); err != nil {
		return rchooksUtil, errors.Wrap(err, "New client token")
	} else if apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, appCfg.ServerUrl); err != nil {
		return rchooksUtil, err
	} else {
		rchooksUtil.Client = apiClient
		return rchooksUtil, nil
	}
}
