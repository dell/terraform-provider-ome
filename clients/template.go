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
	"io"
	"strconv"
	"strings"
	"terraform-provider-ome/models"
)

// GetTemplateAttributes returns the editable attributes for a given templateID
func (c *Client) GetTemplateAttributes(templateID int64, stateAttributes []models.Attribute, refreshAll bool) ([]models.OmeAttribute, error) {
	attributesResp, err := c.Get(fmt.Sprintf(TemplateAPI+"(%d)/%s", templateID, "AttributeDetails"), nil, nil)
	if err != nil {
		return nil, err
	}
	respBody := attributesResp.Body
	attributeRespBody, _ := io.ReadAll(respBody)
	err = respBody.Close()
	if err != nil {
		return nil, err
	}
	attrGroups := models.OMETemplateAttrGroups{}
	err = c.JSONUnMarshal(attributeRespBody, &attrGroups)
	if err != nil {
		return nil, err
	}

	omeAttributes := []models.OmeAttribute{}
	if refreshAll {
		refreshAllAttributes(attrGroups.AttributeGroups, &omeAttributes)
	} else {

		for _, stateAttribute := range stateAttributes {
			attribute, err := getUpdatedAttribute(stateAttribute.DisplayName.ValueString(), stateAttribute.AttributeID.ValueInt64(), attrGroups.AttributeGroups)
			if err != nil {
				return nil, err
			}
			updatedOMEAttribute := models.OmeAttribute{
				AttributeID: attribute.AttributeID,
				DisplayName: stateAttribute.DisplayName.ValueString(),
				Value:       attribute.Value,
				IsIgnored:   attribute.IsIgnored,
			}

			omeAttributes = append(omeAttributes, updatedOMEAttribute)
		}
	}

	return omeAttributes, nil
}

func findSubAttrGroup(displayName string, counter *int, subAttrGroups []models.SubAttributeGroup) (models.SubAttributeGroup, error) {
	attributeHierarchy := strings.Split(displayName, ",")
	for _, subAttrGroup := range subAttrGroups {
		if subAttrGroup.DisplayName == attributeHierarchy[*counter] {
			*counter = *counter + 1
			if *counter == len(attributeHierarchy)-1 {
				return subAttrGroup, nil
			}
			return findSubAttrGroup(displayName, counter, subAttrGroup.SubAttributeGroups)
		}
	}
	return models.SubAttributeGroup{}, fmt.Errorf("SubAttributeGroup could not be found")
}

func getUpdatedAttribute(displayName string, attributeID int64, omeAttributes []models.AttributeGroup) (models.OmeAttribute, error) {
	attributeHierarchy := strings.Split(displayName, ",")
	counter := 0
	for _, attrGroup := range omeAttributes {
		if attrGroup.DisplayName == attributeHierarchy[counter] {
			counter++
			subAttrGroup, err := findSubAttrGroup(displayName, &counter, attrGroup.SubAttributeGroups)
			if err != nil {
				return models.OmeAttribute{}, err
			}
			for _, attribute := range subAttrGroup.Attributes {
				if attributeID == attribute.AttributeID && attribute.DisplayName == attributeHierarchy[counter] {
					return attribute, nil
				}
			}
			return models.OmeAttribute{}, fmt.Errorf("attribute could not be found")
		}
	}
	return models.OmeAttribute{}, fmt.Errorf("AttributeGroup could not be found")
}

func refreshAllAttributes(omeAttributeGroups []models.AttributeGroup, omeAttributes *[]models.OmeAttribute) {

	for _, attrGroup := range omeAttributeGroups {
		findAttributes(attrGroup.DisplayName, attrGroup.SubAttributeGroups, omeAttributes)
	}

}

