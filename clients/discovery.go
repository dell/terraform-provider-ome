package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

func (c *Client) CreateDiscoveryJob(discoveryJob models.DiscoveryJobPayload) (models.DiscoveryJobResponse, error) {
	data, _ := c.JSONMarshal(discoveryJob)
	response, err := c.Post(DiscoveryJobAPI, nil, data)
	if err != nil {
		return models.DiscoveryJobResponse{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeDiscoveryJob := models.DiscoveryJobResponse{}
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	return omeDiscoveryJob, err
}

func (c *Client) UpdateDiscoveryJob(discoveryJob models.DiscoveryJobPayload) (models.DiscoveryJobResponse, error){
	data, _ := c.JSONMarshal(discoveryJob)
	response, err := c.Post(fmt.Sprintf(DiscoveryJobAPI+"?groupId=%d",discoveryJob.DiscoveryConfigGroupID), nil, data)
	if err != nil {
		return models.DiscoveryJobResponse{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeDiscoveryJob := models.DiscoveryJobResponse{}
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	return omeDiscoveryJob, err
}

func (c *Client) DeleteDiscoveryJob(discoveryGroupIds models.DiscoveryJobDeletePayload) error {
	data, _ := c.JSONMarshal(discoveryGroupIds)
	_, err := c.Post(DiscoveryJobRemoveAPI, nil, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetDiscoveryJobByGroupID(groupId int64) (models.DiscoveryJobResponse, error) {
	omeDiscoveryJob := models.DiscoveryJobResponse{}
	response, err := c.Get(fmt.Sprintf(DiscoveryJobByGroupIDAPI,groupId),nil,nil)
	if err != nil {
		return omeDiscoveryJob, err
	}
	respData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	return omeDiscoveryJob, err
}