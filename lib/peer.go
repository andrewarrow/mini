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
	ReadVersion(conn)
}

func ReadVersion(conn net.Conn) {
	inNetworkType, _ := ReadUvarint(conn)
	fmt.Println(inNetworkType)
	inMsgType, _ := ReadUvarint(conn)
	fmt.Println(inMsgType)
	checksum := make([]byte, 8)
	io.ReadFull(conn, checksum)
	payloadLength, _ := ReadUvarint(conn)
	payload := make([]byte, payloadLength)
	io.ReadFull(conn, payload)
	fmt.Println(len(payload))
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

	hdr := []byte{}
	hdr = append(hdr, UintToBuf(uint64(1))...)
	hdr = append(hdr, UintToBuf(uint64(1))...)
	payload := ToBytes(version)
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

func ToBytes(msg MsgBitCloutVersion) []byte {
	retBytes := []byte{}
	retBytes = append(retBytes, UintToBuf(msg.Version)...)
	retBytes = append(retBytes, UintToBuf(uint64(msg.Services))...)
	retBytes = append(retBytes, IntToBuf(msg.TstampSecs)...)
	retBytes = append(retBytes, UintToBuf(msg.Nonce)...)
	retBytes = append(retBytes, UintToBuf(uint64(len(msg.UserAgent)))...)
	retBytes = append(retBytes, msg.UserAgent...)
	retBytes = append(retBytes, UintToBuf(uint64(msg.StartBlockHeight))...)
	retBytes = append(retBytes, UintToBuf(uint64(msg.MinFeeRateNanosPerKB))...)
	retBytes = append(retBytes, UintToBuf(uint64(0))...)
	return retBytes
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
