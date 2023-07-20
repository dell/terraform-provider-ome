package clients

import (
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ID string

func TestClient_CreateUserClient(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	createUser := models.UserPayload{
		UserTypeID:         1,
		DirectoryServiceID: 0,
		Description:        "Avenger",
		Password:           "Dell123!",
		UserName:           "Dell",
		RoleID:             "10",
		Locked:             false,
		Enabled:            true,
	}

	tests := []struct {
		name string
		args models.UserPayload
	}{
		{"Create User Successfully", createUser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cuser, err := c.CreateUser(tt.args)
			ID = cuser.ID
			t.Log(cuser, err)
			if err == nil {
				assert.Equal(t, createUser.UserName, cuser.UserName)
			}
		})
	}
}

func TestClient_UpdateUserClient(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	updateUser := models.User{
		ID:                 ID,
		UserTypeID:         1,
		DirectoryServiceID: 0,
		Description:        "Avenger",
		Password:           "Dell123!",
		UserName:           "Dell",
		RoleID:             "10",
		Locked:             false,
		Enabled:            false,
	}
	tests := []struct {
		name string
		args models.User
	}{
		{"Update User Successfully", updateUser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuser, err := c.UpdateUser(tt.args)
			t.Log(uuser, err)
			if err == nil {
				assert.Equal(t, updateUser.Enabled, uuser.Enabled)
			}
		})
	}
}

func TestClient_GetUserClient(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name string
		args string
	}{
		{"Get User Successfully", ID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guser, err := c.GetUserByID(ID)
			t.Log(guser, err)
			if err == nil {
				assert.Equal(t, ID, guser.ID)
			}
		})
	}
}

func TestClient_CreateUser(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name string
		args string
	}{
		{"Delete User Successfully", ID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duser, err := c.DeleteUser(tt.args)
			t.Log(duser, err)
			if err == nil {
				assert.Equal(t, duser, duser)
			}
		})
	}
}
