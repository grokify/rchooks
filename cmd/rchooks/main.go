package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/grokify/goauth/credentials"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/type/stringsutil"
	flags "github.com/jessevdk/go-flags"

	"github.com/grokify/rchooks"
)

type Options struct {
	credentials.Options
	//EnvFile   string `short:"e" long:"env" description:"Env filepath"`
	//WhichEnv  []bool `short:"w" long:"which" description:"Which .env path"`
	List      []bool   `short:"l" long:"list" description:"List subscriptions"`
	Create    string   `short:"c" long:"create" description:"Create subscription"`
	CreateEnv []bool   `long:"createenv" description:"Create subscription from environment var"`
	Delete    string   `short:"d" long:"delete" description:"Delete subscription"`
	Recreate  string   `short:"r" long:"recreate" description:"Recreate subscription"`
	URL       string   `short:"u" long:"url" description:"URL for webhook"`
	Events    []string `short:"e" long:"events" description:"EventFilters"`
}

func (opts *Options) CanonicalEvents() ([]string, error) {
	sep := ","
	events := stringsutil.SliceCondenseSpace(
		strings.Split(
			strings.Join(opts.Events, sep),
			sep,
		),
		true, true)
	return rchooks.ConvertEvents(events...)
}

func handleResponse(info interface{}, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(info)
}

func GetCredentials(credspath, account string) (credentials.Credentials, error) {
	set, err := credentials.ReadFileCredentialsSet(credspath, true)
	if err != nil {
		return credentials.Credentials{},
			errorsutil.Wrap(err, fmt.Sprintf("creds file [%s]", credspath))
	}
	creds, err := set.Get(account)
	if err != nil {
		accts := strings.Join(set.Accounts(), ",")
		return credentials.Credentials{},
			errorsutil.Wrap(err, fmt.Sprintf("use [%s]", accts))
	}
	return creds, err
}

// This code takes a bot token and creates a permanent webhook.
func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := GetCredentials(opts.Options.CredsPath, opts.Options.Account)
	if err != nil {
		log.Fatal(err)
	}

	canEvts, err := opts.CanonicalEvents()
	if err != nil {
		log.Fatal(err)
	}
	hook := rchooks.WebhookDefinitionThin{
		URL:          opts.Create,
		EventFilters: canEvts}
	hookdef := string(jsonutil.MustMarshal(hook.Full(), true))

	ctx := context.Background()

	appCfg, err := rchooks.NewRcHooksConfigCreds(creds, hookdef)
	if err != nil {
		log.Fatal(err)
	}

	rch, err := appCfg.InitializeRcHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := appCfg.WebhookDefinition

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
		handleResponse(rch.DeleteByIDOrURL(ctx, opts.Delete))
	}

	if len(opts.Recreate) > 0 {
		fmt.Println("RECREATE")
		handleResponse(rch.RecreateSubscriptionIDOrURL(ctx, opts.Recreate))
	}

	fmt.Println("DONE")
}
