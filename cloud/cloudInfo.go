package cloud

import "net/http"

type tpLinkCloud struct {
	client  *http.Client
	token   string
	devices *[]TPLinkDevice
}

var theCloud *tpLinkCloud

func init() {
	theCloud = &tpLinkCloud{
		client: http.DefaultClient,
	}
}
