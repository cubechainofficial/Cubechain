package core

import (
//	"fmt"
	//"strings"
	"sync"
	"strconv"
	"time"
	//"os"
	//"encoding/gob"
)

type Cube struct {
	Cubeno		int
	Timestamp	int
	Blocks		[27]Block
	PrevHash    string
	CHash       string
	Nonce       int
}

var waitGroup sync.WaitGroup


func (cube *Cube) Input(cubeno int) {
	cube.Cubeno=cubeno
	cube.Cubeno=cubeno
	cube.Timestamp=int(time.Now().Unix())
	var blocks [27]Block

	Sumfee=0.0
	for i:=0;i<27;i++ {
		blocks[i].Input(cubeno,i+1)
	}
	cube.Blocks=blocks
	cube.SetPrevHash()
	cube.SetHash()
	cube.PowMining()

	if cube.Nonce>0 {
		br:=cube.Broadcast()
		if br {
			cube.Save()
			decho ("save cube "+strconv.Itoa(cubeno))
			cube.MineBroadcast()
		}
	}
}

func (cube *Cube) InputChanel(cubeno int) {
	cube.Cubeno=cubeno
	cube.Cubeno=cubeno
	cube.Timestamp=int(time.Now().Unix())
    bno:=make(chan int)
	var blocks [27]Block
	waitGroup.Add(27)

	Sumfee=0.0
	for i:=0;i<27;i++ {
		go blockInput(cubeno,bno,&blocks[i])
		bno<-(i+1)
	}
	waitGroup.Wait()

	cube.Blocks=blocks
	cube.SetPrevHash()
	cube.SetHash()
	cube.PowMining()

	if cube.Nonce>0 {
		br:=cube.Broadcast()
		if br {
			cube.Save()
			decho ("save cube "+strconv.Itoa(cubeno))
			cube.MineBroadcast()
		}
	}
}


func (cube *Cube) SetCubing(cubing *Cubing) {
	cubing.Cubeno=cube.Cubeno
	cubing.Timestamp=cube.Timestamp
	cubing.PrevHash=cube.PrevHash
	cubing.CHash=cube.CHash
	cubing.Nonce=cube.Nonce
	for i:=0;i<27;i++ {
		cubing.Hash1[i]=cube.Blocks[i].Hash
		cubing.Hash2[i]=cube.Blocks[i].PatternHash
	}
}

func (cube *Cube) Read() error {
	filename:=cube.FileName()
	datapath:=cube.FilePath()
	decho(datapath+string(filepathSeparator)+filename)
	err:=FileRead(datapath+string(filepathSeparator)+filename,cube)
	decho(err)
	return err
}

func (cube *Cube) Save() error {
	filename:=cube.FileName()
	path:=FilePath(cube.Cubeno)
	err:=FileWrite(path+filepathSeparator+filename,cube)
	if err==nil {
		//var cubing Cubing
		cube.SetCubing(&PrvCubing)
		CubingFileWrite(PrvCubing)
	}
	return err
}

func (cube *Cube) String() string {
	bhash:=""
	phash:=""
	for i:=0;i<27;i++ {
		bhash+=cube.Blocks[i].Hash+","
		phash+=cube.Blocks[i].PatternHash+","
	}
	toStr:=strconv.Itoa(cube.Cubeno)+CubeDelim+strconv.Itoa(cube.Timestamp)+CubeDelim+bhash+CubeDelim+phash+CubeDelim+cube.PrevHash+CubeDelim+cube.CHash+CubeDelim+strconv.Itoa(cube.Nonce)
	return toStr
}

func (cube *Cube) MineString() string {
	bhash:=""
	phash:=""
	for i:=0;i<27;i++ {
		bhash+=cube.Blocks[i].Hash+","
		phash+=cube.Blocks[i].PatternHash+","
	}
	toStr:=strconv.Itoa(cube.Cubeno)+CubeDelim+strconv.Itoa(cube.Timestamp)+CubeDelim+bhash+CubeDelim+phash+CubeDelim+cube.PrevHash+CubeDelim+cube.CHash+CubeDelim+strconv.Itoa(cube.Nonce)
	toStr+=CubeDelim+strconv.FormatFloat(Pratio.BlockHash+Sumfee,'f',-1,64)+CubeDelim+Configure.Address+CubeDelim+strconv.FormatInt(cube.FileSize(),10)

	return toStr
}

func (cube *Cube) HashString() string {
	bhash:=""
	for i:=0;i<27;i++ {
		bhash+=cube.Blocks[i].Hash
		bhash+=cube.Blocks[i].PatternHash
	}
	toStr:=strconv.Itoa(cube.Cubeno)+CubeDelim+strconv.Itoa(cube.Timestamp)+CubeDelim+bhash+CubeDelim+cube.PrevHash
	return toStr
}

func (cube *Cube) SetPrevHash() {
	cube.PrevHash=cube.GetPrevHash()
}

