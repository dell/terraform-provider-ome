package ome

import (
	"fmt"
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func getSweeperClient(region string) (*clients.Client, error) {
	url := clients.GetURL(omeHost, defaultPort)
	clientOptions := clients.ClientOptions{
		Username:       omeUserName,
		Password:       omePassword,
		URL:            url,
		SkipSSL:        true,
		Timeout:        defaultTimeout,
		Retry:          clients.Retries,
		PreRequestHook: clients.ClientPreReqHook,
	}
	omeClient, err := clients.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Unable to create sweeper client %s", err.Error())
	}
	return omeClient, nil
}
