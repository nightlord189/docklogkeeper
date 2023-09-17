package entity

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const TriggerMethodWebhook = "webhook"

const DefaultHeaders = `Content-Type:application/json`

type TriggerDB struct {
	ID             int64  `json:"id" gorm:"id"`
	Name           string `json:"name" gorm:"trigger_name"`
	ContainerName  string `json:"containerName" gorm:"container_name"`
	Contains       string `json:"contains" gorm:"contains"`
	NotContains    string `json:"notContains" gorm:"not_contains"`
	Regexp         string `json:"regexp" gorm:"regexp"`
	Method         string `json:"method" gorm:"method"`
	WebhookURL     string `json:"webhookURL" gorm:"webhook_url"`
	WebhookHeaders string `json:"webhookHeaders" gorm:"webhook_headers"`
	WebhookBody    string `json:"webhookBody" gorm:"webhook_body"`
}

/*
-- internal variables:
-- $dlk_container_full_name
-- $dlk_container_name
-- $dlk_log
-- $dlk_timestamp
*/

func (TriggerDB) TableName() string {
	return "trigger"
}

func (r *TriggerDB) IsValid() error {
	if r.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if r.ContainerName == "" {
		return fmt.Errorf("container name is empty")
	}
	if r.Contains == "" && r.NotContains == "" && r.Regexp == "" {
		return fmt.Errorf("all search criteria are empty")
	}
	if r.Contains == r.NotContains {
		return fmt.Errorf("contains and not contains values are empty")
	}
	if r.Regexp != "" {
		if _, err := regexp.Compile(r.Regexp); err != nil {
			return fmt.Errorf("invalid regexp: %w", err)
		}
	}
	if r.Method != TriggerMethodWebhook {
		return fmt.Errorf("invalid method")
	}
	if _, err := url.ParseRequestURI(r.WebhookURL); err != nil {
		return fmt.Errorf("invalid webhook url: %w", err)
	}
	if r.WebhookBody != "" && !json.Valid([]byte(r.WebhookBody)) {
		return fmt.Errorf("invalid webhook body")
	}
	return nil
}

func (r *TriggerDB) Match(logText string, reg *regexp.Regexp) bool {
	if r.Contains != "" && !strings.Contains(logText, r.Contains) {
		return false
	}
	if r.NotContains != "" && strings.Contains(logText, r.NotContains) {
		return false
	}
	if r.Regexp != "" && reg != nil {
		return reg.MatchString(logText)
	}
	return false
}
