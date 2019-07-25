package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/gob"
	"bytes"
	"strconv"
	"time"
	"strings"
	"math"
	"log"
)

type TxData struct {
	Timestamp	int
	From		string
	To			string
	Amount		float64
	Fee			float64
	Tax			float64
	Hash		string
	Sign		string
	Nonce		int
	Message		string
	Etc			string
	Datatype	string
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
	toStr:=td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Tax,'f',-1,64)+TxDelim+td.Sign+TxDelim+strconv.Itoa(td.Nonce)+TxDelim+td.Message+TxDelim+td.Datatype
	return setHash(toStr)
}

func (td *TxData) String() string {
	toStr:=strconv.Itoa(td.Timestamp)+TxDelim+td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Tax,'f',-1,64)+TxDelim+td.Hash+TxDelim+td.Sign+TxDelim+strconv.Itoa(td.Nonce)+TxDelim+td.Message+TxDelim+td.Etc+TxDelim+td.Datatype
	return toStr
}

func (td *TxData) TxString() string {
	toStr:=strconv.Itoa(td.Timestamp)+BlockDelim+td.From+BlockDelim+td.To+BlockDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+BlockDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+BlockDelim+strconv.FormatFloat(td.Tax,'f',-1,64)+BlockDelim+td.Hash+BlockDelim+td.Sign+BlockDelim+strconv.Itoa(td.Nonce)+BlockDelim+td.Message+BlockDelim+td.Etc+BlockDelim+td.Datatype
	return toStr
}

func (td *TxData) BtyeString() string {
	toStr:=strconv.Itoa(td.Timestamp)+TxDelim+td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Tax,'f',-1,64)+TxDelim+td.Hash+TxDelim+td.Sign+TxDelim+strconv.Itoa(td.Nonce)+TxDelim+td.Message
	return toStr
}

func (td *TxData) SetHash() {
	td.Hash=td.GetHash()
}

func (td *TxData) GetHash() string {
	return td.HashString()
}

func (td *TxData) GetHash00() string {
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
	echo ("Tax =",td.Tax)
	echo ("Hash =",td.Hash)
	echo ("Sign =",td.Sign)
	echo ("Nonce =",td.Nonce)
	echo ("Message =",td.Message)
	echo ("Etc =",td.Etc)
}

func (td *TxData) Transaction() string {
	result:=NodeSend("tx",td.String())
	return result
}

func GenesisTx(blockno int) string {
	result:=GenBlock[blockno-1]
	return result
}

func GenCubeTx(cubeno int,blockno int) string {
	var cgBlock [27]string
	path:="./config/genfile"+strconv.Itoa(cubeno)
	GenFile=FileReadString(path)
	line:=strings.Split(GenFile,"\r\n")
	for _,v:=range line {
		result:=strings.Split(v, "|")
		genno,ok:=strconv.Atoi(result[0])
		if ok==nil {
			cgBlock[genno-1]+=v+"\r\n"
		}
	}
	result:=cgBlock[blockno-1]
	return result
}

func SpecialTx(blockno int) string {
	result:=""
	return result
}

func TxPool(cubeno int,blockno int) string {
	if cubeno==1 {
		return GenesisTx(blockno)
	} else if cubeno==2 {
		return GenCubeTx(cubeno,blockno)
	} else	if blockno==Configure.Indexing+1 || blockno==Configure.Statistics+1 || blockno==Configure.Format+1 || blockno==Configure.Edit+1 {
		return SpecialTx(blockno)
	} else {
		result:=NodeSend("txpool","0&cubeno="+strconv.Itoa(cubeno)+"&blockno="+strconv.Itoa(blockno))
		decho("cubeno="+strconv.Itoa(cubeno)+"&blockno="+strconv.Itoa(blockno)+"\n"+result)
		return result
	}
}

func TxpoolToBst(cubeno int,blockno int,TxStr string) (TxBST,MineResult) {
	if cubeno==1 || cubeno==2 {
		return TxpoolToGenesis(TxStr)
	} else	if blockno==Configure.Indexing+1 || blockno==Configure.Statistics+1 || blockno==Configure.Format+1 || blockno==Configure.Edit+1 {
		return TxpoolToSpecial(cubeno,blockno)
	} else {
		return TxpoolToTr(TxStr,blockno)
	}
}

func TxpoolToGenesis(TxStr string) (TxBST,MineResult) {
	var txData TxData
	var tbst TxBST
	var mresult MineResult

	tbst.Init()

	if strings.Index(TxStr,"|")<0 {
		return tbst,mresult
	}
	line:=strings.Split(TxStr, "\r\n")
	for _,v:=range line {
		if len(v)<20 {
			continue
		}
		decho (v)
		result:=strings.Split(v, "|")
		mresult.TxMine=result[0]+","
		txData.Timestamp=int(time.Now().Unix())
		txData.From=result[1]
		txData.To=result[2]
		txData.Amount,_=strconv.ParseFloat(result[3],64)
		txData.Fee,_=strconv.ParseFloat(result[4],64)
		txData.Tax,_=strconv.ParseFloat(result[5],64)
		txData.Datatype=result[6]
		txData.Nonce,_=strconv.Atoi(result[7])
		txData.Message=result[8]
		txData.Sign=setHash(setHash(strconv.Itoa(txData.Timestamp)+txData.From+txData.To+txData.Datatype+strconv.Itoa(txData.Nonce)))
		txData.SetHash()

		mresult.Sumfee+=txData.Fee
		mresult.Sumfee=math.Round(mresult.Sumfee*100000000)/100000000
		if txData.Datatype=="QUB" || txData.Datatype=="POW" || txData.Datatype=="POS"  {
			mresult.Txcnt++
			mresult.Txamount+=txData.Amount+txData.Tax
			tbst.treeInsert(txData,"Coin")
		} else if txData.Datatype=="Contract" {
			mresult.Concnt++
			tbst.treeInsert(txData,"Contract")
		} else if len(txData.Datatype)>=2 && len(txData.Datatype)<=5 {
			mresult.Tkcnt++
			tbst.treeInsert(txData,"Token")
		} else {
			tbst.treeInsert(txData,"Data")
		}
	}
	return tbst,mresult
}

