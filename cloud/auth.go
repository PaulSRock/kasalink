package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const method = "login"
const authURL = "https://wap.tplinkcloud.com"

type resultFromAuth struct {
	RegTime string `json:"regTime"`
	Email   string `json:"email"`
	Token   string `json:"token"`
}

type authResponse struct {
	ErrorCode int            `json:"error_code"`
	Result    resultFromAuth `json:"result"`
	Message   string         `json:"msg"`
}

type paramsForAuth struct {
	AppType       string `json:"appType"`
	CloudUserName string `json:"cloudUserName"`
	CloudPassword string `json:"cloudPassword"`
	TerminalUUID  string `json:"terminalUUID"`
}

type auth struct {
	Method string        `json:"method"`
	Params paramsForAuth `json:"params"`
}

// GetCloudToken takes your username and password and returns a token for followup TP-Link API calls.
func GetCloudToken(user, pass string) (err error) {
	var (
		UUID uuid.UUID
		req  *http.Request
		resp *http.Response
		a    = auth{
			Method: method,
			Params: paramsForAuth{
				AppType:       "Kasa_Android",
				CloudUserName: user,
				CloudPassword: pass,
			},
		}
		jsonPayload []byte
		jsonDecoder *json.Decoder
		ar          authResponse
	)
	UUID = uuid.New()
	log.Printf("New UUID: %s\n", UUID)
	a.Params.TerminalUUID = fmt.Sprintf("%s", UUID)
	//a.Params.TerminalUUID = "bb24ac75-11ac-4faa-88e4-617abca4e53c"
	if jsonPayload, err = json.Marshal(a); err != nil {
		return
	}
	log.Printf("sending:\n%s\n", jsonPayload)
	if req, err = http.NewRequest("POST", authURL, bytes.NewReader(jsonPayload)); err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	if resp, err = theCloud.client.Do(req); err != nil {
		return
	}
	jsonDecoder = json.NewDecoder(resp.Body)
	if err = jsonDecoder.Decode(&ar); err != nil {
		return
	}
	defer closer(resp.Body)
	log.Printf("%+v", ar)
	if ar.ErrorCode != 0 {
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("error from tplink: %s; %s", b, err)
		return
	}
	theCloud.token = ar.Result.Token
	return
}

func closer(thingToClose io.Closer) {
	var err = thingToClose.Close()
	if err != nil {
		log.Println(err)
	}
}
