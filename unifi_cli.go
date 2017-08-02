package main

import (
	"github.com/sequoiia/unifi-proper-portal/controller"
	"github.com/sequoiia/unifi-proper-portal/unifi"
	"log"
)

func main() {
	controller.LoadConfig()

	cli := unifi.NewClient(&unifi.ClientConfig{
		Username:        controller.Config.UniFi.Username,
		Password:        controller.Config.UniFi.Password,
		UniFiBaseUrlRaw: controller.Config.UniFi.BaseUrl,
		Site:            controller.Config.UniFi.Site,
		InsecureMode:    controller.Config.UniFi.InsecureMode,
	})

	err := cli.Login()
	if err != nil {
		log.Fatal(err)
	}

	vouchers, err := cli.GetVouchers()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(vouchers); i++ {
		log.Println(vouchers[i].Code)
	}
}
