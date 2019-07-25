package core

import (
	"strings"
	"strconv"
	"math"
)

func CubeBalance(addr string,cubeno int) float64 {
	var c Cube 
	c.Cubeno=cubeno
	c.Read()

	return c.Balance(addr)
}

func GetBalance(addr string) float64 {
	result:=0.0
	ci:=GetIndexBlock(addr)
	for _,v:=range ci {
		vv:=TxScanArr(addr,v)
		for _,v:=range vv {
			if v.Datatype=="QUB" {
				if v.From==addr {
					result=result-v.Amount-v.Fee-v.Tax
					result=math.Round(result*100000000)/100000000
				} else if v.To==addr {
					result=result+v.Amount+v.Tax
					result=math.Round(result*100000000)/100000000
				}
			}
		}
	}
	return result
}

func GetBalanceRS(addr string,rs string) float64 {
	result:=0.0
	ci:=GetIndexBlock(addr)
	for _,v:=range ci {
		vv:=TxScanArr(addr,v)
		for _,v:=range vv {
			if v.Datatype=="QUB" {
				if v.From==addr && rs=="From" {
					result=result+v.Amount+v.Fee+v.Tax
					result=math.Round(result*100000000)/100000000
				} else if v.To==addr && rs=="To" {
					result=result+v.Amount+v.Tax
					result=math.Round(result*100000000)/100000000
				}
			}
		}
	}
	return result
}

func GetTransactionCount(addr string) int {
	result:=0
	ci:=GetIndexBlock(addr)
	for _,v:=range ci {
		vv:=TxScanArr(addr,v)
		for _,v:=range vv {
			if v.Datatype=="QUB" {
				if v.From==addr || v.To==addr {
					result++
				}
			}
		}
	}
	return result
}

func GetStaticValue(addr string) (float64,float64,int) {
	result1:=0.0
	result2:=0.0
	cnt:=0
	ci:=GetIndexBlock(addr)
	for _,v:=range ci {
		vv:=TxScanArr(addr,v)
		for _,v:=range vv {
			if v.Datatype=="QUB" {
				if v.From==addr {
					result1=result1+v.Amount+v.Fee+v.Tax
					result1=math.Round(result1*100000000)/100000000
					cnt++
				} else if v.To==addr {
					result2=result2+v.Amount+v.Tax
					result2=math.Round(result2*100000000)/100000000
					cnt++
				}
			}
		}
	}
	return result2,result1,cnt
}

func GetStaticVar(addr string) (float64,float64,float64,float64,int,int,int) {
	gs:=GetStatisticAddr(addr)
	gsitem:=strings.Split(gs,"||")
	gsv:=strings.Split(gsitem[1],",")
	
	balance,_:=strconv.ParseFloat(gsv[0],64)
	rbalance,_:=strconv.ParseFloat(gsv[1],64)
	sbalance,_:=strconv.ParseFloat(gsv[2],64)
	tbalance,_:=strconv.ParseFloat(gsv[3],64)
	txcnt,_:=strconv.Atoi(gsv[4])
	tkcnt,_:=strconv.Atoi(gsv[5])
	gcubeno,_:=strconv.Atoi(gsitem[0])

	return balance,rbalance,sbalance,tbalance,txcnt,tkcnt,gcubeno
}

func GetStaticVarCube(addr string,cubeno int) (float64,float64,int,int) {
	var iData []TxData

	rbalance:=0.0
	sbalance:=0.0
	txcnt:=0
	tkcnt:=0

	if cubeno<=0 || addr[0:31]=="C"+strings.Repeat("0",30) {
		return rbalance,sbalance,txcnt,tkcnt
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
				if v.Datatype!="NULL" && v.Datatype!="Data" && v.Datatype!="Contract" {
					if v.Datatype=="QUB" {
						if v.From==addr {
							sbalance=sbalance+v.Amount+v.Fee+v.Tax
							sbalance=math.Round(sbalance*100000000)/100000000
							txcnt++
						}
						if v.To==addr {
							rbalance=rbalance+v.Amount+v.Tax
							rbalance=math.Round(rbalance*100000000)/100000000
							txcnt++
						}
					} else {
						tkcnt++
					}
				}
			}
		}
	}
	return rbalance,sbalance,txcnt,tkcnt
}

func GetTransactionList(addr string) string {
	result:=""
	var res []string
	ci:=GetIndexBlock(addr)
	for _,v:=range ci {
		vv:=TxScanArr(addr,v)
		for _,v:=range vv {
			if v.From==addr || v.To==addr {
				res=append(res,v.Hash)
			}
		}
	}
	result=strings.Join(res,",")
	return result
}