func findAttributes(displayName string, subAttrGroups []models.SubAttributeGroup, omeAttributes *[]models.OmeAttribute) {
	for _, currentSubAttrGroup := range subAttrGroups {
		if len(currentSubAttrGroup.SubAttributeGroups) > 0 {
			displayName := fmt.Sprintf("%s,%s", displayName, currentSubAttrGroup.DisplayName)
			findAttributes(displayName, currentSubAttrGroup.SubAttributeGroups, omeAttributes)

		} else {
			attributes := currentSubAttrGroup.Attributes
			for _, attr := range attributes {

				attribute := models.OmeAttribute{
					DisplayName: fmt.Sprintf("%s,%s,%s", displayName, currentSubAttrGroup.DisplayName, attr.DisplayName),
					AttributeID: attr.AttributeID,
					Value:       attr.Value,
					IsIgnored:   attr.IsIgnored,
				}
				*omeAttributes = append(*omeAttributes, attribute)
			}
		}
	}
}

// CreateTemplate creates a template from a reference device id.
func (c *Client) CreateTemplate(ut models.CreateTemplate) (int64, error) {
	data, _ := c.JSONMarshal(ut)
	response, err := c.Post(TemplateAPI, nil, data)
	if err != nil {
		return -1, err
	}
	respData, _ := c.GetBodyData(response.Body)
	val, _ := strconv.ParseInt(string(respData), 10, 64)
	return val, nil
}

// GetViewTypeID gets the viewTypeID based on the view_type.
func (c *Client) GetViewTypeID(viewType string) (int64, error) {
	var viewTypeID int64 = -1
	var err error

	response, err := c.Get(TemplateViewTypeAPI, nil, nil)
	if err != nil {
		return -1, err
	}

	respData, _ := c.GetBodyData(response.Body)
	vt := models.ViewTypes{}

	err = c.JSONUnMarshal(respData, &vt)
	if err != nil {
		return -1, err
	}

	for _, vt2 := range vt.Value {
		if strings.EqualFold(vt2.Description, viewType) {
			viewTypeID = vt2.ID
			break
		}
	}

	return viewTypeID, err
}

// GetDeviceTypeID gets the viewTypeID based on the view_type.
func (c *Client) GetDeviceTypeID(deviceType string) (int64, error) {
	var deviceTypeID int64 = -1
	var err error

	response, err := c.Get(TemplateDeviceTypeAPI, nil, nil)
	if err != nil {
		return -1, err
	}

	respData, _ := c.GetBodyData(response.Body)
	dt := models.DeviceTypes{}

	err = c.JSONUnMarshal(respData, &dt)
	if err != nil {
		return -1, err
	}

	for _, dt2 := range dt.Value {
		if strings.EqualFold(dt2.Name, deviceType) {
			deviceTypeID = dt2.ID
			break
		}
	}

	return deviceTypeID, err
}

// GetTemplateByID gets the viewTypeID based on the view_type .
func (c *Client) GetTemplateByID(id int64) (models.OMETemplate, error) {
	omeTemplate := models.OMETemplate{}
	response, err := c.Get(fmt.Sprintf(TemplateAPI+"(%d)", id), nil, nil)
	if err != nil {
		return omeTemplate, err
	}

	respData, _ := c.GetBodyData(response.Body)

	err = c.JSONUnMarshal(respData, &omeTemplate)
	if err != nil {
		return omeTemplate, err
	}

	return omeTemplate, nil
}

// GetTemplateByName returns the template for the given template name
func (c *Client) GetTemplateByName(name string) (models.OMETemplate, error) {
	omeTemplateResponse := []models.OMETemplate{}
	err := c.GetPaginatedDataWithQueryParam(TemplateAPI, map[string]string{"$filter": fmt.Sprintf("%s eq '%s'", "Name", name)}, &omeTemplateResponse)
	if err != nil {
		return models.OMETemplate{}, err
	}
	if len(omeTemplateResponse) == 0 {
		return models.OMETemplate{}, nil
	}
	for _, template := range omeTemplateResponse {
		if template.Name == name {
			return template, nil
		}
	}
	return models.OMETemplate{}, nil
}

// UpdateTemplate updates a template from a reference template id.
func (c *Client) UpdateTemplate(ut models.UpdateTemplate) error {
	data, _ := c.JSONMarshal(ut)
	uri := fmt.Sprintf(TemplateAPI+"(%d)", ut.ID)
	_, err := c.Put(uri, nil, data)
	return err
}

