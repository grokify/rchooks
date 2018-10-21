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
	"github.com/jessevdk/go-flags"

	"github.com/grokify/rchooks"
)

type Options struct {
	EnvFile   string `short:"e" long:"env" description:"Env filepath"`
	WhichEnv  []bool `short:"w" long:"which" description:"Which .env path"`
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
		fmt.Println("H")
		log.Fatal(err)
	}

	ctx := context.Background()

	appCfg := rchooks.NewRcHooksConfigEnv(
		"RINGCENTRAL_TOKEN",
		"RINGCENTRAL_SERVER_URL",
		"RINGCENTRAL_WEBHOOK_DEFINITION_JSON")

	rchooksUtil, err := appCfg.InitilizeRcHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := appCfg.WebhookDefinition

	if len(opts.WhichEnv) > 0 {
		fmt.Printf("Using .env file [%v]\n", opts.EnvFile)
	}

	if len(opts.List) > 0 {
		handleResponse(rchooksUtil.GetSubscriptions(ctx))
	}

	if len(opts.CreateEnv) > 0 {
		handleResponse(rchooksUtil.CreateSubscription(ctx, req))
	}

	if len(opts.Create) > 0 {
		req.DeliveryMode.Address = opts.Create
		handleResponse(rchooksUtil.CreateSubscription(ctx, req))
	}

	if len(opts.Delete) > 0 {
		handleResponse(rchooksUtil.DeleteByIdOrUrl(ctx, opts.Delete))
	}

	if len(opts.Recreate) > 0 {
		handleResponse(rchooksUtil.RecreateSubscriptionIdOrUrl(ctx, opts.Recreate))
	}

	fmt.Println("DONE")
}