func GetTxListDetail(addr string) string {
	result:=""
	var res []string
	ci:=GetIndexBlock(addr)
	for _,v:=range ci {
		vv:=TxScanArr(addr,v)
		for _,v:=range vv {
			if v.From==addr || v.To==addr {
				res=append(res,v.String())
			}
		}
	}
	result=strings.Join(res,",")
	return result
}

func GetTransactionDetail(txhash string) (TxData,CubeIndex) {
	var tdata TxData
	c:=CurrentHeight()-1
	for i:=1;i<=c;i++ {
		cdata,j:=CubeScanHash(txhash,i,"QUB")
		if cdata.Datatype!="NULL" && cdata.Datatype!="Data" && cdata.Datatype!="Contract" {
			return cdata,CubeIndex{i,j}
		}
	}
	return tdata,CubeIndex{0,0}
}

func GetTransactionData(txhash string) (TxData,CubeIndex) {
	var tdata TxData
	c:=CurrentHeight()-1
	for i:=1;i<=c;i++ {
		cdata,j:=CubeScanHash(txhash,i,"Data")
		if cdata.Datatype!="NULL" && cdata.Datatype!="Data" && cdata.Datatype!="Contract" {
			return cdata,CubeIndex{i,j}
		}
	}
	return tdata,CubeIndex{0,0}
}

func GetTxCount(addr string) int {
	c:=CurrentHeight()-1
	ci:=GetIndexBlock(addr)
	result:=len(ci)
	cdata:=CubeScan(addr,c)
	if cdata.Datatype=="QUB" {
		result++
	}
	return result
}

func GetTkCount(addr string) int {
	c:=CurrentHeight()-1
	ci:=GetIndexBlock(addr)
	result:=len(ci)
	cdata:=CubeScan(addr,c)
	if cdata.Datatype!="QUB" && cdata.Datatype!="NULL" && cdata.Datatype!="Data" && cdata.Datatype!="Contract" {
		result++
	}
	return result
}

func GetTxList(addr string,symbol string) string {
	result:=""
	var res []string
	c:=CurrentHeight()-1
	ci:=GetIndexBlock(addr)
	for _,v:=range ci {
		vv:=TxScan(addr,v)
		if vv.From==addr || vv.To==addr {
			res=append(res,vv.Hash)
		}
	}
	cdata:=CubeScan(addr,c)
	if cdata.Datatype==symbol {
		res=append(res,cdata.Hash)
	}
	result=strings.Join(res,",")
	return result
}

func GetTxDetail(txhash string,symbol string) (TxData,CubeIndex) {
	var tdata TxData
	c:=CurrentHeight()-1
	for i:=1;i<=c;i++ {
		cdata,j:=CubeScanHash(txhash,i,symbol)
		if cdata.Datatype==symbol {
			return cdata,CubeIndex{i,j}
		}
	}
	return tdata,CubeIndex{0,0}
}

func GeTxData(txhash string,symbol string) (TxData,CubeIndex) {
	var tdata TxData
	c:=CurrentHeight()-1
	for i:=1;i<=c;i++ {
		cdata,j:=CubeScanHash(txhash,i,"Data")
		if cdata.Datatype==symbol {
			return cdata,CubeIndex{i,j}
		}
	}
	return tdata,CubeIndex{0,0}
}

func CubeScan(addr string,idx int) TxData {
	var iTxData TxData
	var ci CubeIndex
	for i:=-1;i<27;i++ {
		ci=CubeIndex{idx,i+1}
		iTxData=TxScan(addr,ci)
		if iTxData.Datatype!="NULL" {
			return iTxData
		}
	}	
	return iTxData
}

func GetIndexBlock(addr string) []CubeIndex {
	var CSindex []CubeIndex
	var CTindex []CubeIndex
	Cindex:=CubeIndex{0,0}
	indexcube,indexStr:=IndexingAddr(addr)
	if strings.Index(indexStr,",")>0 {
		cin:=strings.Split(indexStr, ",")
		for _,v:=range cin {
			if len(v)==0 {
				continue
			}
			Cindex=StrToCubeIndex(v)
			CSindex=append(CSindex,Cindex)
		}
	}

	c:=CurrentHeight()-1
	for indexcube<c {
		CTindex=GetIndexCube(addr,indexcube)
		for _,v:=range CTindex {
			CSindex=append(CSindex,v)
		}
		indexcube++
	}
	return CSindex
}