func (cube *Cube) GetPrevHash() string {
	if cube.Cubeno<2 {
		return CubingHash("GenesisCubehash")
	} else if PrvCubing.Cubeno==cube.Cubeno-1 && PrvCubing.CHash>"" {
	} else if CubingFileName(cube.Cubeno-1)>"" {
		PrvCubing=CubingFileRead(cube.Cubeno-1)
	} else if CurrCube.Cubeno==cube.Cubeno-1 {
		CurrCube.SetCubing(&PrvCubing)
	} else if CubeFileName(cube.Cubeno-1)>"" {
		CurrCube.Cubeno=cube.Cubeno-1
		CurrCube.Read()
		CurrCube.SetCubing(&PrvCubing)
	} else {
		downpath:=CubeDownload(cube.Cubeno-1)
		if downpath>"" {
			FileRead(downpath,&CurrCube)
			if CurrCube.Cubeno==cube.Cubeno-1 {
				CurrCube.SetCubing(&PrvCubing)
			}  else {
				PrvCubing=GetCubing(cube.Cubeno-1)
			}             
		} else {
			PrvCubing=GetCubing(cube.Cubeno-1)
		}
	} 
	return PrvCubing.CHash
}

func (cube *Cube) SetHash() {
	cube.CHash=cube.GetHash()
}

func (cube *Cube) GetHash() string {
	cube.Nonce=0
	hashstr:=cube.HashString()
	result:=CallHash(hashstr,4)
	return result
}

func (cube *Cube) Mining() {
	phs:=PohSet(cube.Cubeno)
	phs.Cubeno=cube.Cubeno
	phs.Blockno=0
	phs.HashStr=cube.CHash	
	phs.Result(0)
	if phs.ResultHash>"" && phs.ResultNonce>0 {
		cube.Nonce=phs.ResultNonce
		cube.CHash=phs.ResultHash
	}
}

func (cube *Cube) PowMining() {
	max:=10000
	cube.CHash,cube.Nonce=PowCubeHashing(cube.CHash,max)
}

func (cube *Cube) Verify() bool {
	if cube.PrevHash!=cube.GetPrevHash() {
		return false
	} else if cube.CHash!=cube.GetHash() {
		return false
	} else {
		for i:=0;i<27;i++ {
			if cube.Blocks[i].Verify() {
				return false
			}
		}
	}
	return true
}

func (cube *Cube) Broadcast() bool {
	result:=false
	r:=NodeCube("cubesave",cube.String())
	decho (r)
	if r=="Success." {
		result=true
	}
	return result
}


func (cube *Cube) MineBroadcast() bool {
	result:=false
	r:=NodeSend("cubesave",cube.MineString())
	decho (r)
	if r=="success." {
		result=true
	}
	return result
}

func (cube *Cube) Balance(addr string) float64 {
	amount:=0.0
	for i:=0;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Statistics || i==Configure.Escrow || i==Configure.Format || i==Configure.Edit {
		} else {
			var TxArr []TxData
			iData:=TreeDeserialize(cube.Blocks[i].Data)
			iData.Coin.Convert(&TxArr);
			iData.Poh.Convert(&TxArr);
			for _,v:=range TxArr {
				if v.From==addr {
					amount+=v.Amount*(-1)
				}
				if v.To==addr {
					amount+=v.Amount
				}
			}
		}
	}
	return amount
}

func (cube *Cube) TxCount(addr string) (int,int) {
	var count=0
	var ecount=0
	for i:=0;i<27;i++ {
		if i==Configure.Indexing || i==Configure.Format || i==Configure.Edit || i==Configure.Statistics {
			var TxArr []TxData
			iData:=TreeDeserialize(cube.Blocks[i].Data)
			iData.Coin.Convert(&TxArr);
			iData.Poh.Convert(&TxArr);
			for _,v:=range TxArr {
				if v.From==addr {
					if i==Configure.Escrow { ecount++ } else { count++ }
				}
			}
		}
	}
	return count,ecount
}

func (cube *Cube) FileName() string {
	filename:=""
	if cube.CHash>"" {
		filename=cube.CHash+".cub"
	} else {
		filename=CubeFileName(cube.Cubeno)
	}
	return filename	
}

func (cube *Cube) FilePath() string {
	filepath:=FilePath(cube.Cubeno)
	return filepath	
}

func (cube *Cube) FileSize() int64 {
	filepath:=cube.FilePath()
	filename:=cube.FileName()
	filesize:=FileSize(filepath+filepathSeparator+filename)
	return filesize	
}

func (cube *Cube) Print() {
	echo ("==============Cube Head==============")
	echo ("Cubeno=",cube.Cubeno)
	echo ("Timestamp=",cube.Timestamp)
	echo ("PrevHash=",cube.PrevHash)
	echo ("CHash=",cube.CHash)
	echo ("==============Cube Blocks==============")

	for i:=0;i<27;i++ {
		j:=i+1
		echo ("==============Cube Blocks["+strconv.Itoa(j)+"]==============")
		cube.Blocks[i].Print()
	}
}

func (cube *Cube) FileInfo() {
	echo ("==============File Info==============")
	echo ("FileName=",cube.FileName())
	echo ("FilePath=",cube.FilePath())
	echo ("FileSize=",cube.FileSize())
}




