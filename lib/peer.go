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
	SendVersion()
	m := ReadMessage()
	cv := m.(*MsgBitCloutVersion)
	fmt.Println("nonce1", cv.Nonce)
	SendNonce(cv.Nonce)
	m = ReadMessage()
	fmt.Println("nonce2", m)

	/*
		for {
			time.Sleep(time.Second * 1)
			fmt.Println("Reading...")
			inNetworkType, _ := ReadUvarint(conn)
			fmt.Println(inNetworkType)
			inMsgType, _ := ReadUvarint(conn)
			fmt.Println(inMsgType)
		}*/
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
	}
	fmt.Println(m)
	return m
}

func SendVersion() {
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
}
func SendNonce(n uint64) {
	m := MsgBitCloutVerack{}
	m.Nonce = n
	payload := m.ToBytes()
	SendPayloadWithType(2, payload)
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

type BlockHash [32]byte

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
