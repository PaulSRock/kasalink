package cloud

import "testing"

func TestGetCloudToken(t *testing.T) {
	var err = GetCloudToken("paulsrock@aol.com", "972YsLVPVP431E%OI8o9nYcE3G")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("My Token: %s", theCloud.token)
}

func TestGetDeviceList(t *testing.T) {
	var deviceList, err = GetDeviceList()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", deviceList)
}
