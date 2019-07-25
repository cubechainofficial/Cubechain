package core

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"
	"../config"
)

type Block struct {
	Cubeno		int
	Blockno		int
	Timestamp	int
	Data		[]byte
	Merkle		string
	Hash		string
	PrevHash	string
	PatternHash	string
	Nonce		int
	Mine		MineResult
}

type MineResult struct {
	TxMine		string
	Txamount	float64
	Txcnt		int
	Tkcnt		int
	Concnt		int
	Sumfee		float64
}


func (block *Block) SetData() {
}

func blockInput(cubeno int,bno chan int,block *Block) {
	for {
		blockno:=<-bno
		block.Input(cubeno,blockno)
		waitGroup.Done()
	}
}

func (block *Block) Input(cubeno int,blockno int) {
	var tbst TxBST
	block.Cubeno=cubeno
	block.Blockno=blockno
	block.Timestamp=int(time.Now().Unix())
	block.Nonce=0
	tbst,block.Mine=TxBlock(cubeno,blockno)
	block.Data=GetBytes(tbst)
	block.Merkle=setHash(tbst.Coin.Merkle())
	block.SetPrevHash()
	block.SetPatternHash()
	block.SetHash()

	ci:=strconv.Itoa(cubeno)+","+strconv.Itoa(blockno)

	decho ("Mining start "+ci)
	block.Mining()
	decho ("Mining end "+ci)
	if block.Nonce>0 {

		br:=block.Broadcast()
		if br==false {
		}

		block.Save()
		decho ("save block "+ci)
		block.MineBroadcast()
	}
	decho (block)
}

func (block *Block) InputString(Str string) {
	var td TxData
	td.Input(Str)
	block.Input(0,0)
}

func (block *Block) ToStamp() string {
	result:=NodeSend("block",block.String())
	return result
}

func (block *Block) Save() error {
	filename:=block.FileName()
	datapath:=block.FilePath()
	err:=FileWrite(datapath+string(filepathSeparator)+filename,block)
	return err
}

func (block *Block) Read() error {
	filename:=block.FileName()
	datapath:=block.FilePath()
	err:=FileRead(datapath+string(filepathSeparator)+filename,block)
	return err
}

func (block *Block) HashString() string {
	txd:=hex.EncodeToString(StrToByte(block.BlockTxString()))
	toStr:=strconv.Itoa(block.Cubeno)+BlockDelim+strconv.Itoa(block.Blockno)+BlockDelim+strconv.Itoa(block.Timestamp)+BlockDelim+txd+BlockDelim+block.Merkle+BlockDelim+block.PrevHash+BlockDelim+block.PatternHash
	return toStr
}

func (block *Block) String() string {
	txd:=hex.EncodeToString(StrToByte(block.BlockTxString()))
	toStr:=strconv.Itoa(block.Cubeno)+BlockDelim+strconv.Itoa(block.Blockno)+BlockDelim+strconv.Itoa(block.Timestamp)+BlockDelim+txd+BlockDelim+block.Merkle+BlockDelim+block.Hash+BlockDelim+block.PrevHash+BlockDelim+block.PatternHash+BlockDelim+strconv.Itoa(block.Nonce)
	return toStr
}

func (block *Block) MineString() string {
	txd:=hex.EncodeToString(StrToByte(block.BlockTxString()))
	txd=setHash(txd)
	toStr:=strconv.Itoa(block.Cubeno)+BlockDelim+strconv.Itoa(block.Blockno)+BlockDelim+strconv.Itoa(block.Timestamp)+BlockDelim+txd+BlockDelim+block.Merkle+BlockDelim+block.Hash+BlockDelim+block.PrevHash+BlockDelim+block.PatternHash+BlockDelim+strconv.Itoa(block.Nonce)
	toStr+=BlockDelim+block.Mine.TxMine+BlockDelim+strconv.FormatFloat(block.Mine.Txamount,'f',-1,64)+BlockDelim+strconv.Itoa(block.Mine.Txcnt)+BlockDelim+strconv.Itoa(block.Mine.Tkcnt)+BlockDelim+strconv.FormatFloat(block.Mine.Sumfee,'f',-1,64)
	toStr+=BlockDelim+strconv.FormatFloat(Pratio.BlockHash+block.Mine.Sumfee,'f',-1,64)+BlockDelim+Configure.Address+BlockDelim+strconv.FormatInt(block.FileSize(),10)
	return toStr
}

