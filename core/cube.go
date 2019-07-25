package core

import (
	"strings"
	"sync"
	"strconv"
	"time"
	"math"
	"os"
	"bytes"
	"net/http"
	"io/ioutil"
)

type Cube struct {
	Cubeno		int
	Timestamp	int
	Blocks		[27]Block
	PrevHash    string
	CHash       string
	Nonce       int
	Mine		CMineResult
}

type CMineResult struct {
	MineAddr	string
	MineMethod	string
	PowReward	float64
	PosReward	float64
	Txamount	float64
	Sumfee		float64
	Txcnt		int
	Tkcnt		int
	Concnt		int
	Hashcnt		int
	Difficulty	string
}


var waitGroup sync.WaitGroup


func (cube *Cube) Input(cubeno int) {
	cube.Cubeno=cubeno
	cube.Cubeno=cubeno
	cube.Timestamp=int(time.Now().Unix())
	var blocks [27]Block
	MakeDir(cubeno)
	difficulty()

	cube.Mine.MineAddr=Configure.Address
	cube.Mine.PowReward=Pratio.POW
	cube.Mine.PosReward=Pratio.POS
	if Configure.MiningMode=="pos" {
		cube.Mine.MineMethod="POS"
	} else {
		cube.Mine.MineMethod="POW"
	}
	cube.Mine.Difficulty=MineDifficulty

	Sumfee=0.0
	for i:=0;i<27;i++ {
		blocks[i].Input(cubeno,i+1)
		cube.Mine.Txamount+=blocks[i].Mine.Txamount
		cube.Mine.Sumfee+=blocks[i].Mine.Sumfee
		cube.Mine.Txcnt+=blocks[i].Mine.Txcnt
		cube.Mine.Tkcnt+=blocks[i].Mine.Tkcnt
		cube.Mine.Concnt+=blocks[i].Mine.Concnt
		cube.Mine.Hashcnt+=blocks[i].Nonce
	}
	cube.Mine.PowReward+=cube.Mine.Sumfee
	cube.Blocks=blocks
	cube.SetPrevHash()
	cube.SetHash()
	cube.PosMining()
	cube.Mine.Hashcnt+=cube.Nonce

	if cube.Nonce>0 {
		br:=cube.Broadcast()
		if br {
			cube.Save()
			decho ("save cube "+strconv.Itoa(cubeno))
			cube.FileBroadcast()
			go cube.MineBroadcast()
		}
	}
}

