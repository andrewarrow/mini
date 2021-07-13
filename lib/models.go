package lib

import (
	"bytes"
	"io"
)

type BlockHash [32]byte

type InvVect struct {
	Type uint32 // 0 tx
	Hash BlockHash
}

type MsgBitCloutInv struct {
	InvList        []*InvVect
	IsSyncResponse bool
}

type MsgBitCloutVersion struct {
	Version              uint64
	Services             uint64
	TstampSecs           int64
	Nonce                uint64
	UserAgent            string
	StartBlockHeight     uint32
	MinFeeRateNanosPerKB uint64
}

func (msg *MsgBitCloutVersion) ToBytes() []byte {
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

func MsgBitCloutVersionFromBytes(data []byte) *MsgBitCloutVersion {
	rr := bytes.NewReader(data)
	m := MsgBitCloutVersion{}

	ver, _ := ReadUvarint(rr)
	m.Version = ver
	services, _ := ReadUvarint(rr)
	m.Services = services
	tstampSecs, _ := ReadVarint(rr)
	m.TstampSecs = tstampSecs

	nonce, _ := ReadUvarint(rr)
	m.Nonce = nonce
	strLen, _ := ReadUvarint(rr)
	userAgent := make([]byte, strLen)
	io.ReadFull(rr, userAgent)
	m.UserAgent = string(userAgent)

	lastBlockHeight, _ := ReadUvarint(rr)
	m.StartBlockHeight = uint32(lastBlockHeight)

	minFeeRateNanosPerKB, _ := ReadUvarint(rr)
	m.MinFeeRateNanosPerKB = minFeeRateNanosPerKB
	//ReadUvarint(rr)
	return &m
}

type MsgBitCloutVerack struct {
	Nonce uint64
}

func (msg *MsgBitCloutVerack) ToBytes() []byte {
	retBytes := []byte{}
	retBytes = append(retBytes, UintToBuf(msg.Nonce)...)
	return retBytes
}

func MsgBitCloutVerackFromBytes(data []byte) *MsgBitCloutVerack {
	rr := bytes.NewReader(data)
	m := MsgBitCloutVerack{}

	nonce, _ := ReadUvarint(rr)
	m.Nonce = nonce
	return &m
}

func _readInvList(rr io.Reader) ([]*InvVect, error) {
	invList := []*InvVect{}
	numInvVects, _ := ReadUvarint(rr)
	for ii := uint64(0); ii < numInvVects; ii++ {
		typeUint, _ := ReadUvarint(rr)

		invHash := BlockHash{}
		io.ReadFull(rr, invHash[:])

		invVect := &InvVect{
			Type: uint32(typeUint),
			Hash: invHash,
		}
		invList = append(invList, invVect)
	}

	return invList, nil
}

func MsgBitCloutInvFromBytes(data []byte) *MsgBitCloutInv {
	rr := bytes.NewReader(data)
	invList, _ := _readInvList(rr)
	//isSyncResponse := _readBoolByte(rr)

	m := MsgBitCloutInv{
		InvList: invList,
	}
	return &m
}