func (block *Block) MineDataString() string {
	txd:=hex.EncodeToString(StrToByte(block.BlockTxString()))
	toStr:=strconv.Itoa(block.Cubeno)+BlockDelim+strconv.Itoa(block.Blockno)+BlockDelim+strconv.Itoa(block.Timestamp)+BlockDelim+txd+BlockDelim+block.Merkle+BlockDelim+block.Hash+BlockDelim+block.PrevHash+BlockDelim+block.PatternHash+BlockDelim+strconv.Itoa(block.Nonce)
	toStr+=BlockDelim+block.Mine.TxMine+BlockDelim+strconv.FormatFloat(block.Mine.Txamount,'f',-1,64)+BlockDelim+strconv.Itoa(block.Mine.Txcnt)+BlockDelim+strconv.Itoa(block.Mine.Tkcnt)+BlockDelim+strconv.FormatFloat(block.Mine.Sumfee,'f',-1,64)
	toStr+=BlockDelim+strconv.FormatFloat(Pratio.BlockHash+block.Mine.Sumfee,'f',-1,64)+BlockDelim+Configure.Address+BlockDelim+strconv.FormatInt(block.FileSize(),10)
	return toStr
}

func (block *Block) BlockString() string {
	txd:=hex.EncodeToString(block.Data)
	txd=setHash(txd)
	toStr:=strconv.Itoa(block.Cubeno)+BlockDelim+strconv.Itoa(block.Blockno)+BlockDelim+strconv.Itoa(block.Timestamp)+BlockDelim+txd+BlockDelim+block.Merkle+BlockDelim+block.Hash+BlockDelim+block.PrevHash+BlockDelim+block.PatternHash+BlockDelim+strconv.Itoa(block.Nonce)
	toStr+=BlockDelim+strconv.FormatFloat(block.Mine.Txamount,'f',-1,64)+BlockDelim+strconv.Itoa(block.Mine.Txcnt)+BlockDelim+strconv.Itoa(block.Mine.Tkcnt)+BlockDelim+strconv.Itoa(block.Mine.Concnt)+BlockDelim+strconv.FormatFloat(block.Mine.Sumfee,'f',-1,64)
	toStr+=BlockDelim+strconv.FormatInt(block.FileSize(),10)
	return toStr
}

func (block *Block) BlockTxString() string {
	iData:=BlockTxData(block.Data)
	toStr:=""
	for _,v := range iData {
		toStr=v.String()+"%%%"
	}
	return toStr
}

func (block *Block) SetPrevHash() {
	block.PrevHash=block.GetPrevHash()
}

func (block *Block) GetPrevHash() string {
	hashnip:=NodeSend("blockhash","0&cubeno="+strconv.Itoa(block.Cubeno-1)+"&blockno="+strconv.Itoa(block.Blockno))
	if hashnip=="0,0" || hashnip=="" {
		return block.GetPrevHash0();
	} else {
		haship:=strings.Split(hashnip, ",")
		return haship[0]
	}
}

func GetPattenHash(cubeno int,blockno int) string {
	hashnip:=NodeSend("blockphash","0&cubeno="+strconv.Itoa(cubeno-1)+"&blockno="+strconv.Itoa(blockno))
	if hashnip=="0,0" || hashnip=="" {
		return GetPattenHash0(cubeno,blockno);
	} else {
		haship:=strings.Split(hashnip, ",")
		return haship[0]
	}
}

func (block *Block) GetPrevHash0() string {
	if block.Cubeno<2 {
		return setHash("GenesisBlockhash"+strconv.Itoa(block.Blockno))
	} else if PrvCubing.Cubeno==block.Cubeno-1 && PrvCubing.Hash1[block.Blockno-1]>"" {
	} else if BlockName(block.Cubeno-1,block.Blockno)>"" {
		return ReadBlockHash(block.Cubeno-1,block.Blockno)
	} else if CubingFileName(block.Cubeno-1)>"" {
		PrvCubing=CubingFileRead(block.Cubeno-1)
	} else if CurrCube.Cubeno==block.Cubeno-1 {
		CurrCube.SetCubing(&PrvCubing)
	} else if CubeFileName(block.Cubeno-1)>"" {
		CurrCube.Cubeno=block.Cubeno-1
		CurrCube.Read()
		CurrCube.SetCubing(&PrvCubing)
	} else {
		downpath:=CubeDownload(block.Cubeno-1)
		if downpath>"" {
			FileRead(downpath,&CurrCube)
			if CurrCube.Cubeno==block.Cubeno-1 {
				CurrCube.SetCubing(&PrvCubing)
			}  else {
				PrvCubing=GetCubing(block.Cubeno-1)
			}
		} else {
			PrvCubing=GetCubing(block.Cubeno-1)
		}
	} 
	return PrvCubing.Hash1[block.Blockno-1]
}



