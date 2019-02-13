package cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const discoveryURI = "https://wap.tplinkcloud.com?token=%s"

type tpLinkDeviceList struct {
	DeviceList []tpLinkDevice `json:"deviceList"`
}

type tpLinkGetDeviceListResponse struct {
	Result    tpLinkDeviceList `json:"result"`
	ErrorCode int              `json:"error_code"`
	Message   string           `json:"msg"`
}

// GetDeviceList takes an auth Token and returns the list of devices registered to the Token's Account
func GetDeviceList() (devices []tpLinkDevice, err error) {
	var (
		req         *http.Request
		resp        *http.Response
		jsonDecoder *json.Decoder
		uri         string
		payload     tpLinkGetDeviceListResponse
	)
	uri = fmt.Sprintf(discoveryURI, theCloud.token)
	if req, err = http.NewRequest("POST", uri, strings.NewReader(`{"method":"getDeviceList"}`)); err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	if resp, err = theCloud.client.Do(req); err != nil {
		return
	}
	defer closer(req.Body)
	jsonDecoder = json.NewDecoder(resp.Body)
	if err = jsonDecoder.Decode(&payload); err != nil {
		return
	}
	//log.Printf("%+v", payload)
	devices = payload.Result.DeviceList
	return
}
