package lib

type MsgBitCloutVersion struct {
	Version              uint64
	Services             uint64
	TstampSecs           int64
	Nonce                uint64
	UserAgent            string
	StartBlockHeight     uint32
	MinFeeRateNanosPerKB uint64
}
