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
}

func IPsForHost(host string) {
	ipAddrs, err := net.LookupIP(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, ip := range ipAddrs {
		fmt.Println(ip)
	}
}
