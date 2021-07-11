package lib

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net"
	"time"
)

func Connect(ip net.IP) {
	fmt.Println("connecting to peer", ip)
	netAddr := net.TCPAddr{
		IP:   ip,
		Port: 17000,
	}
	fmt.Println(netAddr)
	conn, err := net.DialTimeout(netAddr.Network(), netAddr.String(), 30*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	SendVersion(conn)
}

func SendVersion(conn net.Conn) {
	version := MsgBitCloutVersion{}
	version.Version = 1
	version.Services = 1
	version.UserAgent = "Architect"
	version.Nonce = uint64(RandInt64(math.MaxInt64))
	version.TstampSecs = time.Now().Unix()
	version.StartBlockHeight = uint32(0)
	version.MinFeeRateNanosPerKB = 0

	fmt.Println(conn)
}

func RandInt64(max int64) int64 {
	val, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return val.Int64()
}
