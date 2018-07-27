package core

import (
	"strconv"
	"../lib"
	"../config"
)

var Configure config.Configuration


func addIndexing() Block {
	c:=CurrentHeight()-1
	cnum:=Configure.Indexing;
	var iBlock [27]Block
	var Txd TransactionData
	var iAdd []IndexType
	var pAdd IBlock
	var CArr []CubeIndex
	CI:=CubeIndex{c,cnum}
	CArr=append(CArr,CI)

	for i:=0;i<27;i++ {
		err:=BlockRead(c,i,iBlock[i])
		lib.Err(err,0)	
		if i==Configure.Statistics || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
		} else if i==Configure.Indexing {
			pAdd=IndexDeserialize(iBlock[i].Data)
		} else {
			iData:=Deserialize(iBlock[i].Data)
			for _,v := range iData.Tdata {
				if v.DataType=="tx" {
					Txd=v.DataTx
					iAdd=append(iAdd,IndexType{lib.ByteToStr(Txd.From),CArr})
					iAdd=append(iAdd,IndexType{lib.ByteToStr(Txd.To),CArr})
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
	if ci.CubeNum==27 {
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
	if string(cno)=="A" {
		cbn=27
	} else {
		cbn=int(cno)-96
	}
	result:=CubeIndex{index,cbn}
	return result
}


func addStatistic() Block {
	c:=CurrentHeight()-1
	cnum:=Configure.Statistics;
	var iBlock [27]Block
	var pAdd SBlock
	var sAddr []string
	var sBal []int
	var iBal int
	var Txd TransactionData

	for i:=0;i<27;i++ {
		err:=BlockRead(c,i,iBlock[i])
		lib.Err(err,0)	
		if i==Configure.Indexing || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
		} else if i==cnum {
			pAdd=StaticDeserialize(iBlock[i].Data)
		} else {
			iData:=Deserialize(iBlock[i].Data)
			for _,v := range iData.Tdata {
				if v.DataType=="tx" {
					Txd=v.DataTx
					sAddr=append(sAddr,lib.ByteToStr(Txd.From))
					sBal=append(sBal,Txd.Amount*(-1))
					sAddr=append(sAddr,lib.ByteToStr(Txd.To))
					sBal=append(sBal,Txd.Amount)
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
	lib.Err(err,0)	
	pAdd=EscrowDeserialize(iBlock.Data)
	
	for k,v := range pAdd.Escrow {
		if v.State!=0 {
			pAdd.Escrow = append(pAdd.Escrow,pAdd.Escrow[k])
		}
	}
	return addBlock(Serialize(pAdd),c)
}
