package core

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"strings"
	"log"
	"time"
	"os"
	"math"
)

type IBlock struct {
	IndexAddress	[]IndexType
}

type IndexType struct {
	Address		string
	Indexing	[]CubeIndex
}

type TxIndexing struct {
	Cubeno		int	
	AddrIndex	map[string]string
}

type TxStatistic struct {
	Cubeno		int	
	AddrIndex	map[string]string
}


type CubeIndex struct {
	Index	int
	CubeNum	int
}

type StaticRule1 struct {
	Address string
	Balance	float64
}

type SBlock struct {
	RuleArr	[]StaticRule1
}

type EscrowData struct {
	EscrowTx	TxData
	EscrowType  int
	EscrowKey	string
	EscrowTime	int 
	State		int
}

type OpenRule struct {
	HashId	    string
	Key2nd		string
	State		int    
}

type EBlock struct {
	Escrow		[]EscrowData 
}



var AddrStatistics	map[string]StatisticData

type TxStatistics struct {
	Cubeno		int	
	AddrIndex	map[string]StatisticData
}

type StatisticData struct {
	Cubeno		int	
	Addr		string
	Balance		float64
	RBalance	float64
	SBalance	float64
	TBalance	float64
	Txcnt		int
	Tkcnt		int
	Pos			string
}

func CubeIndexing(cubeno int) map[string]string {
	var iData []TxData
	var cubeIndexing map[string]string
	cubeIndexing=make(map[string]string)
	if cubeno<=0 {
		return cubeIndexing
	}
	var c Cube 
	c.Cubeno=cubeno
	c.Read()
	for i:=-1;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			if i==-1 {
				tx1,tx2:=GetCubePoh(cubeno)
				iData=append(iData,tx1)
				iData=append(iData,tx2)
			} else {
				iData=BlockTxData(c.Blocks[i].Data)
			}
			cIndexValue:=CubeIndexToStr(CubeIndex{cubeno,i+1})
			for _,v := range iData {
				if v.Datatype!="NULL" && v.Datatype!="Data" && v.Datatype!="Contract"   {
					if(len(v.From)==34) {
						if v.From[0:31]=="C"+strings.Repeat("0",30) {
						} else if cubeIndexing[v.From]=="" {
							cubeIndexing[v.From]=cIndexValue+","
						} else if strings.Contains(cubeIndexing[v.From],cIndexValue+",")==false && cIndexValue>"1" {
							cubeIndexing[v.From]+=cIndexValue+","
						}
					}
					if(len(v.To)==34) {
						if v.To[0:31]=="C"+strings.Repeat("0",30) {
						} else if cubeIndexing[v.To]=="" {
							cubeIndexing[v.To]=cIndexValue+","
						} else if strings.Contains(cubeIndexing[v.To],cIndexValue+",")==false && cIndexValue>"1" {
							cubeIndexing[v.To]+=cIndexValue+","
						}
					}
				}
			}
		}
	}
	return cubeIndexing
}

func CubeStatistic(cubeno int) map[string]string {
	var iData []TxData
	var statAddr map[string]bool
	statAddr=make(map[string]bool)

	var cubeStatistic map[string]string
	cubeStatistic=make(map[string]string)

	if cubeno<=0 {
		return cubeStatistic
	}

	var c Cube 
	c.Cubeno=cubeno
	c.Read()
	for i:=-1;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			if i==-1 {
				tx1,tx2:=GetCubePoh(cubeno)
				iData=append(iData,tx1)
				iData=append(iData,tx2)
			} else {
				iData=BlockTxData(c.Blocks[i].Data)
			}
			for _,v := range iData {
				if v.Datatype!="Contract" && v.Datatype!="Data" {
					if(len(v.From)==34) {
						if v.From[0:31]!="C"+strings.Repeat("0",30) {
							statAddr[v.From]=true
							
						}
					}
					if(len(v.To)==34) {
						if v.To[0:31]!="C"+strings.Repeat("0",30) {
							statAddr[v.To]=true
						}
					}
				}
			}
		}
	}
	for k,_ := range statAddr {
		cubeStatistic[k]=StatisticToStr(k,cubeno)
	}
	return cubeStatistic
}

