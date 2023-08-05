/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://mozilla.org/MPL/2.0/
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

const (
	GetNetworkAdapterAPI    = "/api/ApplicationService/Network/AdapterConfigurations('%s')"
	UpdateNetworkAdapterAPI = "/api/ApplicationService/Actions/Network.ConfigureNetworkAdapter"
	GetNetworkSessions      = "/api/SessionService/SessionConfiguration"
	UpdateNetworkSessions   = "/api/SessionService/Actions/SessionService.SessionConfigurationUpdate"
	GetTimeConfiguration    = "/api/ApplicationService/Network/TimeConfiguration"
	GetTimeZone             = "/api/ApplicationService/Network/TimeZones"
	UpdateTimeConfiguration = "/api/ApplicationService/Network/TimeConfiguration"
	ProxyConfigurationAPI   = "/api/ApplicationService/Network/ProxyConfiguration"
)

func (c *Client) GetNetworkAdapterConfigByInterface(interfaceName string) (models.NetworkAdapterSetting, error) {
	path := fmt.Sprintf(GetNetworkAdapterAPI, interfaceName)
	response, err := c.Get(path, nil, nil)
	if err != nil {
		return models.NetworkAdapterSetting{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	networkAdapterSetting := models.NetworkAdapterSetting{}
	err = c.JSONUnMarshal(bodyData, &networkAdapterSetting)
	if err != nil {
		return models.NetworkAdapterSetting{}, err
	}
	return networkAdapterSetting, nil
}

func (c *Client) UpdateNetworkAdapterConfig(networkAdapter models.UpdateNetworkAdapterSetting) (models.Job, error) {
	data, _ := c.JSONMarshal(networkAdapter)
	response, err := c.Post(UpdateNetworkAdapterAPI, nil, data)
	if err != nil {
		return models.Job{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	jobResponse := models.Job{}
	err = c.JSONUnMarshal(respData, &jobResponse)
	fmt.Println(string(respData))
	return jobResponse, err
}

func (c *Client) GetNetworkSessions() (models.NetworkSessions, error) {
	response, err := c.Get(GetNetworkSessions, nil, nil)
	if err != nil {
		return models.NetworkSessions{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	networkSessions := models.NetworkSessions{}
	err = c.JSONUnMarshal(bodyData, &networkSessions)
	if err != nil {
		return models.NetworkSessions{}, err
	}
	return networkSessions, nil
}

func (c *Client) UpdateNetworkSessions(sessionPayload models.UpdateNetworkSessions) ([]models.SessionInfo, error) {
	data, _ := c.JSONMarshal(sessionPayload)
	response, err := c.Post(UpdateNetworkSessions, nil, data)
	if err != nil {
		return []models.SessionInfo{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	sessionResponse := []models.SessionInfo{}
	err = c.JSONUnMarshal(respData, &sessionResponse)
	fmt.Println(string(respData))
	return sessionResponse, err
}

func (c *Client) GetTimeConfiguration() (models.TimeConfiguration, error) {
	response, err := c.Get(GetTimeConfiguration, nil, nil)
	if err != nil {
		return models.TimeConfiguration{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	timeConfig := models.TimeConfiguration{}
	err = c.JSONUnMarshal(bodyData, &timeConfig)
	if err != nil {
		return models.TimeConfiguration{}, err
	}
	return timeConfig, nil
}

func (c *Client) UpdateTimeConfiguration(payloadTC models.TimeConfigPayload) (models.TimeConfigResponse, error) {
	data, _ := c.JSONMarshal(payloadTC)
	response, err := c.Put(UpdateTimeConfiguration, nil, data)
	if err != nil {
		return models.TimeConfigResponse{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)
	timeConfig := models.TimeConfigResponse{}
	err = c.JSONUnMarshal(bodyData, &timeConfig)
	if err != nil {
		return models.TimeConfigResponse{}, err
	}
	return timeConfig, nil
}

func (c *Client) GetTimeZone() (models.TimeZones, error) {
	response, err := c.Get(GetTimeZone, nil, nil)
	if err != nil {
		return models.TimeZones{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	timeZones := models.TimeZones{}
	err = c.JSONUnMarshal(bodyData, &timeZones)
	if err != nil {
		return models.TimeZones{}, err
	}
	return timeZones, nil
}

func (c *Client) GetProxyConfig() (models.ProxyConfiguration, error) {
	response, err := c.Get(ProxyConfigurationAPI, nil, nil)
	if err != nil {
		return models.ProxyConfiguration{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	proxyConfig := models.ProxyConfiguration{}
	err = c.JSONUnMarshal(bodyData, &proxyConfig)
	if err != nil {
		return models.ProxyConfiguration{}, err
	}
	return proxyConfig, nil
}

func (c *Client) UpdateProxyConfig(payloadProxy models.PayloadProxyConfiguration) (models.ProxyConfiguration, error) {
	data, _ := c.JSONMarshal(payloadProxy)
	response, err := c.Put(ProxyConfigurationAPI, nil, data)
	if err != nil {
		return models.ProxyConfiguration{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)
	proxyConfig := models.ProxyConfiguration{}
	err = c.JSONUnMarshal(bodyData, &proxyConfig)
	if err != nil {
		return models.ProxyConfiguration{}, err
	}
	return proxyConfig, nil
}
