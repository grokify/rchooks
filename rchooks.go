package rchooks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	rc "github.com/grokify/go-ringcentral-client/office/v1/client"
	clientutil "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	WebhookStatusBlacklisted     = "Blacklisted"
	RingCentralAPIResponseFormat = `RingCentral_API_Status_Code [%v]`
	ExpiresMax                   = 499999999 // 15 years
	TransportTypePubNub          = "PubNub"
	TransportTypeWebHook         = "WebHook"
)

func ParseCreateSubscriptionRequest(data []byte) (rc.CreateSubscriptionRequest, error) {
	var req rc.CreateSubscriptionRequest
	err := json.Unmarshal(data, &req)
	return req, err
}

type RcHooks struct {
	Client *rc.APIClient
}

func (util *RcHooks) GetSubscriptions(ctx context.Context) (rc.RecordsCollectionResourceSubscriptionResponse, error) {
	info, resp, err := util.Client.PushNotificationsApi.GetSubscriptions(ctx)
	if err != nil && resp.StatusCode >= 300 {
		err = errorsutil.Wrap(err, string(clientutil.ApiResponseErrorBody(err)))
	}
	return info, httputilmore.CondenseResponseNot2xxToError(resp, err, "ERROR - Get Subscriptions API")
}

func (util *RcHooks) CreateSubscription(ctx context.Context, req rc.CreateSubscriptionRequest) (rc.SubscriptionInfo, error) {
	info, resp, err := util.Client.PushNotificationsApi.CreateSubscription(ctx, req)

	if err != nil && resp.StatusCode >= 300 {
		err = errorsutil.Wrap(err, string(clientutil.ApiResponseErrorBody(err)))
	}

	return info, httputilmore.CondenseResponseNot2xxToError(resp, err, "ERROR - Create Subscription API")
}

func (util *RcHooks) DeleteSubscription(ctx context.Context, subscriptionID string) error {
	resp, err := util.Client.PushNotificationsApi.DeleteSubscription(ctx, subscriptionID)
	return httputilmore.CondenseResponseNot2xxToError(resp, err,
		fmt.Sprintf("ERROR - Dete Subscription API id [%v]", subscriptionID))
}

func (util *RcHooks) DeleteBlacklisted(ctx context.Context, matches []rc.SubscriptionResponse) ([]rc.SubscriptionResponse, error) {
	goodMatches := []rc.SubscriptionResponse{}

	for _, match := range matches {
		if match.Status == WebhookStatusBlacklisted {
			err := util.DeleteSubscription(ctx, match.Id)
			if err != nil {
				return goodMatches, err
			}
		} else {
			goodMatches = append(goodMatches, match)
		}
	}
	return goodMatches, nil
}

func FilterSubscriptionsForRequest(ress []rc.SubscriptionResponse, req rc.CreateSubscriptionRequest) []rc.SubscriptionResponse {
	reqFiltersStringSorted := stringsutil.JoinStringsTrimSpaceToLowerSort(req.EventFilters, ",")
	reqWebhookURL := strings.TrimSpace(req.DeliveryMode.Address)

	matches := []rc.SubscriptionResponse{}
	for _, res := range ress {
		resFiltersStringSorted := stringsutil.JoinStringsTrimSpaceToLowerSort(res.EventFilters, ",")
		resWebhookURL := strings.TrimSpace(res.DeliveryMode.Address)
		if reqWebhookURL == resWebhookURL &&
			reqFiltersStringSorted == resFiltersStringSorted {
			matches = append(matches, res)
		}
	}
	return matches
}

func (util *RcHooks) RecreateSubscriptionIDOrURL(ctx context.Context, subIDOrURL string) ([]rc.SubscriptionInfo, error) {
	recreated := []rc.SubscriptionInfo{}

	subs, err := util.GetSubscriptions(ctx)
	if err != nil {
		return recreated, err
	}

	matches := []rc.SubscriptionResponse{}
	for _, sub := range subs.Records {
		if subIDOrURL == sub.Id || subIDOrURL == sub.DeliveryMode.Address {
			matches = append(matches, sub)
		}
	}
	if len(matches) == 0 {
		return recreated, fmt.Errorf("no matches found for [%v]", subIDOrURL)
	}

	for _, sub := range matches {
		newSubReq := NewCreateSubscriptionRequestPermahook(sub.EventFilters, sub.DeliveryMode.Address)
		if err := util.DeleteSubscription(ctx, sub.Id); err != nil {
			return recreated, err
		}
		if newSubInfo, err := util.CreateSubscription(ctx, newSubReq); err != nil {
			return recreated, err
		} else {
			recreated = append(recreated, newSubInfo)
		}
	}

	return recreated, nil
}

func (util *RcHooks) CheckAndFixSubscription(ctx context.Context, req rc.CreateSubscriptionRequest) (rc.SubscriptionInfo, error) {
	recreated := rc.SubscriptionInfo{}

	subs, err := util.GetSubscriptions(ctx)
	if err != nil {
		return recreated, err
	}

	matches := FilterSubscriptionsForRequest(subs.Records, req)

	if remaining, err := util.DeleteBlacklisted(ctx, matches); err != nil {
		return recreated, err
	} else if len(remaining) == 0 {
		return util.CreateSubscription(ctx, req)
	}
	return recreated, nil
}

func (util *RcHooks) DeleteByIDOrURL(ctx context.Context, idOrURLToDelete string) ([]rc.SubscriptionResponse, error) {
	deleted := []rc.SubscriptionResponse{}
	info, err := util.GetSubscriptions(ctx)
	if err != nil {
		return deleted, err
	}

	for _, subscription := range info.Records {
		if idOrURLToDelete == subscription.Id ||
			idOrURLToDelete == subscription.DeliveryMode.Address {
			resp, err := util.Client.PushNotificationsApi.DeleteSubscription(
				ctx, subscription.Id)
			err = httputilmore.CondenseResponseNot2xxToError(
				resp, err,
				fmt.Sprintf("ERROR - Delete Subscription API Id [%v]", subscription.Id))
			if err != nil {
				return deleted, err
			}
			deleted = append(deleted, subscription)
		}
	}
	return deleted, nil
}

func NewCreateSubscriptionRequestPermahook(eventFilters []string, hookURL string) rc.CreateSubscriptionRequest {
	return rc.CreateSubscriptionRequest{
		EventFilters: eventFilters,
		DeliveryMode: rc.NotificationDeliveryModeRequest{
			TransportType: TransportTypeWebHook,
			Address:       hookURL},
		ExpiresIn: int32(ExpiresMax)}
}

func NewCreateSubscriptionRequestPermahookBotSimple(hookURL string) rc.CreateSubscriptionRequest {
	return NewCreateSubscriptionRequestPermahook([]string{EventFilterMessagePosts}, hookURL)
}

type WebhookDefinitionThin struct {
	URL          string   `json:"url"`
	EventFilters []string `json:"eventFilters"`
}

func (thin *WebhookDefinitionThin) Full() rc.CreateSubscriptionRequest {
	return rc.CreateSubscriptionRequest{
		EventFilters: thin.EventFilters,
		DeliveryMode: rc.NotificationDeliveryModeRequest{
			TransportType: TransportTypeWebHook,
			Address:       thin.URL},
		ExpiresIn: int32(ExpiresMax)}
}