func CubeOnlyStatistic(cubeno int) map[string]string {
	var iData []TxData
	var statAddr map[string]bool
	statAddr=make(map[string]bool)

	var cubeStatistic map[string]string
	cubeStatistic=make(map[string]string)

	if cubeno<=0 {
		return cubeStatistic
	}

	var c Cube 
	c.Cubeno=cubeno
	c.Read()
	for i:=-1;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			if i==-1 {
				tx1,tx2:=GetCubePoh(cubeno)
				iData=append(iData,tx1)
				iData=append(iData,tx2)
			} else {
				iData=BlockTxData(c.Blocks[i].Data)
			}
			for _,v := range iData {
				if v.Datatype!="Contract" && v.Datatype!="Data" {
					if(len(v.From)==34) {
						if v.From[0:31]!="C"+strings.Repeat("0",30) {
							statAddr[v.From]=true
							
						}
					}
					if(len(v.To)==34) {
						if v.To[0:31]!="C"+strings.Repeat("0",30) {
							statAddr[v.To]=true
						}
					}
				}
			}
		}
	}
	for k,_ := range statAddr {
		cubeStatistic[k]=StatisticToStrCube(k,cubeno)
	}
	return cubeStatistic
}

func CubeStatistic01(cubeno int) map[string]string {
	var iData []TxData
	var statAddr map[string]bool
	statAddr=make(map[string]bool)

	var cubeStatistic map[string]string
	cubeStatistic=make(map[string]string)

	if cubeno<=0 {
		return cubeStatistic
	}

	var c Cube 
	c.Cubeno=cubeno
	c.Read()
	for i:=-1;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			if i==-1 {
				tx1,tx2:=GetCubePoh(cubeno)
				iData=append(iData,tx1)
				iData=append(iData,tx2)
			} else {
				iData=BlockTxData(c.Blocks[i].Data)
			}
			for _,v := range iData {
				if v.Datatype!="Contract" && v.Datatype!="Data" {
					if v.From[0:31]!="C"+strings.Repeat("0",30) {
						statAddr[v.From]=true
						
					}
					if v.To[0:31]!="C"+strings.Repeat("0",30) {
						statAddr[v.To]=true
					}
				}
			}
		}
	}
	for k,_ := range statAddr {
		cStatisticValue:=CubeStatisticToStr(k)
		cubeStatistic[k]=cStatisticValue
	}
	return cubeStatistic
}

func CubeIndexing00(cubeno int) map[string]string {
	var iBlock [27]Block
	var cubeIndexing map[string]string
	cubeIndexing=make(map[string]string)
	c:=cubeno
	if c<=0 {
		return cubeIndexing
	}
	for i:=0;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			err:=BlockRead(c,i+1,&iBlock[i])
			Err(err,0)	
			iData:=BlockTxData(iBlock[i].Data)
			cIndexValue:=CubeIndexToStr(CubeIndex{c,i+1})
			for _,v := range iData {
				if v.Datatype!="Contract" && v.Datatype!="Data"  {
					if cubeIndexing[v.From]=="" {
						cubeIndexing[v.From]=cIndexValue+","
					} else if strings.Contains(cubeIndexing[v.From],cIndexValue+",")==false && cIndexValue>"1" {
						cubeIndexing[v.From]+=cIndexValue+","
					}
					if cubeIndexing[v.To]=="" {
						cubeIndexing[v.To]=cIndexValue+","
					} else if strings.Contains(cubeIndexing[v.To],cIndexValue+",")==false && cIndexValue>"1" {
						cubeIndexing[v.To]+=cIndexValue+","
					}
				}
			}
		}
	}
	return cubeIndexing
}

func CubeStatistic00(cubeno int) map[string]string {
	var iBlock [27]Block
	var statAddr map[string]bool
	statAddr=make(map[string]bool)

	var cubeStatistic map[string]string
	cubeStatistic=make(map[string]string)

	c:=cubeno
	if c<=0 {
		return cubeStatistic
	}

	for i:=0;i<27;i++ {
		err:=BlockRead(c,i+1,&iBlock[i])
		Err(err,0)	
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			iData:=BlockTxData(iBlock[i].Data)
			for _,v := range iData {
				if v.Datatype!="Contract" && v.Datatype!="Data" {
					statAddr[v.From]=true
					statAddr[v.To]=true
				}
			}
		}
	}
	for k,_ := range statAddr {
		cStatisticValue:=CubeStatisticToStr(k)
		cubeStatistic[k]=cStatisticValue
	}
	return cubeStatistic
}

