package clients

import "terraform-provider-ome/models"

// GetAllVlanNetworks returns the vlan data from OME
func (c *Client) GetAllVlanNetworks() ([]models.VLanNetworks, error) {
	vlanData := []models.VLanNetworks{}
	err := c.GetPaginatedData(VlanNetworksAPI, &vlanData)
	if err != nil {
		return []models.VLanNetworks{}, err
	}
	return vlanData, nil
}
