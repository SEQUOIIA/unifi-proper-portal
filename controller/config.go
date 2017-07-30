package controller

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"net"
)

type ConfigS struct {
	ClientId     string
	ClientSecret string
	ProxyMode    bool
	Domain       string
	UniFi        struct {
		Username     string
		Password     string
		BaseUrl      string
		Site         string
		InsecureMode bool
	}
	AllowedSubnets []Subnet
}

type Subnet struct {
	Ip    net.IP
	IpNet *net.IPNet
	Count uint8
}

var Config ConfigS

func AddAllowedSubnet(s string) error {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}

	subnet := Subnet{}
	subnet.Ip = ip
	subnet.IpNet = ipnet
	count, err := subnetCount4(subnet.Ip)
	if err != nil {
		return err
	}
	subnet.Count = count
	Config.AllowedSubnets = append(Config.AllowedSubnets, subnet)

	return nil
}

func subnetCount4(ip net.IP) (uint8, error) {
	var count uint8
	if p4 := ip.To4(); len(p4) == net.IPv4len {
		for i := 1; i < len(p4); i++ {
			if uint(p4[i]) != 0 {
				count++
			}
		}
	} else {
		return count, errors.New("Invalid IPv4 address")
	}
	return count, nil
}

func LoadConfig() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	Config = ConfigS{}
	Config.ClientId = viper.GetString("clientid")
	Config.ClientSecret = viper.GetString("clientsecret")
	Config.ProxyMode = viper.GetBool("proxymode")
	Config.Domain = viper.GetString("domain")
	unifiVars := viper.GetStringMap("unifi")
	Config.UniFi.Username = unifiVars["username"].(string)
	Config.UniFi.Password = unifiVars["password"].(string)
	Config.UniFi.BaseUrl = unifiVars["baseurl"].(string)
	Config.UniFi.Site = unifiVars["site"].(string)
	Config.UniFi.InsecureMode = unifiVars["insecuremode"].(bool)
	allowedSubnets := viper.GetStringSlice("allowedsubnets")

	for i := 0; i < len(allowedSubnets); i++ {
		AddAllowedSubnet(allowedSubnets[i])
	}
}