func AllIndexing(cubeno int) map[string]string {
	var aIndexing TxIndexing
	aIndexing.AddrIndex=make(map[string]string)	

	if cubeno==0 {
		cubeno=CubeHeight()-1
	}
	if cubeno<=0 {
		cubeno=1
	}

	path:=Configure.Datafolder+filepathSeparator+"special"
	if DirExist(path)==false {
		if err:=os.MkdirAll(path, os.FileMode(0755)); err!=nil {
			echo ("Special block directory not found")
		}	
	}
	pathfile:=path+filepathSeparator+"Indexing.cbs"
	if DirExist(pathfile) {
		err:=FileRead(pathfile,&aIndexing)
		Err(err,0)	
	}
	if aIndexing.Cubeno>=cubeno {
	} else {
		for i:=aIndexing.Cubeno+1;i<cubeno+1;i++ {
			tmpIndexing:=CubeIndexing(i)
			for k,v := range tmpIndexing {
				if aIndexing.AddrIndex[k]=="" {
					aIndexing.AddrIndex[k]=v+","
				} else if v>"1" {
					aIndexing.AddrIndex[k]+=v+","
				}
			}
		}
		aIndexing.Cubeno=cubeno
		err:=FileWrite(pathfile,aIndexing)
		Err(err,0)	
	}
	return aIndexing.AddrIndex
}

func AllStatistic(cubeno int) map[string]string {
	var aStatistic TxStatistic
	aStatistic.AddrIndex=make(map[string]string)	

	if cubeno==0 {
		cubeno=CubeHeight()-1
	}
	if cubeno<=0 {
		cubeno=1
	}

	path:=Configure.Datafolder+filepathSeparator+"special"
	if DirExist(path)==false {
		if err:=os.MkdirAll(path, os.FileMode(0755)); err!=nil {
			echo ("Special block directory not found")
		}	
	}
	pathfile:=path+filepathSeparator+"Statistic.cbs"
	if DirExist(pathfile) {
		err:=FileRead(pathfile,&aStatistic)
		Err(err,0)	
	}
	if aStatistic.Cubeno>=cubeno {
	} else {
		tmpStatistic:=CubeStatistic(cubeno)
		for k,v := range tmpStatistic {
			aStatistic.AddrIndex[k]=v
		}
		aStatistic.Cubeno=cubeno
		err:=FileWrite(pathfile,aStatistic)
		Err(err,0)	
	}
	return aStatistic.AddrIndex
}

func CubeStatistic_pp00(cubeno int) [4][]string {
	var iBlock [27]Block
	var dStatistic [4][]string
	c:=cubeno
	if c<=0 {
		return dStatistic
	}
	for i:=0;i<27;i++ {
		err:=BlockRead(c,i+1,&iBlock[i])
		Err(err,0)	
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			iData:=BlockTxData(iBlock[i].Data)
			for _,v := range iData {
				if v.Datatype!="Contract" && v.Datatype!="Data"  {
					if GetBalance(v.From)>=5000 {
						dStatistic[0]=append(dStatistic[0],v.From)
					} else {
						dStatistic[1]=append(dStatistic[1],v.From)
					}
					if GetBalance(v.To)>=5000 {
						dStatistic[0]=append(dStatistic[0],v.To)
					} else {
						dStatistic[1]=append(dStatistic[1],v.To)
					}
				}
			}
		}
	}
	return dStatistic
}

func IndexingAddr(addr string) (int,string) {
	var aIndexing TxIndexing
	aIndexing.AddrIndex=make(map[string]string)	
	IndexingRead(&aIndexing)
	result:=strings.Replace(aIndexing.AddrIndex[addr],",,",",",-1)
	return aIndexing.Cubeno,result
}

