package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/ringcentral"
	"github.com/jessevdk/go-flags"

	"github.com/grokify/rchooks"
)

type Options struct {
	EnvFile   string `short:"e" long:"env" description:"Env filepath"`
	WhichEnv  []bool `short:"w" long:"which" description:"Which .env path"`
	NewToken  []bool `short:"n" long:"newToken" description:"Get New Token"`
	List      []bool `short:"l" long:"list" description:"List subscriptions"`
	Create    string `short:"c" long:"create" description:"Create subscription"`
	CreateEnv []bool `long:"createenv" description:"Create subscription from environment"`
	Delete    string `short:"d" long:"delete" description:"Delete subscription"`
	Recreate  string `short:"r" long:"recreate" description:"Recreate subscription"`
}

func isUrl(s string) bool { return regexp.MustCompile(`^https?://`).MatchString(strings.ToLower(s)) }

func handleResponse(info interface{}, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(info)
}

func GetCredentials() (ringcentral.Credentials, error) {
	return ringcentral.NewCredentialsJSONs(
		[]byte(os.Getenv("RC_APP")),
		[]byte(os.Getenv("RC_USER")),
		[]byte(os.Getenv("RINGCENTRAL_TOKEN")))
}

// This code takes a bot token and creates a permanent webhook.
func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	err = config.LoadEnvPathsPrioritized(opts.EnvFile, os.Getenv("ENV_PATH"))
	if err != nil {
		fmt.Printf("E [%v] ENV [%v]\n", opts.EnvFile, os.Getenv("ENV_PATH"))
		log.Fatal(err)
	}

	creds, err := GetCredentials()
	if err != nil {
		log.Fatal(err)
	}

	if len(opts.NewToken) > 0 {
		_, err = creds.NewToken()
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()

	appCfg := rchooks.NewRcHooksConfigCreds(
		creds,
		os.Getenv("RINGCENTRAL_WEBHOOK_DEFINITION_JSON"))

	rch, err := appCfg.InitilizeRcHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := appCfg.WebhookDefinition

	if len(opts.WhichEnv) > 0 {
		fmt.Printf("Using .env file [%v]\n", opts.EnvFile)
	}

	if len(opts.List) > 0 {
		handleResponse(rch.GetSubscriptions(ctx))
	}

	if len(opts.CreateEnv) > 0 {
		handleResponse(rch.CreateSubscription(ctx, req))
	}

	if len(opts.Create) > 0 {
		req.DeliveryMode.Address = opts.Create
		handleResponse(rch.CreateSubscription(ctx, req))
	}

	if len(opts.Delete) > 0 {
		handleResponse(rch.DeleteByIdOrUrl(ctx, opts.Delete))
	}

	if len(opts.Recreate) > 0 {
		handleResponse(rch.RecreateSubscriptionIdOrUrl(ctx, opts.Recreate))
	}

	fmt.Println("DONE")
}
