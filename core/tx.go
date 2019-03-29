package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/gob"
	"bytes"
	"strconv"
	"time"
	"strings"
	"log"
)

type TxData struct {
	Timestamp	int
	From		string
	To			string
	Amount		float64
	Fee			float64
	Hash		string
	Sign		string
	Nonce		int
	Message		string
	Datatype	string
}

func (td *TxData) Set() {
	// Input
	// Verify
}

func (td *TxData) Input(TxStr string) {
	if len(TxStr)<20 {
		return
	}
	result:=strings.Split(TxStr, "|")
	td.Timestamp=int(time.Now().Unix())
	td.From=result[1]
	td.To=result[2]
	td.Amount,_=strconv.ParseFloat(result[3],64)
	td.Fee,_=strconv.ParseFloat(result[4],64)
	td.Sign=result[5]
	td.Nonce=GetTxCount(td.From)
	td.SetHash()
}

func (td *TxData) HashString() string {
	toStr:=strconv.Itoa(td.Timestamp)+TxDelim+td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+td.Sign+TxDelim+strconv.Itoa(td.Nonce)+TxDelim+td.Message
	return setHash(toStr)
}

func (td *TxData) String() string {
	toStr:=strconv.Itoa(td.Timestamp)+TxDelim+td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+td.Hash+TxDelim+td.Sign+TxDelim+strconv.Itoa(td.Nonce)+TxDelim+td.Message
	return toStr
}

func (td *TxData) BtyeString() string {
	toStr:=strconv.Itoa(td.Timestamp)+TxDelim+td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+td.Hash+TxDelim+td.Sign+TxDelim+strconv.Itoa(td.Nonce)+TxDelim+td.Message
	return toStr
}

func (td *TxData) SetHash() {
	td.Hash=td.GetHash()
}

func (td *TxData) GetHash() string {
	hashstr:=td.HashString()
	h:=sha256.New()
	h.Write([]byte(hashstr))
	return hex.EncodeToString(h.Sum(nil))
}

func (td *TxData) Verify() bool {
	return td.Hash==td.GetHash()
}

func (td *TxData) ToByte() []byte {
	var buff bytes.Buffer
	enc:=gob.NewEncoder(&buff)
	err:=enc.Encode(td)
	if err!=nil {
		echo (err)
	}
	return buff.Bytes()
}

func ByteToTx(data []byte) TxData {
	var Tdata TxData
	gob.Register(Tdata)
	decoder:=gob.NewDecoder(bytes.NewReader(data))
	err:=decoder.Decode(&Tdata)
	if err != nil {
		log.Panic(err)
	}
	return Tdata
}

func (td *TxData) GetCount() int {
	return GetTxCount(td.From)
}

func (td *TxData) Print() {
	echo ("Timestamp =",td.Timestamp)
	echo ("From =",td.From)
	echo ("To =",td.To)
	echo ("Amount =",td.Amount)
	echo ("Fee =",td.Fee)
	echo ("Hash =",td.Hash)
	echo ("Sign =",td.Sign)
	echo ("Nonce =",td.Nonce)
	echo ("Message =",td.Message)
}

func (td *TxData) Transaction() string {
	result:=NodeSend("tx",td.String())
	return result
}

func TxPool(cubeno int,blockno int) string {
	result:=NodeSend("txpool","cubeno="+strconv.Itoa(cubeno)+"&blockno="+strconv.Itoa(blockno))
	return result
}

func TxpoolToBst(TxStr string) (TxBST,MineResult) {
	var txData TxData
	var tbst TxBST
	var mresult MineResult

	pt:=PowTx()
	ph:=tbst.treeInsertNode(pt,"Poh")

	if strings.Index(TxStr,"|")<0 {
		return tbst,mresult
	}
	line:=strings.Split(TxStr, "\n")
	for _,v:=range line {
		if len(v)<20 {
			continue
		}
 		//echo (v)
		result:=strings.Split(v, "|")
		mresult.TxMine+=result[0]+","
		txData.Timestamp,_=strconv.Atoi(result[8])
		txData.From=result[1]
		txData.To=result[2]
		txData.Amount,_=strconv.ParseFloat(result[4],64)
		txData.Fee,_=strconv.ParseFloat(result[5],64)
		txData.Hash=result[6]
		txData.Sign=result[7]
		txData.Nonce,_=strconv.Atoi(result[8])
		txData.Datatype=result[3];

		mresult.Sumfee+=txData.Fee
		if txData.Datatype=="QUB" || txData.Datatype=="POW" || txData.Datatype=="POS"  {
			mresult.Txcnt++
			mresult.Txamount+=txData.Amount
			tbst.treeInsert(txData,"Coin")
		} else if txData.Datatype=="Contract" {
			tbst.treeInsert(txData,"Contract")
		} else if len(txData.Datatype)>=2 && len(txData.Datatype)<=5 {
			mresult.Tkcnt++
			tbst.treeInsert(txData,"Token")
		} else {
			tbst.treeInsert(txData,"Data")
		}
	}
	
	txData=pt
	if mresult.Sumfee>10 { mresult.Sumfee=10 }
	txData.Fee=mresult.Sumfee
	txData.Amount=txData.Amount+mresult.Sumfee
	ph.Val=txData
	Sumfee+=mresult.Sumfee

	return tbst,mresult
}

func TxBlock(cubeno int,blockno int) (TxBST,MineResult) {
	r:=TxPool(cubeno,blockno)
	btx,mresult:=TxpoolToBst(r)
	decho(r)
	return btx,mresult
}


func TxBlockData(cubeno int,blockno int) ([]byte,MineResult) {
	btx,mresult:=TxBlock(cubeno,blockno)
	txp:=GetBytes(btx)
	return txp,mresult
}