func addIndexing() Block {
	c:=CurrentHeight()-1
	cnum:=Configure.Indexing;
	var iBlock [27]Block
	var iAdd []IndexType
	var pAdd IBlock
	var CArr []CubeIndex
	CI:=CubeIndex{c,cnum}
	CArr=append(CArr,CI)

	for i:=0;i<27;i++ {
		err:=BlockRead(c,i,iBlock[i])
		Err(err,0)	
		if i==Configure.Statistics || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
		} else if i==Configure.Indexing {
			pAdd=IndexDeserialize(iBlock[i].Data)
		} else {
			iData:=BlockTxData(iBlock[i].Data)
			for _,v := range iData {
				if v.Datatype=="QUB" {
					iAdd=append(iAdd,IndexType{v.From,CArr})
					iAdd=append(iAdd,IndexType{v.To,CArr})
				}
			}
		}
	}
	for k,v := range pAdd.IndexAddress {
		for _,i := range iAdd {
			if v.Address==i.Address {
				pAdd.IndexAddress[k].Indexing=append(pAdd.IndexAddress[k].Indexing,i.Indexing[0])
			} else {
				pAdd.IndexAddress=append(pAdd.IndexAddress,i)
			}
		}
	}
	return addBlock(Serialize(pAdd),c)
}

func CubeIndexToStr(ci CubeIndex) string {
	result:=strconv.Itoa(ci.Index)
	NumStr:=""
	if ci.CubeNum==0 {
		NumStr="B"
	} else if ci.CubeNum==27 {
		NumStr="A"
	} else if ci.CubeNum>=1 && ci.CubeNum<=26 {
		NumStr=string(ci.CubeNum+96)
	}
	result=result+NumStr
	return result
}

func StrToCubeIndex(str string) CubeIndex {
	idx,cno:=str[:len(str)-1],str[len(str)-1]
	index,_:=strconv.Atoi(idx)
	cbn:=0
	if string(cno)=="B" {
		cbn=0
	} else if string(cno)=="A" {
		cbn=27
	} else {
		cbn=int(cno)-96
	}
	result:=CubeIndex{index,cbn}
	return result
}

func AllStatisticToStr(addr string) string {
	result:=""
	pos:="F"
	balance:=GetBalance(addr)
	txcnt:=GetTxCount(addr)
	tkcnt:=GetTkCount(addr)
	if balance>=5000.0 {
		pos="T"
	}
	result=strconv.FormatFloat(balance,'f',-1,64)+","
	result+=strconv.Itoa(txcnt)+","
	result+=strconv.Itoa(tkcnt)+","
	result+=pos
	return result
}

func StatisticToStr(addr string,cubeno int) string {
	result:=""
	pos:="F"
	balance,rbalance,sbalance,tbalance,txcnt,tkcnt,gcubeno:=GetStaticVar(addr)
	if gcubeno<cubeno {
		for i:=gcubeno+1;i<=cubeno;i++ {
			rbalance1,sbalance1,txcnt1,tkcnt1:=GetStaticVarCube(addr,i)
			sbalance+=sbalance1
			sbalance=math.Round(sbalance*100000000)/100000000
			rbalance+=rbalance1
			rbalance=math.Round(rbalance*100000000)/100000000
			txcnt+=txcnt1
			tkcnt+=tkcnt1
		}
	}
	balance=rbalance-sbalance
	balance=math.Round(balance*100000000)/100000000
	tbalance=rbalance+sbalance
	tbalance=math.Round(tbalance*100000000)/100000000
	if balance>=5000.0 {
		pos="T"
	}
	result=strconv.FormatFloat(balance,'f',-1,64)+","
	result+=strconv.FormatFloat(rbalance,'f',-1,64)+","
	result+=strconv.FormatFloat(sbalance,'f',-1,64)+","
	result+=strconv.FormatFloat(tbalance,'f',-1,64)+","
	result+=strconv.Itoa(txcnt)+","
	result+=strconv.Itoa(tkcnt)+","
	result+=pos
	return result
}


