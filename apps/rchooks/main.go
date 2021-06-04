package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/oauth2more/ringcentral"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/jessevdk/go-flags"

	"github.com/grokify/rchooks"
)

type Options struct {
	CredsPath string `long:"creds" description:"Credentials filepath"`
	CredsUser string `long:"user" description:"Credentials user key"`
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

func GetCredentials(opts Options) (ringcentral.Credentials, error) {
	if len(opts.CredsPath) > 0 {
		credsSet, err := ringcentral.ReadFileCredentialsSet(opts.CredsPath)
		if err != nil {
			return ringcentral.Credentials{}, err
		}
		return credsSet.Get(opts.CredsUser)
	}
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

	if len(opts.CredsPath) == 0 {
		err = config.LoadEnvPathsPrioritized(opts.EnvFile, os.Getenv("ENV_PATH"))
		if err != nil {
			fmt.Printf("E [%v] ENV [%v]\n", opts.EnvFile, os.Getenv("ENV_PATH"))
			log.Fatal(err)
		}
	}

	creds, err := GetCredentials(opts)
	if err != nil {
		log.Fatal(err)
	}

	if len(opts.NewToken) > 0 {
		_, err = creds.NewToken()
		if err != nil {
			log.Fatal(err)
		}
	}

	hookdef := os.Getenv("RINGCENTRAL_WEBHOOK_DEFINITION_JSON")

	if 1 == 1 {
		hook := rchooks.WebhookDefinitionThin{
			URL: "https://320ea976a693.ngrok.io/webhook",
			EventFilters: []string{
				"/restapi/v1.0/account/~/a2p-sms/messages?direction=Inbound",
				"/restapi/v1.0/account/~/a2p-sms/batch",
				"/restapi/v1.0/account/~/a2p-sms/opt-outs"}}
		hookdef = string(jsonutil.MustMarshal(hook.Full(), true))
	}

	fmtutil.PrintJSON(creds)
	fmtutil.PrintJSON(hookdef)

	ctx := context.Background()

	appCfg, err := rchooks.NewRcHooksConfigCreds(creds, hookdef)
	if err != nil {
		log.Fatal(err)
	}

	rch, err := appCfg.InitilizeRcHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := appCfg.WebhookDefinition

	if len(opts.WhichEnv) > 0 {
		fmt.Printf("Using .env file [%v]\n", opts.EnvFile)
	}

	if len(opts.List) > 0 {
		fmt.Println("LIST")
		handleResponse(rch.GetSubscriptions(ctx))
	}

	if len(opts.CreateEnv) > 0 {
		fmt.Println("CREATE_ENV")
		handleResponse(rch.CreateSubscription(ctx, req))
	}

	if len(opts.Create) > 0 {
		req.DeliveryMode.Address = opts.Create
		fmt.Println("CREATE")
		handleResponse(rch.CreateSubscription(ctx, req))
	}

	if len(opts.Delete) > 0 {
		fmt.Println("DELETE")
		handleResponse(rch.DeleteByIdOrUrl(ctx, opts.Delete))
	}

	if len(opts.Recreate) > 0 {
		fmt.Println("RECREATE")
		handleResponse(rch.RecreateSubscriptionIdOrUrl(ctx, opts.Recreate))
	}

	fmt.Println("DONE")
}
