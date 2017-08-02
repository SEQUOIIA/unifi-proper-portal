package unifi

import (
	"fmt"
	"github.com/sequoiia/unifi-proper-portal/model"
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

	_, err = c.doRequest(req, nil, true)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UnauthoriseGuest(data model.UniFiGuestUnauthoriseRequest) error {
	var payload struct {
		Command string `json:"cmd"`
		Mac     string `json:"mac"`
	}
	payload.Command = "unauthorize-guest"
	payload.Mac = data.Mac

	req, err := c.createRequest(http.MethodPost, fmt.Sprintf("/api/s/%s/cmd/stamgr", c.Config.Site), payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, nil, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetVouchers() ([]model.UniFiVoucherResponse, error) {
	req, err := c.createRequest(http.MethodGet, fmt.Sprintf("/api/s/%s/stat/voucher", c.Config.Site), nil)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Data []model.UniFiVoucherResponse `json:"data"`
	}

	_, err = c.doRequest(req, &payload, true)
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}

func (c *Client) RemoveVoucher(voucherId string) error {
	var payload struct {
		Command string `json:"cmd"`
		Id      string `json:"_id"`
	}
	payload.Command = "delete-voucher"
	payload.Id = voucherId

	req, err := c.createRequest(http.MethodPost, fmt.Sprintf("/api/s/%s/cmd/hotspot", c.Config.Site), payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, nil, true)
	if err != nil {
		return err
	}
	return nil
}