func StatisticToStrCube(addr string,cubeno int) string {
	result:=""
	pos:="F"

	rbalance,sbalance,txcnt,tkcnt:=GetStaticVarCube(addr,cubeno)
	sbalance=math.Round(sbalance*100000000)/100000000
	rbalance=math.Round(rbalance*100000000)/100000000
	balance:=rbalance-sbalance
	balance=math.Round(balance*100000000)/100000000
	tbalance:=rbalance+sbalance
	tbalance=math.Round(tbalance*100000000)/100000000
	if balance>=5000.0 {
		pos="T"
	}
	result=strconv.FormatFloat(balance,'f',-1,64)+","
	result+=strconv.FormatFloat(rbalance,'f',-1,64)+","
	result+=strconv.FormatFloat(sbalance,'f',-1,64)+","
	result+=strconv.FormatFloat(tbalance,'f',-1,64)+","
	result+=strconv.Itoa(txcnt)+","
	result+=strconv.Itoa(tkcnt)+","
	result+=pos
	return result
}


func CubeStatisticToStr(addr string) string {
	result:=""
	pos:="F"

	rbalance,sbalance,txcnt:=GetStaticValue(addr)
	balance:=rbalance-sbalance
	balance=math.Round(balance*100000000)/100000000
	tbalance:=rbalance+sbalance
	tbalance=math.Round(tbalance*100000000)/100000000
	tkcnt:=GetTkCount(addr)
	if balance>=5000.0 {
		pos="T"
	}
	result=strconv.FormatFloat(balance,'f',-1,64)+","
	result+=strconv.FormatFloat(rbalance,'f',-1,64)+","
	result+=strconv.FormatFloat(sbalance,'f',-1,64)+","
	result+=strconv.FormatFloat(tbalance,'f',-1,64)+","
	result+=strconv.Itoa(txcnt)+","
	result+=strconv.Itoa(tkcnt)+","
	result+=pos
	return result
}

func addStatistic() Block {
	c:=CurrentHeight()-1
	cnum:=Configure.Statistics;
	var iBlock [27]Block
	var pAdd SBlock
	var sAddr []string
	var sBal []float64
	var iBal float64

	for i:=0;i<27;i++ {
		err:=BlockRead(c,i,iBlock[i])
		Err(err,0)	
		if i==Configure.Indexing || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
		} else if i==cnum {
			pAdd=StaticDeserialize(iBlock[i].Data)
		} else {
			iData:=BlockTxData(iBlock[i].Data)
			for _,v := range iData {
				if v.Datatype=="QUB" {
					sAddr=append(sAddr,v.From)
					sBal=append(sBal,v.Amount*(-1))
					sAddr=append(sAddr,v.To)
					sBal=append(sBal,v.Amount)
				}
			}
		}
	}
	for k,v := range pAdd.RuleArr {
		for j,i := range sAddr {
			if v.Address==i {
				pAdd.RuleArr[k].Balance+=sBal[j]
			} else {
				iBal=GetBalance(i)
				if iBal>=5000 {
					pAdd.RuleArr=append(pAdd.RuleArr,StaticRule1{i,iBal})
				}
			}
		}
	}
	return addBlock(Serialize(pAdd),c)
}

