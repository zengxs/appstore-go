package apple

import "github.com/go-resty/resty/v2"

type AppleClient struct {
	Cred *AppleCredentials
}

func NewAppleClient() *AppleClient {
	return &AppleClient{
		Cred: nil,
	}
}

func NewAppleClientWithCred(cred *AppleCredentials) *AppleClient {
	return &AppleClient{
		Cred: cred,
	}
}

func (c *AppleClient) defaultRequest() *resty.Request {
	req := resty.New().R()
	req.SetHeader("User-Agent", "Configurator/2.15 (Macintosh; OS X 11.0.0; 16G29) AppleWebKit/2603.3.8")
	if c.Cred != nil {
		req.SetHeader("X-Dsid", c.Cred.DSID)
		req.SetHeader("iCloud-DSID", c.Cred.DSID)
		req.SetCookies(c.Cred.Cookies)
	}
	return req
}
