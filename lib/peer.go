package lib

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
	"net"
	"time"

	merkletree "github.com/laser/go-merkle-tree"
)

var conn net.Conn

func Connect(ip net.IP) {
	fmt.Println("connecting to peer", ip)
	netAddr := net.TCPAddr{
		IP:   ip,
		Port: 17000,
	}
	fmt.Println(netAddr)
	var err error
	conn, err = net.DialTimeout(netAddr.Network(), netAddr.String(), 30*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	nonce1 := SendVersion()
	fmt.Println("nonce1", nonce1)
	m := ReadMessage()
	cv := m.(*MsgBitCloutVersion)
	SendNonce(cv.Nonce)
	m = ReadMessage()
	fmt.Println("nonce2", m)
	SendMempool()

	for {
		time.Sleep(time.Second * 1)
		fmt.Println("Reading...")
		ReadMessage()
	}
}

func ReadMessage() interface{} {
	inNetworkType, _ := ReadUvarint(conn)
	inMsgType, _ := ReadUvarint(conn)
	fmt.Println(inNetworkType, inMsgType)
	checksum := make([]byte, 8)
	io.ReadFull(conn, checksum)
	payloadLength, _ := ReadUvarint(conn)
	payload := make([]byte, payloadLength)
	io.ReadFull(conn, payload)

	var m interface{}
	if inMsgType == 1 {
		m = MsgBitCloutVersionFromBytes(payload)
	} else if inMsgType == 2 {
		m = MsgBitCloutVerackFromBytes(payload)
	} else if inMsgType == 10 {
		inv := MsgBitCloutInvFromBytes(payload)
		for _, item := range inv.InvList {
			fmt.Println(item.Type)
		}
	}
	return m
}

func SendVersion() uint64 {
	version := MsgBitCloutVersion{}
	version.Version = 1
	version.Services = 1
	version.UserAgent = "Architect"
	version.Nonce = uint64(RandInt64(math.MaxInt64))
	version.TstampSecs = time.Now().Unix()
	version.StartBlockHeight = uint32(0)
	version.MinFeeRateNanosPerKB = 0
	payload := version.ToBytes()
	SendPayloadWithType(1, payload)
	return version.Nonce
}
func SendNonce(n uint64) {
	m := MsgBitCloutVerack{}
	m.Nonce = n
	payload := m.ToBytes()
	SendPayloadWithType(2, payload)
}
func SendMempool() {
	payload := []byte{}
	SendPayloadWithType(14, payload)
}

func SendPayloadWithType(mType int, payload []byte) {
	hdr := []byte{}
	hdr = append(hdr, UintToBuf(uint64(1))...)
	hdr = append(hdr, UintToBuf(uint64(mType))...)
	hash := Sha256DoubleHash(payload)
	hdr = append(hdr, hash[:8]...)
	hdr = append(hdr, UintToBuf(uint64(len(payload)))...)
	_, err := conn.Write(hdr)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = conn.Write(payload)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Sha256DoubleHash(input []byte) *BlockHash {
	hashBytes := merkletree.Sha256DoubleHash(input)
	ret := &BlockHash{}
	copy(ret[:], hashBytes[:])
	return ret
}

func RandInt64(max int64) int64 {
	val, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return val.Int64()
}