func addEscrow() Block {
	c:=CurrentHeight()-1
	cnum:=Configure.Escrow;
	var iBlock Block
	var pAdd EBlock

	err:=BlockRead(c,cnum,iBlock)
	Err(err,0)	
	pAdd=EscrowDeserialize(iBlock.Data)
	
	for k,v := range pAdd.Escrow {
		if v.State!=0 {
			pAdd.Escrow = append(pAdd.Escrow,pAdd.Escrow[k])
		}
	}
	return addBlock(Serialize(pAdd),c)
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

func Deserialize(data []byte) []TxData {
	var transaction []TxData
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

func TdDeserialize(data []byte) TxData {
	var Tdata TxData
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&Tdata)
	if err != nil {
		log.Panic(err)
	}
	return Tdata
}


func TxDeserialize(data []byte) TxData {
	var transaction TxData
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

func DataDeserialize(data []byte) Block {
	var idata Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&idata)
	if err != nil {
		log.Panic(err)
	}
	return idata
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


func addBlock(data []byte,Cubenum int) Block{
	var block Block
	var oblock Block
	ch:=CurrentHeight()
	if ch<1 {
		echo ("Genesis block not create")
		return block
	}
	err:=BlockRead(ch-1,Cubenum,&oblock)
	Err(err,0)
	block.Cubeno=ch
	block.Blockno=Cubenum
	block.Timestamp=int(time.Now().Unix())
	block.Nonce=0
	block.Data=data
	block.PrevHash=oblock.Hash
	block.PatternHash=block.GetPatternHash() 
	block.SetHash()
	return block
}


func MakeStat(filepath string) {
	var bStatistic TxStatistic

	aStatistic:=make(map[string]StatisticData)
	bStatistic.AddrIndex=make(map[string]string)	

	ch:=CubeHeight()
	
	path:=filepath+filepathSeparator+"special"

	if DirExist(path)==false {
		if err:=os.MkdirAll(path, os.FileMode(0755)); err!=nil {
			echo ("Special block directory not found")
		}	
	}
	
	for i:=1;i<ch;i++ {
		aStatistic=GetStaticData(i,aStatistic)
	}
	
	for k,v := range aStatistic {
		result:=""
		result=strconv.FormatFloat(v.Balance,'f',-1,64)+","
		result+=strconv.FormatFloat(v.RBalance,'f',-1,64)+","
		result+=strconv.FormatFloat(v.SBalance,'f',-1,64)+","
		result+=strconv.FormatFloat(v.TBalance,'f',-1,64)+","
		result+=strconv.Itoa(v.Txcnt)+","
		result+=strconv.Itoa(v.Tkcnt)+","
		result+=v.Pos
		bStatistic.AddrIndex[k]=result
	}
	
	bStatistic.Cubeno=ch-1
	pathfile:=path+filepathSeparator+"Statistic.cbs"
	err:=FileWrite(pathfile,bStatistic)
	Err(err,0)	
}

func GetStaticData(cubeno int,aStatistic map[string]StatisticData) map[string]StatisticData {
	var iData []TxData
	var sdata1 StatisticData
	var sdata2 StatisticData

	var c Cube 
	c.Cubeno=cubeno
	c.Read()

	for i:=-1;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			if i==-1 {
				tx1,tx2:=GetCubePoh(cubeno)
				iData=append(iData,tx1)
				iData=append(iData,tx2)
			} else {
				iData=BlockTxData(c.Blocks[i].Data)
			}
			for _,v := range iData {
				if v.Datatype!="NULL" && v.Datatype!="Data" && v.Datatype!="Contract" {
					if v.Datatype=="QUB" {
						if _, ok := aStatistic[v.From]; ok==false {
							aStatistic[v.From]=StatisticData{}
						}
						sdata1=aStatistic[v.From]
						sdata1.Cubeno=cubeno
						sdata1.SBalance=sdata1.SBalance+v.Amount+v.Fee+v.Tax
						sdata1.SBalance=math.Round(sdata1.SBalance*100000000)/100000000
						sdata1.Txcnt++
						if _, ok := aStatistic[v.To]; ok==false {
							aStatistic[v.To]=StatisticData{}
						}
						sdata2=aStatistic[v.To]
						sdata2.Cubeno=cubeno
						sdata2.RBalance=sdata2.RBalance+v.Amount+v.Tax
						sdata2.RBalance=math.Round(sdata2.RBalance*100000000)/100000000
						sdata2.Txcnt++
					} else {
						sdata1.Tkcnt++
						sdata2.Tkcnt++
					}
					sdata1.Balance=sdata1.RBalance-sdata1.SBalance
					sdata1.TBalance=sdata1.RBalance+sdata1.SBalance
					if sdata1.Balance>=5000.0 {
						sdata1.Pos="T"
					}	
					sdata2.Balance=sdata2.RBalance-sdata2.SBalance
					sdata2.TBalance=sdata2.RBalance+sdata2.SBalance
					if sdata2.Balance>=5000.0 {
						sdata2.Pos="T"
					}
					aStatistic[v.From]=sdata1
					aStatistic[v.To]=sdata2

				}
			}
		}
	}
	return aStatistic
}
