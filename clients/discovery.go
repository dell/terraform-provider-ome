package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

func (c *Client) CreateDiscoveryJob(discoveryJob models.DiscoveryJobPayload) (models.DiscoveryJob, error) {
	data, _ := c.JSONMarshal(discoveryJob)
	response, err := c.Post(DiscoveryJobAPI, nil, data)
	if err != nil {
		return models.DiscoveryJob{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeDiscoveryJob := models.DiscoveryJob{}
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	return omeDiscoveryJob, err
}

func (c *Client) UpdateDiscoveryJob(discoveryJob models.DiscoveryJobPayload) (models.DiscoveryJob, error) {
	data, _ := c.JSONMarshal(discoveryJob)
	x := DiscoveryJobAPI + "?groupId=" + fmt.Sprint(discoveryJob.DiscoveryConfigGroupID)
	response, err := c.Post(x, nil, data)
	if err != nil {
		return models.DiscoveryJob{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeDiscoveryJob := models.DiscoveryJob{}
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	fmt.Println(string(respData))
	return omeDiscoveryJob, err
}

func (c *Client) DeleteDiscoveryJob(discoveryGroupIds models.DiscoveryJobDeletePayload) (string, error) {
	data, _ := c.JSONMarshal(discoveryGroupIds)
	resp, err := c.Post(DiscoveryJobRemoveAPI, nil, data)
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

func (c *Client) GetDiscoveryJobByGroupID(groupId int64) (models.DiscoveryJob, error) {
	omeDiscoveryJob := models.DiscoveryJob{}
	endpoint := fmt.Sprintf(DiscoveryJobByGroupIDAPI, groupId)
	// h := addHeaders()
	response, err := c.Get(endpoint, nil, nil)
	if err != nil {
		return omeDiscoveryJob, err
	}
	respData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	fmt.Println(string(respData))
	return omeDiscoveryJob, err
}

// func addHeaders() map[string]string{
// 	headers := make(map[string]string)
// 	headers["Content-Type"] = "application/json"
// 	headers["X-Auth-Token"] = "get your x-auth token"
// 	return headers
// }
