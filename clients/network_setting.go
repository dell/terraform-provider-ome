/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

// GetNetworkAdapterConfigByInterface to get adapter configuration of the interface.
func (c *Client) GetNetworkAdapterConfigByInterface(interfaceName string) (models.NetworkAdapterSetting, error) {
	path := fmt.Sprintf(GetNetworkAdapterAPI, interfaceName)
	response, err := c.Get(path, nil, nil)
	if err != nil {
		return models.NetworkAdapterSetting{}, err
	}

	bodyData, err := c.GetBodyData(response.Body)
	if err != nil {
		return models.NetworkAdapterSetting{}, err
	}
	networkAdapterSetting := models.NetworkAdapterSetting{}
	err = c.JSONUnMarshal(bodyData, &networkAdapterSetting)
	return networkAdapterSetting, err
}

// UpdateNetworkAdapterConfig to update the network adapter.
func (c *Client) UpdateNetworkAdapterConfig(networkAdapter models.UpdateNetworkAdapterSetting) (JobResp, error) {
	data, _ := c.JSONMarshal(networkAdapter)
	response, err := c.Post(UpdateNetworkAdapterAPI, nil, data)
	if err != nil {
		return JobResp{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	jobResponse := JobResp{}
	err = c.JSONUnMarshal(respData, &jobResponse)
	return jobResponse, err
}

// GetNetworkSessions to get the all sessions setting in the OME.
func (c *Client) GetNetworkSessions() (models.NetworkSessions, error) {
	response, err := c.Get(GetNetworkSessions, nil, nil)
	if err != nil {
		return models.NetworkSessions{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	networkSessions := models.NetworkSessions{}
	err = c.JSONUnMarshal(bodyData, &networkSessions)
	return networkSessions, err
}

// UpdateNetworkSessions to update the network session setting in the OME.
func (c *Client) UpdateNetworkSessions(sessionPayload []models.SessionInfo) ([]models.SessionInfo, error) {
	data, _ := c.JSONMarshal(sessionPayload)
	response, err := c.Post(UpdateNetworkSessions, nil, data)
	if err != nil {
		return []models.SessionInfo{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	sessionResponse := []models.SessionInfo{}
	err = c.JSONUnMarshal(respData, &sessionResponse)
	return sessionResponse, err
}

// GetTimeConfiguration to get the time configuration of the OME.
func (c *Client) GetTimeConfiguration() (models.TimeConfig, error) {
	response, err := c.Get(TimeConfigurationAPI, nil, nil)
	if err != nil {
		return models.TimeConfig{}, err
	}
	bodyData, _ := c.GetBodyData(response.Body)
	timeConfig := models.TimeConfig{}
	err = c.JSONUnMarshal(bodyData, &timeConfig)
	return timeConfig, err
}

// UpdateTimeConfiguration to update the time configuration of the OME.
func (c *Client) UpdateTimeConfiguration(payloadTC models.TimeConfig) (models.TimeConfig, error) {
	data, _ := c.JSONMarshal(payloadTC)
	response, err := c.Put(TimeConfigurationAPI, nil, data)
	if err != nil {
		return models.TimeConfig{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)
	timeConfig := models.TimeConfig{}
	err = c.JSONUnMarshal(bodyData, &timeConfig)
	return timeConfig, err
}

// GetTimeZone to get all time zone.
func (c *Client) GetTimeZone() (models.TimeZones, error) {
	response, err := c.Get(GetTimeZone, nil, nil)
	if err != nil {
		return models.TimeZones{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	timeZones := models.TimeZones{}
	err = c.JSONUnMarshal(bodyData, &timeZones)
	return timeZones, err
}

// GetProxyConfig to get the proxy configuration of the OME.
func (c *Client) GetProxyConfig() (models.ProxyConfiguration, error) {
	response, err := c.Get(ProxyConfigurationAPI, nil, nil)
	if err != nil {
		return models.ProxyConfiguration{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	proxyConfig := models.ProxyConfiguration{}
	err = c.JSONUnMarshal(bodyData, &proxyConfig)
	return proxyConfig, err
}

// UpdateProxyConfig to update the proxy configuration of the OME.
func (c *Client) UpdateProxyConfig(payloadProxy models.PayloadProxyConfiguration) (models.ProxyConfiguration, error) {
	data, _ := c.JSONMarshal(payloadProxy)
	response, err := c.Put(ProxyConfigurationAPI, nil, data)
	if err != nil {
		return models.ProxyConfiguration{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)
	proxyConfig := models.ProxyConfiguration{}
	err = c.JSONUnMarshal(bodyData, &proxyConfig)
	return proxyConfig, err
}
