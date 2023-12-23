package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"syscall/js"

	"github.com/jonfriesen/geoip-wasm/assets"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("getIPLocation", js.FuncOf(lookupIP))

	<-c
}

func lookupIP(this js.Value, args []js.Value) interface{} {
	if len(args) <= 0 {
		return fmt.Sprintf("missing ip arguments")
	}
	stringIP := args[0].String()

	ipLoc, err := getIPLocation(stringIP)
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}

	ipLocB, err := json.MarshalIndent(ipLoc, "", "\t")
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}

	jsData := js.Global().Get("Uint8Array").New(len(ipLocB))
	js.CopyBytesToJS(jsData, ipLocB)

	return jsData
}

type IPLocation struct {
	CountryCode      string `json:"country_code"`
	Country          string `json:"country"`
	Region           string `json:"region"`
	City             string `json:"city"`
	Timezone         string `json:"timezone"`
	IsAnonymousProxy bool   `json:"is_anonymous_proxy"`
}

func getIPLocation(ipstr string) (*IPLocation, error) {
	db, err := geoip2.FromBytes(assets.GeoLite2City)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(ipstr)
	if ip == nil {
		return nil, fmt.Errorf("bad ip")
	}

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	city := &IPLocation{
		CountryCode:      record.Country.IsoCode,
		Country:          record.Country.Names["en"],
		City:             record.City.Names["en"],
		Timezone:         record.Location.TimeZone,
		IsAnonymousProxy: record.Traits.IsAnonymousProxy,
	}
	if len(record.Subdivisions) > 0 {
		city.Region = record.Subdivisions[0].Names["en"]
	}

	return city, nil
}
