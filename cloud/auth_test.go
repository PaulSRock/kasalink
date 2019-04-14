package cloud

import (
	"flag"
	"testing"
)

var (
	uid, pass string
)

func init() {
	flag.StringVar(&uid, "tpLinkUser", "", "User ID to test TPLink functionality")
	flag.StringVar(&pass, "tpLinkPass", "", "Password for the TPLink User")
	flag.Parse()
}

func TestGetCloudToken(t *testing.T) {
	if uid == "" || pass == "" {
		t.Skip("No user/pass given, so skipping.")
	}
	var err = GetCloudToken("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("My Token: %s", theCloud.token)
}

func TestGetDeviceList(t *testing.T) {
	if uid == "" || pass == "" {
		t.Skip("No user/pass given, so skipping.")
	}
	var deviceList, err = GetDeviceList()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", deviceList)
}
