package mock

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"google.golang.org/protobuf/types/known/anypb"

	"cosmossdk.io/core/transaction"
)

var _ transaction.Tx = Tx{}

type Tx struct {
	Sender   []byte
	Msg      transaction.Type
	GasLimit uint64
}

func (t Tx) Hash() [32]byte {
	return sha256.Sum256(t.Bytes())
}

func (t Tx) GetMessages() ([]transaction.Type, error) {
	if t.Msg == nil {
		return nil, errors.New("messages not available or are nil")
	}
	return []transaction.Type{t.Msg}, nil
}

func (t Tx) GetSenders() ([]transaction.Identity, error) {
	if t.Sender == nil {
		return nil, errors.New("senders not available or are nil")
	}
	return []transaction.Identity{t.Sender}, nil
}

func (t Tx) GetGasLimit() (uint64, error) {
	return t.GasLimit, nil
}

type encodedTx struct {
	Sender   []byte     `json:"sender"`
	Msg      *anypb.Any `json:"message"`
	GasLimit uint64     `json:"gas_limit"`
}

func (t Tx) Bytes() []byte {
	msg, err := anypb.New(t.Msg)
	if err != nil {
		panic(err)
	}
	tx, err := json.Marshal(encodedTx{
		Sender:   t.Sender,
		Msg:      msg,
		GasLimit: t.GasLimit,
	})
	if err != nil {
		panic(err)
	}
	return tx
}

func (t *Tx) Decode(b []byte) {
	rawTx := new(encodedTx)
	err := json.Unmarshal(b, rawTx)
	if err != nil {
		panic(err)
	}
	msg, err := rawTx.Msg.UnmarshalNew()
	if err != nil {
		panic(err)
	}
	t.Msg = msg
	t.Sender = rawTx.Sender
	t.GasLimit = rawTx.GasLimit
}

type txCodec struct{}

func (txCodec) Decode(bytes []byte) (Tx, error) {
	t := new(Tx)
	t.Decode(bytes)
	return *t, nil
}