func (cube *Cube) InputChanel(cubeno int) {
	cube.Cubeno=cubeno
	cube.Timestamp=int(time.Now().Unix())
    bno:=make(chan int)
	var blocks [27]Block
	MakeDir(cubeno)
	difficulty()

	Sumfee=0.0
	for i:=0;i<27;i++ {
		waitGroup.Add(1)
		go blockInput(cubeno,bno,&blocks[i])
		bno<-(i+1)
	}
	
	waitGroup.Wait()
	
	cube.Mine.MineAddr=Configure.Address
	cube.Mine.PowReward=Pratio.POW
	cube.Mine.PosReward=Pratio.POS
	if Configure.MiningMode=="pos" {
		cube.Mine.MineMethod="POS"
	} else {
		cube.Mine.MineMethod="POW"
	}
	cube.Mine.Difficulty=MineDifficulty
	
	for i:=0;i<27;i++ {
		cube.Mine.Txamount+=blocks[i].Mine.Txamount
		cube.Mine.Sumfee+=blocks[i].Mine.Sumfee
		cube.Mine.Txcnt+=blocks[i].Mine.Txcnt
		cube.Mine.Tkcnt+=blocks[i].Mine.Tkcnt
		cube.Mine.Concnt+=blocks[i].Mine.Concnt
		cube.Mine.Hashcnt+=blocks[i].Nonce
	}
	cube.Mine.PowReward+=cube.Mine.Sumfee

	cube.Blocks=blocks
	cube.SetPrevHash()
	cube.SetHash()
	cube.PowMining()
	cube.Mine.Hashcnt+=cube.Nonce

	if cube.Nonce>0 {
		br:=cube.Broadcast()
		if br {
			cube.Save()
			decho ("save cube "+strconv.Itoa(cubeno))
			cube.FileBroadcast()
			go cube.MineBroadcast()
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
		cube.SetCubing(&PrvCubing)
		CubingFileWrite(PrvCubing)
	}
	return err
}

func (cube *Cube) Downloading() error {
	path:=FilePath(cube.Cubeno)
	url:="http://"+Configure.PosServer+"/download/"+strconv.Itoa(cube.Cubeno)
	_,err:=DownloadFile(path+filepathSeparator,url)
	if err!=nil {
		echo (err)
	}
	return err
}

func (cube *Cube) Download() error {
	path:=FilePath(cube.Cubeno)
	PathDelete(path)
	return cube.Downloading()
}

func (cube *Cube) FileBroadcast() error {
	filename:=cube.FileName()
	path:=FilePath(cube.Cubeno)
	file:=path+filepathSeparator+filename

	extraParams := map[string]string{
		"cmode": "filebroadcast",
		"cubeno": strconv.Itoa(cube.Cubeno),
	}
	request, err := newfileUploadRequest("http://"+Configure.PosServer+"/filebroadcast", extraParams, "file", file)
	if err != nil {
		echo(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		echo(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			echo(err)
		}
		resp.Body.Close()
		decho(resp.StatusCode)
		decho(resp.Header)
		echo(body)
	}
	return err
}

func (cube *Cube) FileBroadcast2()  {
	filename:=cube.FileName()
	path:=FilePath(cube.Cubeno)
	file, err := os.Open(path+filepathSeparator+filename)
	if err != nil {
		echo(err)
	}
	defer file.Close()
	res, err := http.Post("http://"+Configure.PosServer+"/filebroadcast", "binary/octet-stream", file)
	if err != nil {
		echo(err)
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	echo(string(message))
}


func (cube *Cube) String() string {
	bhash:=""
	phash:=""
	for i:=0;i<27;i++ {
		bhash+=cube.Blocks[i].Hash+","
		phash+=cube.Blocks[i].PatternHash+","
	}
	toStr:=strconv.Itoa(cube.Cubeno)+CubeDelim+strconv.Itoa(cube.Timestamp)+CubeDelim+bhash+CubeDelim+phash+CubeDelim+cube.PrevHash+CubeDelim+cube.CHash+CubeDelim+strconv.Itoa(cube.Nonce)+CubeDelim+cube.Mine.MineAddr
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

func (cube *Cube) CubeString() string {
	toStr:=strconv.Itoa(cube.Cubeno)+BlockDelim+strconv.Itoa(cube.Timestamp)+BlockDelim+cube.PrevHash+BlockDelim+cube.CHash+BlockDelim+strconv.Itoa(cube.Nonce)
	toStr+=BlockDelim+cube.Mine.MineAddr+BlockDelim+cube.Mine.MineMethod+BlockDelim+strconv.FormatFloat(cube.Mine.PowReward,'f',-1,64)+BlockDelim+strconv.FormatFloat(cube.Mine.PosReward,'f',-1,64)
	toStr+=BlockDelim+strconv.FormatFloat(cube.Mine.Txamount,'f',-1,64)+BlockDelim+strconv.FormatFloat(cube.Mine.Sumfee,'f',-1,64)
	toStr+=BlockDelim+strconv.Itoa(cube.Mine.Txcnt)+BlockDelim+strconv.Itoa(cube.Mine.Tkcnt)+BlockDelim+strconv.Itoa(cube.Mine.Concnt)+BlockDelim+strconv.Itoa(cube.Mine.Hashcnt)
	toStr+=BlockDelim+cube.Mine.Difficulty+BlockDelim+strconv.FormatInt(cube.FileSize(),10)
	return toStr
}

func (cube *Cube) SetPrevHash() {
	cube.PrevHash=cube.GetPrevHash()
}

func (cube *Cube) GetPrevHash() string {
	hashnip:=NodeCube("cubehash","0&cubeno="+strconv.Itoa(cube.Cubeno-1))
	if hashnip=="0,0" || hashnip=="" {
		return cube.GetPrevHash0();
	} else {
		cf:=CubeFileName(cube.Cubeno-1)
		haship:=strings.Split(hashnip, ",")
		if cf==strconv.Itoa(cube.Cubeno-1)+"_"+haship[0]+".cub" {
		} else if haship[1]>"0" {
		}
		return haship[0]
	}
}	

func (cube *Cube) GetPrevHash0() string {
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

func (cube *Cube) PosMining() {
	max:=100
	cube.CHash,cube.Nonce=PowCubeHashing(cube.CHash,max)
	cube.CHash="F"+cube.CHash[1:len(cube.CHash)]
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
			iData:=BlockTxData(cube.Blocks[i].Data)

			for _,v:=range iData {
				if v.Datatype=="QUB" {
					if v.From==addr {
						amount=amount-v.Amount-v.Fee-v.Tax
						amount=math.Round(amount*100000000)/100000000
					}
					if v.To==addr {
						amount+=v.Amount+v.Tax
						amount=math.Round(amount*100000000)/100000000
					}
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