// GetIdentityPoolByName returns the identityPool for the given identityPoolName
func (c *Client) GetIdentityPoolByName(name string) (models.IdentityPool, error) {
	response, err := c.Get(IdentityPoolAPI, nil, nil)
	if err != nil {
		return models.IdentityPool{}, err
	}
	b, _ := c.GetBodyData(response.Body)

	omeIdentityPools := models.OMEIdentityPools{}
	err = c.JSONUnMarshal(b, &omeIdentityPools)
	if err != nil {
		return models.IdentityPool{}, err
	}
	for _, omeIdentityPool := range omeIdentityPools.Value {
		if name == omeIdentityPool.Name {
			return omeIdentityPool, nil
		}
	}

	return models.IdentityPool{}, fmt.Errorf(ErrInvalidIdentityPool, name)
}

// GetIdentityPoolByID returns the identityPool for the given identityPoolID
func (c *Client) GetIdentityPoolByID(id int64) (models.IdentityPool, error) {
	response, err := c.Get(fmt.Sprintf(IdentityPoolAPI+"(%d)", id), nil, nil)
	if err != nil {
		return models.IdentityPool{}, err
	}
	b, _ := c.GetBodyData(response.Body)

	omeIdentityPool := models.IdentityPool{}
	err = c.JSONUnMarshal(b, &omeIdentityPool)
	if err != nil {
		return models.IdentityPool{}, err
	}

	return omeIdentityPool, nil
}

// GetPayloadVlanAttribute returns the vlan attribute for a specific (nicIdentifier,port) combination
func (c *Client) GetPayloadVlanAttribute(networkView models.NetworkSpecificView, nicIdentifier string, port int64) (models.PayloadVlanAttribute, error) {
	vlanAttrs := models.PayloadVlanAttribute{}
	for _, nag := range networkView.NetworkAttributeGroups {
		if nag.DisplayName == NICModel {
			for _, nics := range nag.SubAttributeGroups {
				if nics.DisplayName == nicIdentifier {
					for _, ports := range nics.SubAttributeGroups {
						if ports.DisplayName == Port && ports.GroupNameID == port && len(ports.SubAttributeGroups) != 0 {
							vlanAttrs := getVlanAttributeForPayload(ports.SubAttributeGroups[0].NetworkAttributes) //0 Assuming ports have only one partition
							return vlanAttrs, nil
						}

					}
				}
			}
		}
	}
	return vlanAttrs, fmt.Errorf(ErrInvalidNetworkDetails)
}

func getVlanAttributeForPayload(networkAttributes []models.NetworkAttribute) models.PayloadVlanAttribute {
	vlanAttr := models.PayloadVlanAttribute{}
	for _, na := range networkAttributes {
		switch na.DisplayName {
		case NICBondingEnabled:
			vlanAttr.IsNICBonded, _ = strconv.ParseBool(na.Value)
			vlanAttr.ComponentID = na.ComponentID
		case VlanTagged:
			vlanAttr.Tagged = getStringArrayToArrayInt(strings.Split(na.Value, ","))
		case VlanUntagged:
			vlanAttr.Untagged, _ = strconv.ParseInt(na.Value, 10, 64)
		}
	}
	return vlanAttr

}

// GetVlanNetworkModel returns the network view of the template returning all Network attributes
func (c *Client) GetVlanNetworkModel(templateID int64) (models.NetworkSpecificView, error) {
	uri := fmt.Sprintf(TemplateAPI+"(%d)/Views(4)/AttributeViewDetails", templateID)
	resp, err := c.Get(uri, nil, nil)
	if err != nil {
		return models.NetworkSpecificView{}, err
	}
	respBody, _ := c.GetBodyData(resp.Body)

	networksView := models.NetworkSpecificView{}
	err = c.JSONUnMarshal(respBody, &networksView)
	if err != nil {
		return models.NetworkSpecificView{}, err
	}
	return networksView, nil
}

// UpdateNetworkConfig updates the network attributes to the template
func (c *Client) UpdateNetworkConfig(nwConfig *models.UpdateNetworkConfig) error {
	data, _ := c.JSONMarshal(nwConfig)
	_, err := c.Post(UpdateNetworkConfigAPI, nil, data)
	return err
}

