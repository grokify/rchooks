package rchooks

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	om "github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/ringcentral"
	"github.com/pkg/errors"

	rc "github.com/grokify/go-ringcentral/office/v1/client"
	ru "github.com/grokify/go-ringcentral/office/v1/util"
)

type RcHooksConfig struct {
	Token                 string `env:"RINGCENTRAL_TOKEN"`
	ServerUrl             string `env:"RINGCENTRAL_SERVER_URL"`
	WebhookDefinitionJson string `env:"RINGCENTRAL_WEBHOOK_DEFINITION_JSON"`
	WebhookDefinition     rc.CreateSubscriptionRequest
}

func NewRcHooksConfigCreds(creds ringcentral.Credentials, hookDefJson string) RcHooksConfig {
	cfg := RcHooksConfig{
		ServerUrl:             creds.Application.ServerURL,
		WebhookDefinitionJson: hookDefJson}
	if creds.Token != nil {
		cfg.Token = strings.TrimSpace(creds.Token.AccessToken)
	}
	return cfg
}

func NewRcHooksConfigEnv(envVarTokenOrJson, envVarServerUrl, envVarHookDef string) RcHooksConfig {
	return RcHooksConfig{
		Token:                 os.Getenv(envVarTokenOrJson),
		ServerUrl:             os.Getenv(envVarServerUrl),
		WebhookDefinitionJson: os.Getenv(envVarHookDef)}
}

func (rchConfig *RcHooksConfig) Inflate() error {
	rchConfig.WebhookDefinitionJson = strings.TrimSpace(rchConfig.WebhookDefinitionJson)
	if len(rchConfig.WebhookDefinitionJson) <= 0 {
		return fmt.Errorf("E_NO_WEBHOOK_DEFINITION")
	}
	req, err := ParseCreateSubscriptionRequest([]byte(rchConfig.WebhookDefinitionJson))
	if err != nil {
		return errors.Wrap(err, "Parse subscription definition")
	}
	rchConfig.WebhookDefinition = req
	return nil
}

func (rchConfig *RcHooksConfig) Client() (*http.Client, error) {
	return om.NewClientBearerTokenSimpleOrJson(context.Background(), []byte(rchConfig.Token))
}

func (rchConfig *RcHooksConfig) ClientUtil() (ringcentral.ClientUtil, error) {
	cu := ringcentral.ClientUtil{
		ServerURL: rchConfig.ServerUrl}
	client, err := rchConfig.Client()
	if err != nil {
		return cu, err
	}
	cu.Client = client
	return cu, err
}

func (rchConfig *RcHooksConfig) InitilizeRcHooks(ctx context.Context) (RcHooks, error) {
	rchooksUtil := RcHooks{}

	err := rchConfig.Inflate()

	// mt.Printf("RCHOOKS_INIT: [%v]\n", rchConfig.Token)

	httpClient, err := rchConfig.Client()
	if err != nil {
		return rchooksUtil, errors.Wrap(err, "New client token")
	}

	apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, rchConfig.ServerUrl)
	if err != nil {
		return rchooksUtil, err
	} else {
		rchooksUtil.Client = apiClient
		return rchooksUtil, nil
	}
}
