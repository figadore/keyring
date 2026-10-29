package main

import (
	"fmt"

	dbus "github.com/godbus/dbus/v5"
	ss "github.com/zalando/go-keyring/secret_service"
)

func secretKeys(account string) ([]string, error) {
	svc, err := ss.NewSecretService()
	if err != nil {
		return nil, err
	}
	// Use the same collection and attributes as go-keyring. Metadata searches
	// do not require opening a secret session or reading any secret values.
	items, err := svc.SearchItems(svc.GetLoginCollection(), map[string]string{"username": account})
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, item := range items {
		value, err := svc.Object("org.freedesktop.secrets", item).GetProperty("org.freedesktop.Secret.Item.Attributes")
		if err != nil {
			return nil, err
		}
		var attributes map[string]string
		if err := dbus.Store([]interface{}{value.Value()}, &attributes); err != nil {
			return nil, fmt.Errorf("decoding keyring attributes: %w", err)
		}
		if service, ok := attributes["service"]; ok && attributes["username"] == account {
			keys = append(keys, service)
		}
	}
	return keys, nil
}
