package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

// CreateBaseline creates a baseline with baseline target devices and notification settings.
func (c *Client) CreateBaseline(baseline models.ConfigurationBaselinePayload) (models.OmeBaseline, error) {
	data, _ := c.JSONMarshal(baseline)
	response, err := c.Post(BaselineAPI, nil, data)
	if err != nil {
		return models.OmeBaseline{}, err
	}
	respData, _ := c.GetBodyData(response.Body)

	omeBaseline := models.OmeBaseline{}
	err = c.JSONUnMarshal(respData, &omeBaseline)
	return omeBaseline, err
}

// UpdateBaseline updates a baseline with baseline target devices and notification settings.
func (c *Client) UpdateBaseline(baseline models.ConfigurationBaselinePayload) (models.OmeBaseline, error) {
	data, _ := c.JSONMarshal(baseline)
	response, err := c.Put(fmt.Sprintf(BaselineAPI+"(%d)", baseline.ID), nil, data)
	if err != nil {
		return models.OmeBaseline{}, err
	}
	respData, _ := c.GetBodyData(response.Body)

	omeBaseline := models.OmeBaseline{}
	err = c.JSONUnMarshal(respData, &omeBaseline)
	return omeBaseline, err
}

// DeleteBaseline deletea a baseline.
func (c *Client) DeleteBaseline(baselineIDs []int64) error {
	baselineIds := models.BaseLineIDsData{BaselineIDs: baselineIDs}
	body, _ := c.JSONMarshal(baselineIds)
	_, err := c.Post(BaseLineRemoveAPI, nil, body)
	if err != nil {
		return err
	}
	return nil
}

// GetBaselineByID gets the baseline details by baseline ID .
func (c *Client) GetBaselineByID(id int64) (models.OmeBaseline, error) {
	omeBaseline := models.OmeBaseline{}
	response, err := c.Get(fmt.Sprintf(BaselineByIDAPI, id), nil, nil)
	if err != nil {
		return omeBaseline, err
	}

	respData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(respData, &omeBaseline)
	return omeBaseline, err
}

// GetBaselineByName gets the baseline details by baseline name .
func (c *Client) GetBaselineByName(name string) (models.OmeBaseline, error) {
	omeBaseline, err := c.getBaseline(BaselineAPI, name)
	if err != nil {
		return models.OmeBaseline{}, err
	}
	return omeBaseline, nil
}

// GetBaselineDevComplianceReportsByID gets baseline device compliance report by baseline ID as string
func (c *Client) GetBaselineDevComplianceReportsByID(baselineID int64) ([]models.OMEComplianceReports, error) {
	cr := []models.OMEComplianceReports{}
	err := c.GetPaginatedData(fmt.Sprintf(BaselineDeviceComplianceReportsAPI, baselineID), &cr)
	if err != nil {
		return []models.OMEComplianceReports{}, err
	}
	return cr, err
}

// GetBaselineDevAttrComplianceReportsByID gets baseline device attribute compliance report by baseline ID and device ID as string
func (c *Client) GetBaselineDevAttrComplianceReportsByID(baselineID int64, deviceID int64) (string, error) {
	response, err := c.Get(fmt.Sprintf(BaselineDeviceAttrComplianceReportsAPI, baselineID, deviceID), nil, nil)
	if err != nil {
		return "", err
	}
	respData, _ := c.GetBodyData(response.Body)
	return string(respData), err
}

func (c *Client) getBaseline(url, name string) (models.OmeBaseline, error) {
	omeBaselines := models.OmeBaselines{}
	response, err := c.Get(url, nil, nil)
	if err != nil {
		return models.OmeBaseline{}, err
	}

	respData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(respData, &omeBaselines)
	if err != nil {
		return models.OmeBaseline{}, err
	}
	for _, omeBaseline := range omeBaselines.Value {
		if omeBaseline.Name == name {
			return omeBaseline, nil
		}
	}
	for omeBaselines.NextLink != "" {
		return c.getBaseline(omeBaselines.NextLink, name)
	}
	return models.OmeBaseline{}, fmt.Errorf(ErrBaselineNameNotFound, name)
}
