package cloud

import "testing"

func TestGetCloudToken(t *testing.T) {
	var err = GetCloudToken("", "")
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
