package ewdxem4dante

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

var rest *resty.Client

func createRestyClient(password string) {
	rest = resty.New().
		SetBasicAuth("api", password).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetRetryCount(0).
		SetTimeout(time.Second * 5)
}
