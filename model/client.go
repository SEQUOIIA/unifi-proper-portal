package model

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"
)

type Client struct {
	Name       string
	Email      string
	Id         string
	Authorised uint8 // 0 = none, 1 = false, 2 = true
	AuthedBy   AuthMethod
	Tokens     Tokens
	Device     string
	Voucher    string
}

type AuthMethod uint8

var AuthedByFacebook AuthMethod = 1
var AuthedByVoucher AuthMethod = 0

type Tokens struct {
	Facebook *OAuth_accessTokenResponse
}

func GenerateClientId(c Client) string {
	hasher := sha512.New()

	preHash := fmt.Sprintf("%s;%s;%d", c.Name, c.Email, time.Now().UnixNano())
	hasher.Write([]byte(preHash))
	//sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	sha := hex.EncodeToString(hasher.Sum(nil))

	return sha
}

func GenerateClientIdVoucher(c Client) string {
	hasher := sha512.New()

	preHash := fmt.Sprintf("%s;%d", c.Device, time.Now().UnixNano())
	hasher.Write([]byte(preHash))
	//sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	sha := hex.EncodeToString(hasher.Sum(nil))

	return sha
}
