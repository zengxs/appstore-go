package apple

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"howett.net/plist"
)

func (c *AppleClient) Download(trackId, dest string) error {
	if c.Cred == nil {
		return ErrNotLoggedIn
	}

	payloadBuf := new(bytes.Buffer)
	payload := map[string]string{
		"creditDisplay": "",
		"guid":          c.Cred.GUID,
		"salableAdamId": trackId,
	}
	if err := plist.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return err
	}

	req := c.defaultRequest().
		SetQueryParams(map[string]string{
			"guid": c.Cred.GUID,
		}).
		SetHeader("Content-Type", "application/x-apple-plist").
		SetHeader("X-Dsid", c.Cred.DSID).
		SetHeader("iCloud-DSID", c.Cred.DSID).
		SetBody(payloadBuf.Bytes())
	resp, err := req.Post("https://p25-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/volumeStoreDownloadProduct")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.Join(ErrHTTPError, errors.New(resp.Status()))
	}

	var data map[string]any
	if err := plist.NewDecoder(bytes.NewReader(resp.Body())).Decode(&data); err != nil {
		return err
	}

	p, _ := json.Marshal(data)
	fmt.Println(string(p))
	return nil
}