func GetPattenHash0(cubeno int,blockno int) string {
	if cubeno<2 {
		return setHash("GenesisPattenhash"+strconv.Itoa(blockno))
	} else if PrvCubing.Cubeno==cubeno {
	} else if BlockName(cubeno,blockno)>"" {
		return ReadBlockPHash(cubeno,blockno)
	} else if CubingFileName(cubeno)>"" {
		PrvCubing=CubingFileRead(cubeno)
	} else if CurrCube.Cubeno==cubeno {
		CurrCube.SetCubing(&PrvCubing)
	} else if CubeFileName(cubeno)>"" {
		CurrCube.Cubeno=cubeno
		CurrCube.Read()
		CurrCube.SetCubing(&PrvCubing)
	} else {
		downpath:=CubeDownload(blockno-1)
		if downpath>"" {
			FileRead(downpath,&CurrCube)
			if CurrCube.Cubeno==blockno-1 {
				CurrCube.SetCubing(&PrvCubing)
			}  else {
				PrvCubing=GetCubing(blockno-1)
			}
		} else {
			PrvCubing=GetCubing(blockno-1)
		}
	}
	return PrvCubing.Hash2[blockno-1]
}

func (block *Block) SetPatternHash() {
	block.PatternHash=block.GetPatternHash()
}

func (block *Block) GetPatternHash() string {
	pstr:=""
	l:=len(config.CubeCo[block.Blockno-1])
	for _,v:=range config.CubeCo[block.Blockno-1] {
		pstr+=GetPattenHash(block.Cubeno-1,v)
	}
	result:=PatternHash(pstr,l)
	return result
}

func (block *Block) SetHash() {
	block.Hash=block.GetHash()
}

func (block *Block) GetHash() string {
	hashstr:=block.HashString()
	result:=setHash(hashstr)
	return result
}

func (block *Block) Mining() {
	if Configure.MiningMode=="miningpool" {
		block.PoolMining()
	} else if Configure.MiningMode=="pos" {
		block.PosMining()
	} else {
		block.PowMining()
	}
}

func (block *Block) PoolMining() {
	phs:=PohSet(block.Cubeno)
	phs.Cubeno=block.Cubeno
	phs.Blockno=block.Blockno
	phs.HashStr=block.Hash	
	phs.Result(0)
	if phs.ResultHash>"" && phs.ResultNonce>0 {
		block.Nonce=phs.ResultNonce
		block.Hash=phs.ResultHash
	}
}

func (block *Block) PowMining() {
	max:=Configure.Maxnonce
	if block.Blockno==Configure.Indexing+1 || block.Blockno==Configure.Statistics+1 || block.Blockno==Configure.Escrow+1 || block.Blockno==Configure.Format+1 || block.Blockno==Configure.Edit+1 {
		block.Hash,block.Nonce=PowSpecialHashing(block.Hash,max)
	} else {
		block.Hash,block.Nonce=PowBlockHashing(block.Hash,max)
	}
}

func (block *Block) PosMining() {
	block.Nonce=((block.Cubeno*block.Blockno)*2+(block.Cubeno+block.Blockno)*3+block.Timestamp*4+5)%1000000000
	bh:=BlockHash(block.Hash+strconv.Itoa(block.Nonce))
	block.Hash="F"+bh[1:len(bh)]
}

func (block *Block) Verify() bool {
	pidx:=block.Cubeno-1
	if pidx>0 {
		if block.PrevHash!=block.GetPrevHash() {
			return false
		} else if block.Hash!=block.GetHash() {
			return false
		}
	}
	return true
}

func (block *Block) Broadcast() bool {
	result:=false
	r:=NodeCube("blocksave",block.String())
	decho (r)
	if r=="Success." {
		result=true
	}
	return result
}

func (block *Block) MineBroadcast() bool {
	result:=false
	r:=NodeSend("blocksave",block.MineString())
	decho (r)
	if r=="Success." {
		result=true
	}
	return result
}

func (block *Block) FileName() string {
	filename:=""
	if block.Hash=="" {
		filename=BlockName(block.Cubeno,block.Blockno)
	} else {
		filename=strconv.Itoa(block.Cubeno) + "_" + strconv.Itoa(block.Blockno) + "_" + block.Hash + ".blk"
	}
	return filename	
}

func (block *Block) FilePath() string {
	filepath:=FilePath(block.Cubeno)
	return filepath	
}

func (block *Block) FileSize() int64 {
	filepath:=block.FilePath()
	filename:=block.FileName()
	filesize:=FileSize(filepath+filepathSeparator+filename)
	return filesize	
}

func (block *Block) Print() {
	echo ("==============Block Head==============")
	echo ("Cubeno=",block.Cubeno)
	echo ("Blockno=",block.Blockno)
	echo ("Timestamp=",block.Timestamp)
	echo ("Merkle=",block.Merkle)
	echo ("Hash=",block.Hash)
	echo ("PrevHash=",block.PrevHash)
	echo ("Nonce=",block.Nonce)
	echo ("==============Block Body==============")
	txd:=string(block.Data)
	echo (txd)
}

func (block *Block) FileInfo() {
	echo ("==============File Info==============")
	echo ("FileName=",block.FileName())
	echo ("FilePath=",block.FilePath())
	echo ("FileSize=",block.FileSize())
}