func GetIndexCube(addr string,cubeno int) []CubeIndex {
	var CSindex []CubeIndex
	Cindex:=CubeIndex{0,0}
	for i:=-1;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
		} else {
			Cindex.Index=cubeno
			Cindex.CubeNum=i+1
			if TxScanResult(addr,Cindex) {
				CSindex=append(CSindex,Cindex)
			}
		}
	}	
	return CSindex
}


func GetIndexBlock_direct(addr string) []CubeIndex {
	var CSindex []CubeIndex
	Cindex:=CubeIndex{0,0}
	c:=CurrentHeight()-1
	for j:=1;j<c+1;j++ {
		for i:=-1;i<27;i++ {
			if i==Configure.Indexing || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
			} else {
				Cindex.Index=j
				Cindex.CubeNum=i+1
				if TxScanResult(addr,Cindex) {
					CSindex=append(CSindex,Cindex)
				}
			}
		}	
	}
	return CSindex
}

func TxScan(addr string,ci CubeIndex) TxData {
	var iBlock Block
	var iBlockData []TxData
	var iTxData TxData
	iTxData.Datatype="NULL"

	if ci.CubeNum==0 {
		tx1,tx2:=GetCubePoh(ci.Index)
		iBlockData=append(iBlockData,tx1)
		iBlockData=append(iBlockData,tx2)
	} else {
		err:=BlockRead(ci.Index,ci.CubeNum,&iBlock)
		Err(err,0)	
		iBlockData=BlockTxData(iBlock.Data)
	}

	for _,v:=range iBlockData {
		if v.From==addr || v.To==addr {
			return v
		}
	}
	return iTxData
}

func TxScanArr(addr string,ci CubeIndex) []TxData {
	var iBlock Block
	var iBlockData []TxData
	var iTxData []TxData
	if ci.CubeNum==0 {
		tx1,tx2:=GetCubePoh(ci.Index)
		iBlockData=append(iBlockData,tx1)
		iBlockData=append(iBlockData,tx2)
	} else {
		err:=BlockRead(ci.Index,ci.CubeNum,&iBlock)
		Err(err,0)	
		iBlockData=BlockTxData(iBlock.Data)
	}
	for _,v:=range iBlockData {
		if v.From==addr || v.To==addr {
			iTxData=append(iTxData,v)
		}
	}
	return iTxData
}

func TxScanResult(addr string,ci CubeIndex) bool {
	var iBlock Block
	var iBlockData []TxData
	var iTxData TxData
	iTxData.Datatype="NULL"
	if ci.CubeNum==0 {
		tx1,tx2:=GetCubePoh(ci.Index)
		iBlockData=append(iBlockData,tx1)
		iBlockData=append(iBlockData,tx2)
	} else {
		err:=BlockRead(ci.Index,ci.CubeNum,&iBlock)
		Err(err,0)	
		iBlockData=BlockTxData(iBlock.Data)
	}
	for _,v:=range iBlockData {
		if v.From==addr || v.To==addr {
			return true
		}
	}
	return false
}

func CubeScanHash(txhash string,idx int,datatype string) (TxData,int) {
	var iTxData TxData
	var ci CubeIndex
	j:=0
	for i:=-1;i<27;i++ {
		ci=CubeIndex{idx,i+1}
		iTxData,j=BlockScanHash(txhash,ci,datatype)
		if j>0 {
			return iTxData,j
		}
	}	
	return iTxData,0
}

func BlockScanHash(txhash string,ci CubeIndex,datatype string) (TxData,int) {
	var iBlock Block
	var iTxData TxData

	if ci.CubeNum==0 {
		tx1,tx2:=GetCubePoh(ci.Index)
		if tx1.Hash==txhash {
			return tx1,ci.CubeNum
		} else if tx1.Hash==txhash {
			return tx2,ci.CubeNum
		} 
		return iTxData,0
	} else {
		err:=BlockRead(ci.Index,ci.CubeNum,&iBlock)
		Err(err,0)	
		iBlockData:=BlockTree(iBlock.Data)
		coinData:=iBlockData.Coin

		if coinData!=nil {
			txdata,_:=coinData.Search(txhash) 
			iTxData,_:=txdata.(TxData)
			return iTxData,ci.CubeNum
		} else {
			return iTxData,0
		}
	}
}

func GetIndexing() string {
	result:=""
	var aIndexing TxIndexing
	aIndexing.AddrIndex=make(map[string]string)	

	IndexingRead(&aIndexing)

	for k,v:=range aIndexing.AddrIndex {
		result+=k+":"+v+"||"
	}
	result=strconv.Itoa(aIndexing.Cubeno)+"||"+result
	result=strings.Replace(result,",,",",",-1)
	return result
}

