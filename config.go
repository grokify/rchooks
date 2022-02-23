package rchooks

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/goauth/ringcentral"
	"github.com/grokify/mogo/errors/errorsutil"

	rc "github.com/grokify/go-ringcentral-client/office/v1/client"
	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
)

type RcHooksConfig struct {
	Token                 string `env:"RINGCENTRAL_TOKEN"`
	ServerUrl             string `env:"RINGCENTRAL_SERVER_URL"`
	WebhookDefinitionJson string `env:"RINGCENTRAL_WEBHOOK_DEFINITION_JSON"`
	WebhookDefinition     rc.CreateSubscriptionRequest
}

func NewRcHooksConfigCreds(creds credentials.Credentials, hookDefJson string) (RcHooksConfig, error) {
	cfg := RcHooksConfig{
		ServerUrl:             creds.OAuth2.ServerURL,
		WebhookDefinitionJson: hookDefJson}
	if creds.Token == nil {
		tok, err := creds.NewToken()
		if err != nil {
			return cfg, err
		}
		creds.Token = tok
	}
	if creds.Token != nil {
		cfg.Token = strings.TrimSpace(creds.Token.AccessToken)
	}
	return cfg, nil
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
		return errorsutil.Wrap(err, "Parse subscription definition")
	}
	rchConfig.WebhookDefinition = req
	return nil
}

func (rchConfig *RcHooksConfig) Client() (*http.Client, error) {
	return goauth.NewClientBearerTokenSimpleOrJson(context.Background(), []byte(rchConfig.Token))
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

func (rchConfig *RcHooksConfig) InitializeRcHooks(ctx context.Context) (RcHooks, error) {
	rchooksUtil := RcHooks{}

	err := rchConfig.Inflate()
	if err != nil {
		return rchooksUtil, err
	}

	httpClient, err := rchConfig.Client()
	if err != nil {
		return rchooksUtil, errorsutil.Wrap(err, "New client token")
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
