package core

import (
    "path/filepath"
	"os"
	"time"
	"strconv"
	"io/ioutil"
	"../lib"
	"../wallet"
)

func genesisBlock(Cubenum int) Block {
	var block=*new(Block)
	if !vaildCubeno(Cubenum) {
		return block
	}
	bFind:=blockFinder("0_"+strconv.Itoa(Cubenum))
	if bFind {
		return block
	}
	var txp=new(TxPool)
	var txd=new(TxData)
	txstr:="-1"
	amount:=Cubenum*1000
	w:=*wallet.CreateWallet()
	a:=w.GetAddress()
	addr:=lib.ByteToStr(a)
	tx:=txd.TxInput(w,txstr,addr,amount)
	txp.Tdata=append(txp.Tdata,tx)
	block.Index=0
	block.Cubeno=Cubenum
	block.Timestamp=int(time.Now().Unix())
	block.Nonce=10
	block.Data=Serialize(*txp)
	block.PrevCubeHash=cubeHash(Cubenum)
	block.PrevHash=prvHash(Cubenum)
	block.Hash=calculateHash(block)
	blockFile(block)
	return block
}

func addBlock(data []byte,Cubenum int) Block{
	var block Block
	var oblock Block

	c:=CurrentHeight()
	if c<1 {
		return block
	}
	bname:=blockName(strconv.Itoa(c-1)+"_"+strconv.Itoa(Cubenum))
	err:= fileRead(bname, oblock)
	lib.Err(err,0)
	block.Index=c
	block.Timestamp=int(time.Now().Unix())
	block.Nonce=10
	block.Cubeno=Cubenum
	block.Data=data
	block.PrevCubeHash=cubeHash(Cubenum) 
	block.PrevHash=oblock.Hash
	block.Hash=calculateHash(block)
	blockFile(block)
	return block
}


func calculateHash(block Block)  string {
	str := strconv.Itoa(block.Index) + strconv.Itoa(block.Cubeno) + strconv.Itoa(block.Timestamp) + lib.ByteToStr(block.Data) + block.PrevHash + block.PrevCubeHash + strconv.Itoa(block.Nonce)
	return setHash(str)
}

func cubeHash(cubenum int) string {
	return setHash(strconv.Itoa(cubenum))
}

func prvHash(cubenum int) string {
	return setHash(strconv.Itoa(cubenum))
}

func fileName(block *Block) string {
	filename:=strconv.Itoa(block.Index) + "_" + strconv.Itoa(block.Cubeno) + "_" + block.Hash + ".blc"
	return filename	
}

func blockFile(block Block) error {
	filename:=fileName(&block)
	err := fileWrite(filename, block)
	return err
}

func blockStrFile(block Block) error {
	filename:=strconv.Itoa(block.Index) + "_" + strconv.Itoa(block.Cubeno) + "_" + block.Hash + ".blc"
	str := strconv.Itoa(block.Index) + strconv.Itoa(block.Cubeno) + strconv.Itoa(block.Timestamp) + lib.ByteToStr(block.Data) + block.PrevHash + block.PrevCubeHash + strconv.Itoa(block.Nonce)
	bytes:=[]byte(str)
	path,_:=os.Executable()
	err:= ioutil.WriteFile(path+ string(filepath.Separator) +"bdata"+ string(filepath.Separator) +filename, bytes, 0)
	return err
}

func GetBalance(addr string) int {
	var amount=0
	var iBlock [27]Block
	var Txd TransactionData
	c:=CurrentHeight()-1
	for i:=0;i<27;i++ {
		err:=blockRead(c,i,iBlock[i])
		lib.Err(err,0)	
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
		} else {
			iData:=Deserialize(iBlock[i].Data)
			for _,v := range iData.Tdata {
				if v.DataType=="tx" {
					if(lib.ByteToStr(Txd.From)==addr) {
						amount+=Txd.Amount*(-1)
					}
					if(lib.ByteToStr(Txd.To)==addr) {
						amount+=Txd.Amount
					}
				}
			}
		}
	}
	return amount
}