func GetIndexingAddr(addr string) string {
	var aIndexing TxIndexing
	aIndexing.AddrIndex=make(map[string]string)	
	IndexingRead(&aIndexing)
	result:=strconv.Itoa(aIndexing.Cubeno)+"||"+addr+"||"+strings.Replace(aIndexing.AddrIndex[addr],",,",",",-1)
	return result
}

func GetStatistic() string {
	result:=""
	var aStatistic TxStatistic
	aStatistic.AddrIndex=make(map[string]string)	
	StatisticRead(&aStatistic)
	for k,v:=range aStatistic.AddrIndex {
		result+=k+":"+v+"||"
	}
	result=strconv.Itoa(aStatistic.Cubeno)+"||"+result
	return result
}

func GetStatisticAddr(addr string) string {
	var aStatistic TxStatistic
	aStatistic.AddrIndex=make(map[string]string)	
	StatisticRead(&aStatistic)
	if aStatistic.AddrIndex[addr]=="" {
		aStatistic.AddrIndex[addr]="0.0,0.0,0.0,0,0,F"
	}
	result:=strconv.Itoa(aStatistic.Cubeno)+"||"+aStatistic.AddrIndex[addr]
	return result
}

func GetStatisticRank(item string,item2 string,count int) string {
	result:=""
	tsplit:=0
	var aStatistic TxStatistic
	aStatistic.AddrIndex=make(map[string]string)	
	iStatistic:=make(map[string]int)
	fStatistic:=make(map[string]float64)

	StatisticRead(&aStatistic)

	switch item {
		case "balance": tsplit=0
		case "rbalance": tsplit=1
		case "sbalance": tsplit=2
		case "tbalance": tsplit=3
		case "txcnt": tsplit=4
		case "tkcnt": tsplit=5
		default : tsplit=0
	}
	for k,v:=range aStatistic.AddrIndex {
		sitem:=strings.Split(v,",")
		if tsplit>3 {
			iStatistic[k],_=strconv.Atoi(sitem[tsplit])
		} else {
			fStatistic[k],_=strconv.ParseFloat(sitem[tsplit],64)
		}
	}
	if tsplit>3 {
		p:=PairSortInt(iStatistic)
		for k,v:=range p {
			if k<count && v.Value>0 {
				if item2>"" {
					result+=v.Key+","+strconv.Itoa(v.Value)+","+GetStatisticValue(item2,v.Key,&aStatistic)+","+strconv.Itoa(k+1)+"||"
				} else {
					result+=v.Key+","+strconv.Itoa(v.Value)+","+strconv.Itoa(k+1)+"||"
				}
			}
		}
	} else {
		p:=PairSortFloat(fStatistic)
		for k,v:=range p {
			if k<count && v.Value>0.0 {
				if item2>"" {
					result+=v.Key+","+strconv.FormatFloat(v.Value,'f',-1,64)+","+GetStatisticValue(item2,v.Key,&aStatistic)+","+strconv.Itoa(k+1)+"||"
				} else {
					result+=v.Key+","+strconv.FormatFloat(v.Value,'f',-1,64)+","+strconv.Itoa(k+1)+"||"
				}
			}
		}		
	}
	result=strconv.Itoa(aStatistic.Cubeno)+"||"+result
	return result
}

func GetStatisticValue(item string,addr string,aStatistic *TxStatistic) string {
	var result string
	tsplit:=0
	switch item {
		case "balance": tsplit=0
		case "rbalance": tsplit=1
		case "sbalance": tsplit=2
		case "tbalance": tsplit=3
		case "txcnt": tsplit=4
		case "tkcnt": tsplit=5
		default : tsplit=0
	}
	sitem:=strings.Split(aStatistic.AddrIndex[addr],",")
	result=sitem[tsplit]
	return result
}

func GetIssue() string {
	c:=CurrentHeight()-2
	issue:=182.65*float64(c)+2400000000.0
	result:=strconv.Itoa(c-1)+"||"+strconv.FormatFloat(issue,'f',-1,64)
	return result
}

func GetBlock(cubeno int,blockno int) string {
	var iBlock Block
	err:=BlockRead(cubeno,blockno,&iBlock)
	if err!=nil {
		return ""
	} else {
		return iBlock.BlockString()
	}
}

func GetBlockAll(cubeno int) string {
	result:=""
	var res []string
	for i:=0;i<27;i++ {
		res=append(res,GetBlock(cubeno,i+1))
	}
	result=strings.Join(res,",")
	return result
}



