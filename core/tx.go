package core

import (
	"time"
	"log"
	"bytes"
	"encoding/gob"
	"strconv"
	"../wallet"
	"../lib"
)

func (tx TxData) TxInput(w wallet.Wallet,from string, to string,amount int) TxData {
	tx.DataType="tx"
	txd:=new(TransactionData)
	txd.Timestamp=int(time.Now().Unix())
	txd.From=lib.StrToByte(from)
	txd.To=lib.StrToByte(to)
	txd.Amount=amount
	txd.Nonce=100
	txd.Hash=lib.StrToByte(setHash(strconv.Itoa(txd.Timestamp)+from+to+strconv.Itoa(amount)+strconv.Itoa(txd.Nonce)))
	txd.Sign,_=w.Sign(txd.Hash)
	tx.DataTx=*txd
	return tx
}

func (tx *TxData) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func TxpoolStr(tx TxPool) string {
	var str string
	for _,v := range tx.Tdata {
		str+=TxpdataStr(v.DataTx)
	}
	return str
}

func TxpdataStr(tx TransactionData) string {
	str:=strconv.Itoa(tx.Timestamp)+ lib.ByteToStr(tx.From) + lib.ByteToStr(tx.To) +strconv.Itoa(tx.Amount) + lib.ByteToStr(tx.Hash) + lib.ByteToStr(tx.Sign) +strconv.Itoa(tx.Nonce)
	return str
}

func Serialize(object interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(object)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func Deserialize(data []byte) TxPool {
	var transaction TxPool
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

func TxDeserialize(data []byte) TransactionData {
	var transaction TransactionData
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

func IndexDeserialize(data []byte) IBlock {
	var idata IBlock
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&idata)
	if err != nil {
		log.Panic(err)
	}
	return idata
}

func StaticDeserialize(data []byte) SBlock {
	var idata SBlock
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&idata)
	if err != nil {
		log.Panic(err)
	}
	return idata
}

func EscrowDeserialize(data []byte) EBlock {
	var idata EBlock
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&idata)
	if err != nil {
		log.Panic(err)
	}
	return idata
}

func EscrowInput(w wallet.Wallet,from string, to string,amount int,Etype int,EKey string,Etime int) EscrowData {
	var tx EscrowData
	txd:=new(TransactionData)
	txd.Timestamp=int(time.Now().Unix())
	txd.From=lib.StrToByte(from)
	txd.To=lib.StrToByte(to)
	txd.Amount=amount
	txd.Nonce=100
	txd.Hash=lib.StrToByte(setHash(strconv.Itoa(txd.Timestamp)+from+to+strconv.Itoa(amount)+strconv.Itoa(txd.Nonce)))
	txd.Sign,_=w.Sign(txd.Hash)
	
	tx.EscrowTx=*txd
	tx.EscrowType=Etype
	tx.EscrowKey=EKey
	tx.EscrowTime=Etime
	tx.State=0
	return tx
}