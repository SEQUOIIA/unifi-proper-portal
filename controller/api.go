package controller

import (
	"bytes"
	"encoding/json"
	"github.com/sequoiia/unifi-proper-portal/model"
	"log"
	"net/http"
	"strings"
)

func ApiStatus(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Error        int    `json:"error"`
		ErrorMessage string `json:"message"`
		Data         struct {
			Authorisation uint8 `json:"authorisation"`
		} `json:"data"`
	}
	payload.Error = 0

	id, err := model.GetCookieValue(r, "UPP_ID")
	if err != nil {
		payload.Error = 1
		payload.ErrorMessage = "User couldn't be found."
	} else {
		if user, ok := Users[id]; ok {
			payload.Data.Authorisation = user.Authorised
		} else {
			payload.Error = 1
			payload.ErrorMessage = "User couldn't be found."
		}
	}

	buf := bytes.NewBuffer(nil)

	err = json.NewEncoder(buf).Encode(payload)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(200)
	w.Write(buf.Bytes())
}

func ApiVoucherConsume(w http.ResponseWriter, r *http.Request) {
	UniFiClientInit()
	var requestPayload struct {
		Code string `json:"code"`
	}

	var responsePayload struct {
		Error        int    `json:"error"`
		ErrorMessage string `json:"error_message"`
		Success      int    `json:"success"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		responsePayload.Error = 1
		responsePayload.ErrorMessage = "Invalid voucher request"
		responsePayload.Success = 0
		err = json.NewEncoder(w).Encode(responsePayload)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	vouchers, err := Uclient.GetVouchers()
	code := strings.Replace(requestPayload.Code, "-", "", -1)
	var authed bool = false
	client, err := model.GetUniFiGuestCookies(r)
	if err != nil {
		responsePayload.Error = 2
		responsePayload.ErrorMessage = "Invalid voucher request"
		responsePayload.Success = 0
		err = json.NewEncoder(w).Encode(responsePayload)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		for i := 0; i < len(vouchers); i++ {
			voucher := vouchers[i]
			if strings.Compare(code, voucher.Code) == 0 {
				// Auth client
				authRequest := model.UniFiGuestAuthoriseRequest{
					Mac:     client.ClientMacAddress,
					Minutes: voucher.Duration,
				}

				if voucher.QosOverwrite {
					authRequest.Up = voucher.QosRateMaxUp
					authRequest.Down = voucher.QosRateMaxDown
					authRequest.Bytes = voucher.QosUsageQuota
				}

				err := Uclient.AuthoriseGuest(authRequest)
				if err != nil {
					responsePayload.Error = 1
					responsePayload.ErrorMessage = "Unknown error. Contact network administrator."
					responsePayload.Success = 0
					err = json.NewEncoder(w).Encode(responsePayload)
					if err != nil {
						log.Fatal(err)
					}
				}
				authed = true
				responsePayload.Success = 1

				if voucher.Quota == 1 {
					err := Uclient.RemoveVoucher(voucher.Id)
					if err != nil {
						log.Fatal(err)
					}
				}

				iClient := model.Client{
					Authorised: 2,
					AuthedBy:   model.AuthedByVoucher,
					Device:     client.ClientMacAddress,
					Voucher:    voucher.Code,
				}
				iClient.Id = model.GenerateClientIdVoucher(iClient)
				Users[iClient.Id] = iClient

				NewUserCookie(w, iClient.Id)

				break
			}
		}
	}

	if !authed {
		responsePayload.Error = 1
		responsePayload.ErrorMessage = "Invalid voucher code"
		responsePayload.Success = 0
		err = json.NewEncoder(w).Encode(responsePayload)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = json.NewEncoder(w).Encode(responsePayload)
		if err != nil {
			log.Fatal(err)
		}
	}
}
