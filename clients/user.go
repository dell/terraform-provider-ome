package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

// CreateUser - to create ome user
func (c *Client) CreateUser(user models.UserPayload) (models.User, error) {
	data, _ := c.JSONMarshal(user)
	response, err := c.Post(UserAPI, nil, data)
	if err != nil {
		return models.User{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeUser := models.User{}
	err = c.JSONUnMarshal(respData, &omeUser)
	return omeUser, err
}

// UpdateUser - to update ome user
func (c *Client) UpdateUser(user models.User) (models.User, error) {
	data, _ := c.JSONMarshal(user)
	x := UserAPI + fmt.Sprintf("('%s')", user.ID)
	response, err := c.Put(x, nil, data)
	if err != nil {
		return models.User{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeUser := models.User{}
	err = c.JSONUnMarshal(respData, &omeUser)
	fmt.Println(string(respData))
	return omeUser, err
}

// DeleteUser - to delete ome user
func (c *Client) DeleteUser(id string) (string, error) {
	endpoint := fmt.Sprintf(UserAPI+"('%s')", id)
	resp, err := c.Delete(endpoint, nil, nil)
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

// GetUserByID - to get user by id
func (c *Client) GetUserByID(id string) (models.User, error) {
	omeUser := models.User{}
	endpoint := fmt.Sprintf(UserAPI+"('%s')", id)
	response, err := c.Get(endpoint, nil, nil)
	if err != nil {
		return omeUser, err
	}
	respData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(respData, &omeUser)
	fmt.Println(string(respData))
	return omeUser, err
}
