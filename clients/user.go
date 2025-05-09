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

// CreateUser - to create ome user
func (c *Client) CreateUser(user models.UserPayload) (models.User, error) {
	omeUser := models.User{}
	data, errMarshal := c.JSONMarshal(user)
	if errMarshal != nil {
		return omeUser, errMarshal
	}
	response, err := c.Post(UserAPI, nil, data)
	if err != nil {
		return omeUser, err
	}
	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return omeUser, getBodyError
	}
	err = c.JSONUnMarshal(respData, &omeUser)
	return omeUser, err
}

// UpdateUser - to update ome user
func (c *Client) UpdateUser(user models.User) (models.User, error) {
	omeUser := models.User{}
	data, errMarshal := c.JSONMarshal(user)
	if errMarshal != nil {
		return omeUser, errMarshal
	}
	x := UserAPI + fmt.Sprintf("('%s')", user.ID)
	response, err := c.Put(x, nil, data)
	if err != nil {
		return models.User{}, err
	}
	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return omeUser, getBodyError
	}
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
	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return omeUser, getBodyError
	}
	err = c.JSONUnMarshal(respData, &omeUser)
	fmt.Println(string(respData))
	return omeUser, err
}