func GetBlockBase(cubeno int,blockno int) string {
	var iBlock Block
	err:=BlockRead(cubeno,blockno,&iBlock)
	Err(err,0)	
	return iBlock.String()
}

func GetCube(cubeno int) string {
	var c Cube 
	c.Cubeno=cubeno
	c.Read()
	return c.CubeString()
}

func GetCubeBase(cubeno int) string {
	var c Cube 
	c.Cubeno=cubeno
	c.Read()
	return c.String()
}

func GetCubePoh(cubeno int) (TxData,TxData) {
	var c Cube 
	var txd1 TxData
	var txd2 TxData
	c.Cubeno=cubeno
	c.Read()
	
	txd1.Timestamp=c.Timestamp
	txd1.From="C000000000000000000000000000000000"
	txd1.To=c.Mine.MineAddr
	txd1.Amount=c.Mine.PowReward
	txd1.Fee=0.0
	txd1.Tax=0.0
	txd1.Sign=setHash(c.CHash)
	txd1.Nonce=c.Mine.Hashcnt
	txd1.Message=""
	txd1.Datatype="QUB"
	txd1.SetHash()
	
	txd2.Timestamp=c.Timestamp
	txd2.From="C000000000000000000000000000000001"
	txd2.To="CPNpEb8jgwTS51f68DJsUu51WZV4HNTM4u"
	txd2.Amount=c.Mine.PosReward
	txd2.Fee=0.0
	txd2.Tax=0.0
	txd2.Sign=setHash(c.CHash)
	txd2.Nonce=c.Mine.Hashcnt
	txd2.Message=""
	txd2.Datatype="QUB"
	txd2.SetHash()

	return txd1,txd2
}


func GetBlockTx(cubeno int,blockno int) string {
	result:=""
	var res []string
	var iBlock Block

	if blockno==0 {
		tx1,tx2:=GetCubePoh(cubeno)
		res=append(res,strconv.Itoa(blockno)+"||"+tx1.TxString())
		res=append(res,strconv.Itoa(blockno)+"||"+tx2.TxString())
	} else {
		err:=BlockRead(cubeno,blockno,&iBlock)
		Err(err,0)	
		iBlockData:=BlockTxData(iBlock.Data)
		for _,v:=range iBlockData {
			res=append(res,strconv.Itoa(iBlock.Blockno)+"||"+v.TxString())
		}
	}
	result=strings.Join(res,",")
	return result
}

func GetBlockTxSearch(cubeno int,blockno int,addr string,u_addr string,coin float64) string {
	result:=""
	var res []string
	var iBlockData []TxData
	var iBlock Block
	
	if blockno==0 {
		tx1,tx2:=GetCubePoh(cubeno)
		iBlockData=append(iBlockData,tx1)
		iBlockData=append(iBlockData,tx2)
	} else {
		err:=BlockRead(cubeno,blockno,&iBlock)
		Err(err,0)	
		iBlockData=BlockTxData(iBlock.Data)
	}
	for _,v:=range iBlockData {
		if (v.From==addr || addr=="") && (v.To==u_addr || u_addr=="") && (v.Amount==coin || coin==0.0) {
			res=append(res,v.TxString())
		}
	}
	result=strings.Join(res,",")
	return result
}

func GetBlockTxHash(cubeno int,blockno int,addr string,u_addr string,coin float64) string {
	result:=""
	var res []string
	var iBlockData []TxData
	var iBlock Block
	if blockno==0 {
		tx1,tx2:=GetCubePoh(cubeno)
		iBlockData=append(iBlockData,tx1)
		iBlockData=append(iBlockData,tx2)
	} else {
		err:=BlockRead(cubeno,blockno,&iBlock)
		Err(err,0)	
		iBlockData=BlockTxData(iBlock.Data)
	}
	for _,v:=range iBlockData {
		if (v.From==addr || addr=="") && (v.To==u_addr || u_addr=="") && (v.Amount==coin || coin==0.0) {
			res=append(res,v.Hash)
		}
	}
	result=strings.Join(res,",")
	return result
}

func GetCubeTx(cubeno int) string {
	result:=""
	var res []string
	
	tx1,tx2:=GetCubePoh(cubeno)
	res=append(res,"0||"+tx1.TxString())
	res=append(res,"0||"+tx2.TxString())
	for i:=0;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Format || i==Configure.Edit {
		} else {
			res=append(res,GetBlockTx(cubeno,i+1))
		}
	}
	result=strings.Join(res,",")
	return result
}


