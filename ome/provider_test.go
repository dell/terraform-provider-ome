package ome

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	//SkipTestMsg
	SkipTestMsg = "Skipping the test because eith TF_ACC or ACC_DETAIL is not set to 1"
)

var testProvider tfsdk.Provider
var testProviderFactory map[string]func() (tfprotov6.ProviderServer, error)
var omeUserName = os.Getenv("OME_USERNAME")
var omeHost = os.Getenv("OME_HOST")
var omePassword = os.Getenv("OME_PASSWORD")
var DeviceSvcTag1 = os.Getenv("DEVICESVCTAG1")
var DeviceSvcTag2 = os.Getenv("DEVICESVCTAG2")
var DeviceID1 = os.Getenv("DEVICEID1")
var DeviceID2 = os.Getenv("DEVICEID2")
var DeviceID3 = os.Getenv("DEVICEID3") // Not capable for deployment
var ShareUser = os.Getenv("SHAREUSERNAME")
var SharePassword = os.Getenv("SHAREPASSWORD")
var ShareIP = os.Getenv("SHAREIP")

func init() {
	testProvider = New("test")()
	testProviderFactory = map[string]func() (tfprotov6.ProviderServer, error){
		// newProvider is an example function that returns a tfsdk.Provider
		"ome": providerserver.NewProtocol6WithError(testProvider),
	}

}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OME_USERNAME"); v == "" {
		t.Fatal("OME_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("OME_PASSWORD"); v == "" {
		t.Fatal("OME_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("OME_HOST"); v == "" {
		t.Fatal("OME_HOST must be set for acceptance tests")
	}

	testProvider.Configure(context.Background(), tfsdk.ConfigureProviderRequest{}, &tfsdk.ConfigureProviderResponse{})

}

func skipTest() bool {
	return os.Getenv("TF_ACC") == "" || os.Getenv("ACC_DETAIL") == ""
}