func getStringArrayToArrayInt(strSlice []string) []int64 {
	var intSlice = make([]int64, len(strSlice))
	if len(strSlice) == 1 && strSlice[0] == "0" {
		return []int64{}
	}
	for i, v := range strSlice {
		intSlice[i], _ = strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	}
	return intSlice
}

// GetSchemaVlanData returns the vlan data from OME for a template as per the schema of the vlan model in the state
func (c *Client) GetSchemaVlanData(templateID int64) (models.OMEVlan, error) {
	omeVlanAttr := models.OMEVlan{}
	netSpecView, err := c.GetVlanNetworkModel(templateID)
	if err != nil {
		return omeVlanAttr, err
	}

	for _, nag := range netSpecView.NetworkAttributeGroups {
		switch nag.DisplayName {
		case NICBondingTechnologyAttributeGrp:
			omeVlanAttr.BondingTechnology = getAttributeValue(nag.NetworkAttributes, NICBondingTechnologyAttribute)
		case NICModel:
			omeVlanAttr.OMEVlanAttributes = getAllVlanAttributes(nag.SubAttributeGroups)
		default:
			continue
		}
	}
	return omeVlanAttr, nil
}

// GetTemplateByIDOrName - method to get template information by ID or name.
func (c *Client) GetTemplateByIDOrName(templateID int64, templateName string) (models.OMETemplate, error) {
	template, err := c.GetTemplateByID(templateID)
	if err != nil && templateName != "" {
		template, err = c.GetTemplateByName(templateName)
	}
	return template, err
}

// CloneTemplateByRefTemplateID - method to clone template using reference template ID.
func (c *Client) CloneTemplateByRefTemplateID(cloneTemplateRequest models.OMECloneTemplate) (int64, error) {
	data, _ := c.JSONMarshal(cloneTemplateRequest)
	response, err := c.Post(CloneTemplateAPI, nil, data)
	if err != nil {
		return -1, err
	}

	respData, _ := c.GetBodyData(response.Body)
	newTemplateID, _ := strconv.ParseInt(string(respData), 10, 64)
	return newTemplateID, nil
}

// ImportTemplate - method to clone template using reference template ID.
func (c *Client) ImportTemplate(importTemplateRequest models.OMEImportTemplate) (int64, error) {
	data, _ := c.JSONMarshal(importTemplateRequest)
	response, err := c.Post(ImportTemplateAPI, nil, data)
	if err != nil {
		return -1, err
	}
	respData, _ := c.GetBodyData(response.Body)
	newTemplateID, _ := strconv.ParseInt(string(respData), 10, 64)
	return newTemplateID, nil
}

func getAllVlanAttributes(nags []models.NetworkAttributeGroup) []models.OMEVlanAttribute {
	vlanAttrs := []models.OMEVlanAttribute{}
	for _, nicIdentifier := range nags { // Loops NIC identifiers
		ports := nicIdentifier.SubAttributeGroups
		for _, port := range ports { // Loops ports of each nic
			vlanAttr := models.OMEVlanAttribute{}
			vlanAttr.Port = port.GroupNameID
			vlanAttr.NicIdentifier = nicIdentifier.DisplayName
			if len(port.SubAttributeGroups) != 0 {
				attrs := port.SubAttributeGroups[0].NetworkAttributes
				vlanAttr.ComponentID = attrs[0].ComponentID
				vlanAttr.IsNICBonded, _ = strconv.ParseBool(getAttributeValue(attrs, NICBondingEnabled))
				vlanAttr.Tagged = getStringArrayToArrayInt(strings.Split(getAttributeValue(attrs, VlanTagged), ","))
				vlanAttr.Untagged, _ = strconv.ParseInt(getAttributeValue(attrs, VlanUntagged), 10, 64)
			}
			vlanAttrs = append(vlanAttrs, vlanAttr)
		}

	}
	return vlanAttrs
}
func getAttributeValue(nas []models.NetworkAttribute, displayName string) string {
	for _, na := range nas {
		if na.DisplayName == displayName {
			return na.Value
		}
	}
	return ""
}
