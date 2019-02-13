package cloud

import "net/http"

type tpLinkCloud struct {
	client  *http.Client
	token   string
	devices *[]tpLinkDevice
}

var theCloud *tpLinkCloud

func init() {
	theCloud = &tpLinkCloud{
		client: http.DefaultClient,
	}
}
