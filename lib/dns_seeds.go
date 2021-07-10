package lib

import (
	"fmt"
	"net"
)

var DNSSeeds = []string{
	"bitclout.coinbase.com",
	"bitclout.gemini.com",
	"bitclout.kraken.com",
	"bitclout.bitstamp.com",
	"bitclout.bitfinex.com",
	"bitclout.binance.com",
	"bitclout.hbg.com",
	"bitclout.okex.com",
	"bitclout.bithumb.com",
	"bitclout.upbit.com",
	"bitclout-seed-0.io",
	"bitclout-seed-1.io",
	"bitclout-seed-2.io",
	"bitclout-seed-3.io",
	"bitclout-seed-4.io",
	"bitclout-seed-5.io",
	"bitclout-seed-6.io",
	"bitclout-seed-7.io",
	"bitclout-seed-8.io",
	"bitclout-seed-9.io",
	"bitclout-seed-10.io",
	"bitclout-seed-11.io",
	"bitclout-seed-12.io",
	"bitclout-seed-13.io",
	"bitclout-seed-14.io",
	"bitclout-seed-15.io",
	"bitclout-seed-16.io",
	"bitclout-seed-17.io",
	"bitclout-seed-18.io",
	"bitclout-seed-19.io",
}

func IPsForHost(host string) []net.IP {
	items := []net.IP{}
	ipAddrs, err := net.LookupIP(host)
	if err != nil {
		fmt.Println(err)
		return items
	}
	for _, ip := range ipAddrs {
		fmt.Println(ip)
		items = append(items, ip)
	}
	return items
}

func GatherValidIPs() []net.IP {
	items := []net.IP{}
	for _, seed := range DNSSeeds {
		items = append(items, IPsForHost(seed)...)
	}
	return items
}
