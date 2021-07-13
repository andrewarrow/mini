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

type MiniPeer struct {
	id   int
	conn net.Conn
}

func Connect(id int, ip net.IP) {
	fmt.Println("connecting to peer", ip)
	netAddr := net.TCPAddr{
		IP:   ip,
		Port: 17000,
	}
	fmt.Println(netAddr)
	var err error
	mp := MiniPeer{}
	mp.id = id
	mp.conn, err = net.DialTimeout(netAddr.Network(), netAddr.String(), 30*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	nonce1 := mp.SendVersion()
	fmt.Println("nonce1", nonce1)
	m := mp.ReadMessage()
	cv := m.(*MsgBitCloutVersion)
	mp.SendNonce(cv.Nonce)
	m = mp.ReadMessage()
	fmt.Println("nonce2", m)
	mp.SendMempool()

	for {
		time.Sleep(time.Second * 1)
		mp.ReadMessage()
	}
}

func (mp *MiniPeer) ReadMessage() interface{} {
	ReadUvarint(mp.conn)
	inMsgType, _ := ReadUvarint(mp.conn)
	//fmt.Println(inNetworkType, inMsgType)
	checksum := make([]byte, 8)
	io.ReadFull(mp.conn, checksum)
	payloadLength, _ := ReadUvarint(mp.conn)
	payload := make([]byte, payloadLength)
	io.ReadFull(mp.conn, payload)

	var m interface{}
	if inMsgType == 1 {
		m = MsgBitCloutVersionFromBytes(payload)
	} else if inMsgType == 2 {
		m = MsgBitCloutVerackFromBytes(payload)
	} else if inMsgType == 10 {
		inv := MsgBitCloutInvFromBytes(payload)
		t := MsgBitCloutGetTransactions{}
		for _, item := range inv.InvList {
			if item.Type != 0 {
				continue
			}
			t.HashList = append(t.HashList, &item.Hash)
		}
		payload := t.ToBytes()
		//fmt.Println("t.HashList", len(t.HashList))
		mp.SendPayloadWithType(12, payload)
	} else if inMsgType == 13 {
		MsgBitCloutTransactionBundleFromBytes(mp.id, payload)
	}
	return m
}

func (mp *MiniPeer) SendVersion() uint64 {
	version := MsgBitCloutVersion{}
	version.Version = 1
	version.Services = 1
	version.UserAgent = "Architect"
	version.Nonce = uint64(RandInt64(math.MaxInt64))
	version.TstampSecs = time.Now().Unix()
	version.StartBlockHeight = uint32(0)
	version.MinFeeRateNanosPerKB = 0
	payload := version.ToBytes()
	mp.SendPayloadWithType(1, payload)
	return version.Nonce
}
func (mp *MiniPeer) SendNonce(n uint64) {
	m := MsgBitCloutVerack{}
	m.Nonce = n
	payload := m.ToBytes()
	mp.SendPayloadWithType(2, payload)
}
func (mp *MiniPeer) SendMempool() {
	payload := []byte{}
	mp.SendPayloadWithType(14, payload)
}

func (mp *MiniPeer) SendPayloadWithType(mType int, payload []byte) {
	hdr := []byte{}
	hdr = append(hdr, UintToBuf(uint64(1))...)
	hdr = append(hdr, UintToBuf(uint64(mType))...)
	hash := Sha256DoubleHash(payload)
	hdr = append(hdr, hash[:8]...)
	hdr = append(hdr, UintToBuf(uint64(len(payload)))...)
	_, err := mp.conn.Write(hdr)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = mp.conn.Write(payload)
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