func TxpoolToSpecial(cubeno int,blockno int) (TxBST,MineResult) {
	var tbst TxBST
	var mresult MineResult
	tbst.Init()
	if blockno==Configure.Indexing+1 {
		ci:=CubeIndexing(cubeno-1)
		tbst.treeInsert(ci,"Data")
	} else if blockno==Configure.Statistics+1 {
		ci:=CubeStatistic(cubeno-1)
		tbst.treeInsert(ci,"Data")
	}
	return tbst,mresult
}

func TxpoolToTr(TxStr string,blockno int) (TxBST,MineResult) {
	var txData TxData
	var tbst TxBST
	var mresult MineResult
	tbst.Init()

	if strings.Index(TxStr,"|")<0 {
		return tbst,mresult
	}

	line:=strings.Split(TxStr,"\n")
	for _,v:=range line {
		if len(v)<20 {
			continue
		}
		result:=strings.Split(v, "|")
		
		mresult.TxMine+=result[0]+","
		txc,_:=strconv.Atoi(result[11])
		
		txData.Timestamp,_=strconv.Atoi(result[8])
		txData.From=result[1]
		txData.To=result[2]
		txData.Fee,_=strconv.ParseFloat(result[5],64)
		txData.Tax,_=strconv.ParseFloat(result[9],64)
		txData.Sign=result[7]
		txData.Nonce,_=strconv.Atoi(result[8])
		txData.Datatype=result[3]
		txData.Message=result[10]
		txData.Etc=result[12]
		txData.Amount,_=strconv.ParseFloat(result[4],64)

		if txc>0 && txData.Datatype!="Contract" && txData.Datatype!="Data" {
			txData.Nonce+=blockno*10000000;
			for i:=0;i<txc;i++ {
				txData.Nonce++
				txData.SetHash()
				mresult.Sumfee+=txData.Fee
				mresult.Sumfee=math.Round(mresult.Sumfee*100000000)/100000000
				if txData.Datatype=="QUB" || txData.Datatype=="POW" || txData.Datatype=="POS"  {
					mresult.Txcnt++
					mresult.Txamount+=txData.Amount+txData.Tax
					tbst.treeInsert(txData,"Coin")
				} else if len(txData.Datatype)>=2 && len(txData.Datatype)<=5 {
					mresult.Tkcnt++
					tbst.treeInsert(txData,"Token")
				} else {
				}
			}
		} else {
			txData.SetHash()
			mresult.Sumfee+=txData.Fee
			mresult.Sumfee=math.Round(mresult.Sumfee*100000000)/100000000
			if txData.Datatype=="QUB" || txData.Datatype=="POW" || txData.Datatype=="POS"  {
				mresult.Txcnt++
				mresult.Txamount+=txData.Amount+txData.Tax
				tbst.treeInsert(txData,"Coin")
			} else if txData.Datatype=="Contract" {
				mresult.Concnt++
				tbst.treeInsert(txData,"Contract")
			} else if txData.Datatype=="Data" {
				tbst.treeInsert(txData,"Data")
			} else if len(txData.Datatype)>=2 && len(txData.Datatype)<=5 {
				mresult.Tkcnt++
				tbst.treeInsert(txData,"Token")
			} else {
			}
		}
	}
	return tbst,mresult
}

func TxpoolToBst00(TxStr string) (TxBST,MineResult) {
	var txData TxData
	var tbst TxBST
	var mresult MineResult

	tbst.Init()

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
		txData.Tax,_=strconv.ParseFloat(result[9],64)
		txData.Hash=result[6]
		txData.Sign=result[7]
		txData.Nonce,_=strconv.Atoi(result[8])
		txData.Datatype=result[3]
		txData.Message=result[10]

		mresult.Sumfee+=txData.Fee
		if txData.Datatype=="QUB" || txData.Datatype=="POW" || txData.Datatype=="POS"  {
			mresult.Txcnt++
			mresult.Txamount+=txData.Amount+txData.Tax
			tbst.treeInsert(txData,"Coin")
		} else if txData.Datatype=="Contract" {
			mresult.Concnt++
			tbst.treeInsert(txData,"Contract")
		} else if len(txData.Datatype)>=2 && len(txData.Datatype)<=5 {
			mresult.Tkcnt++
			tbst.treeInsert(txData,"Token")
		} else {
			tbst.treeInsert(txData,"Data")
		}
	}
	return tbst,mresult
}

func TxBlock(cubeno int,blockno int) (TxBST,MineResult) {
	var btx TxBST
	var mresult MineResult 
	r:=TxPool(cubeno,blockno)
	btx,mresult=TxpoolToBst(cubeno,blockno,r)
	decho(r)
	return btx,mresult
}

func TxBlockData(cubeno int,blockno int) ([]byte,MineResult) {
	btx,mresult:=TxBlock(cubeno,blockno)
	txp:=GetBytes(btx)
	return txp,mresult
}


