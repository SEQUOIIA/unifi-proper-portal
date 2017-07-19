package model

import (
	"crypto/sha512"
	"fmt"
	"time"
	"encoding/hex"
)

type Client struct {
	Name		string
	Email		string
	Id			string
}

func GenerateClientId(c Client) string{
	hasher := sha512.New()

	preHash := fmt.Sprintf("%s;%s;%d", c.Name, c.Email, time.Now().UnixNano())
	hasher.Write([]byte(preHash))
	//sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	sha := hex.EncodeToString(hasher.Sum(nil))

	return sha
}