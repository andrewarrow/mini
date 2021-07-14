package lib

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/btcsuite/btcd/btcec"
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

type MsgBitCloutGetTransactions struct {
	HashList []*BlockHash
}

func (msg *MsgBitCloutGetTransactions) ToBytes() []byte {
	data := []byte{}

	data = append(data, UintToBuf(uint64(len(msg.HashList)))...)
	for _, hash := range msg.HashList {
		data = append(data, hash[:]...)
	}

	return data
}

type BitCloutOutput struct {
	PublicKey   []byte
	AmountNanos uint64
}

func NewBitCloutInput() *BitCloutInput {
	return &BitCloutInput{
		TxID: BlockHash{},
	}
}

type UtxoKey struct {
	TxID  BlockHash
	Index uint32
}

type BitCloutInput UtxoKey

type MsgBitCloutTxn struct {
	TxInputs    []*BitCloutInput
	TxOutputs   []*BitCloutOutput
	TxnMeta     string //BitCloutTxnMetadata
	PublicKey   []byte
	ExtraData   map[string][]byte
	Signature   string
	TxnTypeJSON uint64
}

type MsgBitCloutTransactionBundle struct {
	Transactions []*MsgBitCloutTxn
}

func _readTransaction(id string, rr io.Reader) {
	//m := MsgBitCloutTxn{}
	numInputs, _ := ReadUvarint(rr)
	for ii := uint64(0); ii < numInputs; ii++ {
		currentInput := NewBitCloutInput()
		io.ReadFull(rr, currentInput.TxID[:])
		ReadUvarint(rr)
	}
	numOutputs, _ := ReadUvarint(rr)
	for ii := uint64(0); ii < numOutputs; ii++ {
		currentOutput := &BitCloutOutput{}
		currentOutput.PublicKey = make([]byte, btcec.PubKeyBytesLenCompressed)
		io.ReadFull(rr, currentOutput.PublicKey)
		ReadUvarint(rr)
	}
	txnMetaType, _ := ReadUvarint(rr)
	metaLen, _ := ReadUvarint(rr)
	metaBuf := make([]byte, metaLen)
	io.ReadFull(rr, metaBuf)
	if txnMetaType == 5 { // TxnTypeSubmitPost
		fmt.Println("txnMetaType", txnMetaType, metaLen)
		meta := SubmitPostMetadataFromBytes(metaBuf)
		ts := int64(meta.TimestampNanos / 1000000000)
		fmt.Println(id, "Timestamp", time.Unix(ts, 0))
		fmt.Println(id, "body", string(meta.Body))
	}
	pkLen, _ := ReadUvarint(rr)
	PublicKey := make([]byte, pkLen)
	io.ReadFull(rr, PublicKey)
	extraDataLen, _ := ReadUvarint(rr)
	if extraDataLen != 0 {
		ExtraData := make(map[string][]byte, extraDataLen)
		for ii := uint64(0); ii < extraDataLen; ii++ {
			var keyLen uint64
			keyLen, _ = ReadUvarint(rr)
			keyBytes := make([]byte, keyLen)
			io.ReadFull(rr, keyBytes)
			key := string(keyBytes)
			var valueLen uint64
			valueLen, _ = ReadUvarint(rr)
			value := make([]byte, valueLen)
			io.ReadFull(rr, value)
			ExtraData[key] = value
		}
	}
	sigLen, _ := ReadUvarint(rr)
	if sigLen != 0 {
		sigBytes := make([]byte, sigLen)
		io.ReadFull(rr, sigBytes)
	}
}

func MsgBitCloutTransactionBundleFromBytes(id string, data []byte) *MsgBitCloutTransactionBundle {
	rr := bytes.NewReader(data)
	m := MsgBitCloutTransactionBundle{}

	numTransactions, _ := ReadUvarint(rr)

	fmt.Println("numTransactions", numTransactions)
	for ii := uint64(0); ii < numTransactions; ii++ {
		_readTransaction(id, rr)
		//m.Transactions = append(m.Transactions, retTransaction)
	}

	return &m
}

type SubmitPostMetadata struct {
	PostHashToModify         []byte
	ParentStakeID            []byte
	Body                     []byte
	CreatorBasisPoints       uint64
	StakeMultipleBasisPoints uint64
	TimestampNanos           uint64
	IsHidden                 bool
}

func ReadVarString(rr io.Reader) []byte {
	StringLen, _ := ReadUvarint(rr)
	ret := make([]byte, StringLen)
	io.ReadFull(rr, ret)
	return ret
}

func SubmitPostMetadataFromBytes(data []byte) *SubmitPostMetadata {
	m := SubmitPostMetadata{}
	rr := bytes.NewReader(data)

	m.PostHashToModify = ReadVarString(rr)
	m.ParentStakeID = ReadVarString(rr)
	m.Body = ReadVarString(rr)
	m.CreatorBasisPoints, _ = ReadUvarint(rr)
	m.StakeMultipleBasisPoints, _ = ReadUvarint(rr)
	m.TimestampNanos, _ = ReadUvarint(rr)
	//m.IsHidden = _readBoolByte(rr)

	return &m
}
