package lib

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
)

var MiniPostChan = make(chan MiniPost, 1024)

type BlockHash [32]byte

func HashToBigint(hash *BlockHash) *big.Int {
	val, _ := new(big.Int).SetString(hex.EncodeToString(hash[:]), 16)
	return val
}

func (bh *BlockHash) String() string {
	return fmt.Sprintf("%064x", HashToBigint(bh))
}

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

func _checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

type MiniPost struct {
	Timestamp     int64
	Body          string
	PosterPub58   string
	PostHashHex   string
	PostExtraData map[string][]byte
	ParentStakeID []byte
}

func _readTransaction(id string, rr io.Reader) {
	mp := MiniPost{}
	history := []byte{}
	numInputs, h := ReadUvarint(rr) // *
	history = append(history, h...)
	for ii := uint64(0); ii < numInputs; ii++ {
		currentInput := NewBitCloutInput()
		io.ReadFull(rr, currentInput.TxID[:]) // *
		history = append(history, currentInput.TxID[:]...)
		_, h := ReadUvarint(rr) // *
		history = append(history, h...)
	}
	numOutputs, h := ReadUvarint(rr) // *
	history = append(history, h...)
	for ii := uint64(0); ii < numOutputs; ii++ {
		currentOutput := &BitCloutOutput{}
		currentOutput.PublicKey = make([]byte, btcec.PubKeyBytesLenCompressed)
		io.ReadFull(rr, currentOutput.PublicKey) // *
		history = append(history, currentOutput.PublicKey...)
		_, h := ReadUvarint(rr) // *
		history = append(history, h...)
	}
	txnMetaType, h := ReadUvarint(rr) // *
	history = append(history, h...)
	metaLen, h := ReadUvarint(rr) // *
	history = append(history, h...)

	metaBuf := make([]byte, metaLen)
	io.ReadFull(rr, metaBuf) // *
	history = append(history, metaBuf...)
	if txnMetaType == 5 { // TxnTypeSubmitPost
		//fmt.Println("txnMetaType", txnMetaType, metaLen)
		meta := SubmitPostMetadataFromBytes(metaBuf)
		mp.Timestamp = int64(meta.TimestampNanos / 1000000000)

		//fmt.Println(id, "Timestamp", time.Unix(mp.Timestamp, 0))
		mp.Body = string(meta.Body)
		//fmt.Println(id, "body", mp.Body)
		mp.ParentStakeID = meta.ParentStakeID
	}
	pkLen, h := ReadUvarint(rr) // *
	history = append(history, h...)
	PublicKey := make([]byte, pkLen)
	io.ReadFull(rr, PublicKey) // *
	history = append(history, PublicKey...)
	if txnMetaType == 5 { // TxnTypeSubmitPost
		PublicKey = append([]byte{205, 20, 0}, PublicKey...)
		suffix := _checksum(PublicKey)
		PublicKey = append(PublicKey, suffix[:]...)
		pub58 := base58.Encode(PublicKey)
		//fmt.Println(id, "PublicKey", pub58, len(PublicKey))
		mp.PosterPub58 = pub58
	}
	extraDataLen, h := ReadUvarint(rr) // *
	history = append(history, h...)
	if extraDataLen != 0 {
		mp.PostExtraData = make(map[string][]byte, extraDataLen)
		for ii := uint64(0); ii < extraDataLen; ii++ {
			var keyLen uint64
			keyLen, h = ReadUvarint(rr) // *
			history = append(history, h...)
			keyBytes := make([]byte, keyLen)
			io.ReadFull(rr, keyBytes) // *
			history = append(history, keyBytes...)
			key := string(keyBytes)
			var valueLen uint64
			valueLen, h = ReadUvarint(rr) // *
			history = append(history, h...)
			value := make([]byte, valueLen)
			io.ReadFull(rr, value) // *
			history = append(history, value...)
			mp.PostExtraData[key] = value
		}
	}
	sigLen, h := ReadUvarint(rr) // *
	history = append(history, h...)
	if sigLen != 0 {
		sigBytes := make([]byte, sigLen)
		io.ReadFull(rr, sigBytes) // *
		history = append(history, sigBytes...)
	}
	if txnMetaType == 5 {
		mp.PostHashHex = Sha256DoubleHash(history).String()
		//fmt.Println("postHash", mp.PostHashHex)
		MiniPostChan <- mp
	}
}

func MsgBitCloutTransactionBundleFromBytes(id string, data []byte) *MsgBitCloutTransactionBundle {
	rr := bytes.NewReader(data)
	m := MsgBitCloutTransactionBundle{}

	numTransactions, _ := ReadUvarint(rr)

	//fmt.Println("numTransactions", numTransactions)
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
