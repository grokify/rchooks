package rchooks

import (
	"fmt"
	"strings"
)

const (
	EventFilterMessagePosts  = "/restapi/v1.0/account/~/extension/~/glip/posts"
	EventFilterMessageGroups = "/restapi/v1.0/account/~/extension/~/glip/groups"
	EventFilterSmsInbound    = "/restapi/v1.0/account/~/a2p-sms/messages?direction=Inbound"
	EventFilterSmsBatch      = "/restapi/v1.0/account/~/a2p-sms/batch"
	EventFilterSmsOptOuts    = "/restapi/v1.0/account/~/a2p-sms/opt-outs"

	SlugMessageGroups = "msggroups"
	SlugMessagePosts  = "msgposts"
	SlugSmsBatch      = "a2psmsbatch"
	SlugSmsInbound    = "a2psmsinbound"
	SlugSmsOptOuts    = "a2psmsoptouts"
)

func SlugToFilterMap() map[string]string {
	return map[string]string{
		SlugMessageGroups: EventFilterMessageGroups,
		SlugMessagePosts:  EventFilterMessagePosts,
		SlugSmsBatch:      EventFilterSmsBatch,
		SlugSmsInbound:    EventFilterSmsInbound,
		SlugSmsOptOuts:    EventFilterSmsOptOuts}
}

func ConvertEvent(slug string) (string, error) {
	slug = strings.ToLower(strings.TrimSpace(slug))
	if strings.Index(slug, "/") == 0 {
		return slug, nil
	}
	s2f := SlugToFilterMap()
	if filter, ok := s2f[slug]; ok {
		return filter, nil
	}
	return "", fmt.Errorf("slug not found [%s]", slug)
}

func ConvertEvents(slugs ...string) ([]string, error) {
	filters := []string{}
	for _, slug := range slugs {
		filter, err := ConvertEvent(slug)
		if err != nil {
			return filters, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}
