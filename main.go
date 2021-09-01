package main

import (
	"cf_ddns/cf"
	"flag"
	"log"
	"time"
)

var key = flag.String("auth_key", "abcdefg", "auth key")
var email = flag.String("auth_email", "xxx@gmail.com", "auth email")
var zoneName = flag.String("zone_name", "abc.com", "zone name")
var recordName = flag.String("record_name", "ddns.abc.com", "record_name")
var ipProvider = flag.String("ip_provider", "https://ip.sendev.cc", "ip provider")

func main() {
	if len(*key) == 0 || len(*email) == 0 || len(*zoneName) == 0 || len(*recordName) == 0 || len(*ipProvider) == 0 {
		log.Fatal("invalid params")
	}

	client := cf.New(
		cf.WithAuthKey(*key),
		cf.WithAuthMail(*email),
		cf.WithIPProvider(*ipProvider),
		cf.WithZoneName(*zoneName),
		cf.WithRecName(*recordName),
	)
	client.StartRefresh(time.Minute)
}
