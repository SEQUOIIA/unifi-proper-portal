package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/sequoiia/unifi-proper-portal/model"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const USERAGENT string = "github.com/sequoiia/unifi-proper-portal - @theSequoiia"

type Client struct {
	Config *ClientConfig

	hclient *http.Client
}

type ClientConfig struct {
	Username        string
	Password        string
	UniFiBaseUrlRaw string
	UniFiBaseUrl    *url.URL
	Site            string
	InsecureMode    bool
}

func NewClient(cc *ClientConfig) *Client {
	var c *Client
	if cc != nil {
		c = &Client{Config: cc}
	} else { // If no config is provided, use default values.
		c = &Client{Config: &ClientConfig{
			Username:        "ubnt",
			Password:        "ubnt",
			UniFiBaseUrlRaw: "https://unifi:8443",
			Site:            "default",
			InsecureMode:    false},
		}
	}

	if len(c.Config.Site) == 0 {
		c.Config.Site = "default"
	}

	c.hclient = http.DefaultClient

	var err error
	c.hclient.Jar, err = cookiejar.New(nil)
	if err != nil {
		log.Println(err)
	}

	c.Config.UniFiBaseUrl, err = url.Parse(c.Config.UniFiBaseUrlRaw)
	if err != nil {
		log.Println(err)
	}

	if c.Config.InsecureMode {
		// Don't verify the TLS cert, use only for internal UniFi controllers if at all possible.
		c.hclient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return c
}

// If body is not nil and the HTTP method is POST, assume the input can be converted to JSON.
func (c *Client) createRequest(httpmethod string, path string, body interface{}) (*http.Request, error) {
	buf := bytes.NewBuffer(nil)

	if httpmethod == http.MethodPost && body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(httpmethod, c.Config.UniFiBaseUrlRaw+path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", USERAGENT)

	return req, nil
}

func (c *Client) doRequest(req *http.Request, t interface{}, deferBody bool) (*http.Response, error) {
	resp, err := c.hclient.Do(req)
	if err != nil {
		return resp, err
	}

	if deferBody {
		defer resp.Body.Close()
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, errors.New("Response was not in the 200 status code range, investigate.")
	}

	c.hclient.Jar.SetCookies(c.Config.UniFiBaseUrl, resp.Cookies())

	if t == nil {
		return resp, nil
	}

	return resp, json.NewDecoder(resp.Body).Decode(t)
}

func (c *Client) Login() error {
	loginBody := model.UniFiLoginRequest{Username: c.Config.Username, Password: c.Config.Password}

	req, err := c.createRequest(http.MethodPost, "/api/login", loginBody)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, nil, true)
	if err != nil {
		return err
	}

	/*
		tmpbody, err := ioutil.ReadAll(resp.Body)

		log.Println(string(tmpbody))

		for i := 0; i < len(resp.Cookies()); i++ {
			cookie := resp.Cookies()[i]
			fmt.Printf("%s : %s\n", cookie.Name, cookie.Value)
			fmt.Println(cookie.String())
		}

	*/
	return nil
}
