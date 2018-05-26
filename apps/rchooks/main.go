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
	iom "github.com/grokify/gotilla/io/ioutilmore"
	"github.com/jessevdk/go-flags"

	"github.com/grokify/rchooks"
)

type Options struct {
	EnvFile  string `short:"e" long:"env" description:"Env filepath"`
	WhichEnv []bool `short:"w" long:"which" description:"Which .env path"`
	List     []bool `short:"l" long:"list" description:"List subscriptions"`
	Create   string `short:"c" long:"create" description:"Create subscription"`
	Delete   string `short:"d" long:"delete" description:"Delete subscription"`
	Recreate string `short:"r" long:"recreate" description:"Recreate subscription"`
}

func isUrl(s string) bool { return regexp.MustCompile(`^https?://`).MatchString(strings.ToLower(s)) }

// This code takes a bot token and creates a permanent webhook.
func main() {
	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	envFilesToCheck := []string{opts.EnvFile, os.Getenv("ENV_PATH"), "./.env"}

	for _, envFileTry := range envFilesToCheck {
		envFileTry = strings.TrimSpace(envFileTry)
		if len(envFileTry) > 0 {
			if exists, err := iom.IsFileWithSizeGtZero(envFileTry); err != nil {
				log.Fatal(err)
			} else if !exists {
				log.Fatal(fmt.Sprintf(".env file [%v] does not exist or is empty.", opts.EnvFile))
			} else {
				opts.EnvFile = envFileTry
				break
			}
		}
	}

	if len(opts.EnvFile) > 0 {
		if err := config.LoadDotEnvSkipEmpty(opts.EnvFile); err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()

	appCfg := rchooks.NewRcHooksConfigEnv("RINGCENTRAL_TOKEN_JSON", "RINGCENTRAL_SERVER_URL", "RINGCENTRAL_WEBHOOK_DEFINITION_JSON")

	rchooksUtil, err := appCfg.InitilizeRcHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := appCfg.WebhookDefinition

	if len(opts.WhichEnv) > 0 {
		fmt.Printf("Using .env file [%v]\n", opts.EnvFile)
	}

	if len(opts.Create) > 0 {
		req.DeliveryMode.Address = opts.Create

		info, err := rchooksUtil.CreateSubscription(ctx, req)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(info)
	}

	if len(opts.Delete) > 0 {
		info, err := rchooksUtil.DeleteByIdOrUrl(ctx, opts.Delete)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(info)
	}

	if len(opts.List) > 0 {
		info, err := rchooksUtil.GetSubscriptions(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(info)
	}

	if len(opts.Recreate) > 0 {
		info, err := rchooksUtil.RecreateSubscriptionIdOrUrl(ctx, opts.Recreate)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(info)
	}

	fmt.Println("DONE")
}
