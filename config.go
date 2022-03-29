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
	ServerURL             string `env:"RINGCENTRAL_SERVER_URL"`
	WebhookDefinitionJSON string `env:"RINGCENTRAL_WEBHOOK_DEFINITION_JSON"`
	WebhookDefinition     rc.CreateSubscriptionRequest
}

func NewRcHooksConfigCreds(creds credentials.Credentials, hookDefJSON string) (RcHooksConfig, error) {
	cfg := RcHooksConfig{
		ServerURL:             creds.OAuth2.ServerURL,
		WebhookDefinitionJSON: hookDefJSON}
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

func NewRcHooksConfigEnv(envVarTokenOrJSON, envVarServerURL, envVarHookDef string) RcHooksConfig {
	return RcHooksConfig{
		Token:                 os.Getenv(envVarTokenOrJSON),
		ServerURL:             os.Getenv(envVarServerURL),
		WebhookDefinitionJSON: os.Getenv(envVarHookDef)}
}

func (rchConfig *RcHooksConfig) Inflate() error {
	rchConfig.WebhookDefinitionJSON = strings.TrimSpace(rchConfig.WebhookDefinitionJSON)
	if len(rchConfig.WebhookDefinitionJSON) <= 0 {
		return fmt.Errorf("E_NO_WEBHOOK_DEFINITION")
	}
	req, err := ParseCreateSubscriptionRequest([]byte(rchConfig.WebhookDefinitionJSON))
	if err != nil {
		return errorsutil.Wrap(err, "Parse subscription definition")
	}
	rchConfig.WebhookDefinition = req
	return nil
}

func (rchConfig *RcHooksConfig) Client() (*http.Client, error) {
	return goauth.NewClientBearerTokenSimpleOrJSON(context.Background(), []byte(rchConfig.Token))
}

func (rchConfig *RcHooksConfig) ClientUtil() (ringcentral.ClientUtil, error) {
	cu := ringcentral.ClientUtil{
		ServerURL: rchConfig.ServerURL}
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
		httpClient, rchConfig.ServerURL)
	if err != nil {
		return rchooksUtil, err
	} else {
		rchooksUtil.Client = apiClient
		return rchooksUtil, nil
	}
}
