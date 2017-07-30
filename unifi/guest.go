package unifi

import (
	"fmt"
	"github.com/sequoiia/unifi-proper-portal/model"
	"io/ioutil"
	"net/http"
)

func (c *Client) AuthoriseGuest(data model.UniFiGuestAuthoriseRequest) error {
	var payload struct {
		Command string   `json:"cmd"`
		Mac     string   `json:"mac"`
		Minutes float64  `json:"minutes"`
		Up      *float64 `json:"up,omitempty"`
		Down    *float64 `json:"down,omitempty"`
		Bytes   *float64 `json:"bytes,omitempty"`
	}

	payload.Command = "authorize-guest"
	payload.Mac = data.Mac
	payload.Minutes = data.Minutes
	payload.Up = data.Up
	payload.Down = data.Down
	payload.Bytes = data.Bytes

	req, err := c.createRequest(http.MethodPost, fmt.Sprintf("/api/s/%s/cmd/stamgr", c.Config.Site), payload)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(req, nil, false)
	if err != nil {
		return err
	}

	tmpbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(tmpbody))

	return nil
}

