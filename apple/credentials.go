package apple

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/afero"
	"howett.net/plist"
)

type AppleCredentials struct {
	// AppleID is the Apple ID to login, usually an email address
	AppleID string `json:"apple_id"`
	// Password is the password for the Apple ID
	Password string `json:"password"`
	// PasswordToken is the access token for Apple services, this is returned after a successful login
	PasswordToken string `json:"password_token"`
	// DSID represents the user's Apple ID, this is returned after a successful login
	DSID string `json:"dsid"`
	// Region is the region of the Apple ID, ISO 3166-1 alpha-2 country code
	Region string `json:"region"`
	// GUID is the MAC address of the device, without colons
	GUID string `json:"guid"`
	// Cookies is the cookies for the Apple services
	Cookies []*http.Cookie `json:"cookies"`
}

type LoginOptions struct {
	AppleID    string
	Password   string
	MacAddress string // MAC address of the device, without colons. If empty, it will use the device's actual MAC address
	Region     string
}

// Login logs in to Apple services with the Apple ID and password
func (c *AppleClient) Login(opt LoginOptions) error {
	if c.Cred != nil {
		return ErrAlreadyLoggedIn
	}

	if opt.MacAddress == "" {
		mac, err := getDeviceMAC()
		if err != nil {
			return err
		}
		opt.MacAddress = mac
	}
	guid := strings.ToLower(strings.ReplaceAll(opt.MacAddress, ":", ""))

	if opt.Region == "" {
		opt.Region = "US"
	}

	// marshal the login payload
	payload := map[string]string{
		"appleId":       opt.AppleID,
		"attempt":       "4",
		"password":      opt.Password,
		"createSession": "true",
		"guid":          guid,
		"rmp":           "0",
		"why":           "signIn",
	}
	buf := new(bytes.Buffer)
	if err := plist.NewEncoder(buf).Encode(payload); err != nil {
		return err
	}

	// login request
	req := c.defaultRequest().
		SetHeader("Content-Type", "application/x-apple-plist").
		SetQueryParams(map[string]string{
			"guid": guid,
			"Pod":  "22",
			"PRH":  "22",
		}).
		SetBody(buf.Bytes())
	resp, err := req.Post("https://p22-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate")
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.Join(ErrHTTPError, errors.New(resp.Status()))
	}

	// unmarshal the login response
	var data map[string]any
	if err := plist.NewDecoder(bytes.NewReader(resp.Body())).Decode(&data); err != nil {
		return err
	}
	cred := AppleCredentials{
		AppleID:       data["accountInfo"].(map[string]any)["appleId"].(string),
		PasswordToken: data["passwordToken"].(string),
		Password:      opt.Password,
		DSID:          data["dsPersonId"].(string),
		Region:        opt.Region,
		GUID:          guid,
		Cookies:       resp.Cookies(),
	}

	c.Cred = &cred

	return nil
}

// LoadCredentials loads the Apple credentials from a file
func (c *AppleClient) LoadCredentials(credentialPath string) error {
	cred := AppleCredentials{}
	data, err := afero.ReadFile(afero.NewOsFs(), credentialPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &cred); err != nil {
		return err
	}

	c.Cred = &cred
	return nil
}

// SaveCredentials saves the Apple credentials to a file
func (c *AppleClient) SaveCredentials(credentialPath string) error {
	if c.Cred == nil {
		return ErrNotLoggedIn
	}

	data, err := json.Marshal(c.Cred)
	if err != nil {
		return err
	}

	return afero.WriteFile(afero.NewOsFs(), credentialPath, data, 0600)
}

// StoreFront returns the Apple StoreFront code for the region
func (c *AppleCredentials) StoreFront() string {
	code, ok := storefrontCodes[c.Region]
	if !ok {
		return storefrontCodes["US"] // default to US
	}
	return code
}

func getDeviceMAC() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && !bytes.Equal(iface.HardwareAddr, nil) {
			// Skip locally administered addresses
			if iface.HardwareAddr[0]&2 == 2 {
				continue
			}
			return iface.HardwareAddr.String(), nil
		}
	}

	return "3c:06:30:0f:0f:0f", nil // if no MAC address is found, return a fake one, this should never happen
}
