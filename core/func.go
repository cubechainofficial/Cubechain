package core

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
    "bytes"
	"net"
	"net/http"
	"io"
	"io/ioutil"
    "path/filepath"
	"../config"
)

var Configure config.Configuration
var filepathSeparator=string(filepath.Separator)
var echo=fmt.Println
var MineCubeno=0
var TxMine string
var TxMerkleHash string
var Sumfee=0.0
var Txamount=0.0
var Txcnt=0
var Tkcnt=0
var CubeSetNum [27]string
var DebugMode=false

var Version="0.98"
//var cubechainInfo Cubechain
var Pratio=Pohr{4.566,4.566,4.566,4.566}
var TxDelim="|"
var BlockDelim="||"
var CubeDelim="|||"
var CNO=1
var PrvCubing Cubing
var CurrCube Cube

func Err(err error, exit int) int {
	if err != nil {
		fmt.Println(err)	
	}
	if exit>=1 {
		os.Exit(exit)
		return 1
	}
	return 0
}

func netError(err error) {
	if err!=nil && err!=io.EOF {
		fmt.Println("Network Error : ", err)
	}
}

func IpCheck() []string {
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil
	}
	addrs=append(addrs,host)
	if len(addrs)==2 {
		addrs2:=make([]string,3)
		addrs2[0]="mac_linux"
		addrs2[1]=addrs[0]
		addrs2[2]=addrs[1]
		addrs=addrs2
	}
	return addrs
}

func GetBytes(key interface{}) []byte {
    var buf bytes.Buffer
	var Tdata TxData
	var Tbst TxBST

	gob.Register(Tdata)  
	gob.Register(Tbst)
	enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
		decho(err)
        return nil
    }
    return buf.Bytes()
}



func GetCubeHeight() string {
	result:=NodeSend("cubeheight","0")
	return result
}

func GetCubeHeight2() string {
	result:=NodeSend2("cubeheight","0")
	return result
}

func GetCubeHeight3() string {
	result:=NodeCube("cubeheight","0")
	return result
}

func CubeHeight() int {
	result,_:=strconv.Atoi(GetCubeHeight3())
	return result
}

func CurrentHeight() int {
	result:=0
	f:=MaxFind(Configure.Datafolder+filepathSeparator)
	if f=="0" {
		return 1
	}
	f2:=MaxFind(Configure.Datafolder+filepathSeparator+f)
	if f2=="0" {
		return 1
	}
	nint,_:=strconv.ParseUint(f,16,32)
	mint,_:=strconv.ParseUint(f2,16,32)
	result=(int(nint)-1)*Configure.Datanumber+int(mint)
	if FileSearch(FilePath(result),".cub")>"" {
		result++
	}
	return result	
}

func GetTxCount(addr string) int {
	return 1
}


func BlockName(idx int,cno int) string {
	find:=strconv.Itoa(idx)+"_"+strconv.Itoa(cno)+"_"	
    dirname:=FilePath(idx)
	result:=FileSearch(dirname,find)
	return result
}

func BlockRead(index int,cubeno int,object interface{}) error {
	filename:=BlockName(index,cubeno)
	datapath:=FilePath(index)+filepathSeparator 

	if filename=="" {
		//fmt.Println(strconv.Itoa(index)+":"+strconv.Itoa(cubeno))
		return nil
	}
	file,err:=os.Open(datapath+filename)
	if err==nil {
		decoder:=gob.NewDecoder(file)
		err=decoder.Decode(object)
	}
	file.Close()
	
	return err
}

func BlockScan(cubeno int,blockno int) Block {
	var block Block
	BlockRead(cubeno,blockno,&block)
	return block
}



func ReadBlockHash(index int,cubeno int) string {
	var iBlock Block
	var hash string
	err:=BlockRead(index,cubeno,&iBlock)
	Err(err,0)
	hash=iBlock.Hash
	return hash
}

func ReadBlockPHash(index int,cubeno int) string {
	var iBlock Block
	var hash string
	err:=BlockRead(index,cubeno,&iBlock)
	Err(err,0)
	hash=iBlock.PatternHash
	return hash
}


func CubeRead(index int,object *Cube) error {
	datapath:=FilePath(index)+filepathSeparator 
	filename:=FileSearch(datapath,".cub")
	file,err:=os.Open(datapath+filename)
	if err==nil {
		decoder:=gob.NewDecoder(file)
		err=decoder.Decode(object)
	}
	file.Close()
	return err
}

func CubeFileName(idx int) string {
	find:=".cub"	
    dirname:=FilePath(idx)
	result:=FileSearch(dirname,find)
	return result
}

/*
func GetTxCount(addr string) int {
	c,count:=0,0
	var block Block
	if c<=0 {
		c=CurrentHeight()-1
	}
	for i:=0;i<c;i++ {
		mblock.Index=i
		err:=block.Read()
		Err(err,0)
		if(Block.Data.From==addr) {
			count++
		}
	}
	return count
}

*/




func NodeSend(cmode string,data string) string {
	arr:=IpCheck()
	reader:=strings.NewReader("cmode="+cmode+"&_token=9X1rK2Z2sofIeFpqg6VBXI5aUWsPOfGPGyzzztgu&data="+data+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&netname="+Configure.Network+"&netset="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&netport="+strconv.Itoa(Configure.Port)+"&ver="+Version)
	request,_:=http.NewRequest("POST","http://"+Configure.MainServer+"/"+cmode, reader)
	request.Header.Add("content-type","application/x-www-form-urlencoded")
	request.Header.Add("cache-control","no-cache")
	client:=&http.Client{}
	res, err := client.Do(request)
	Err(err,0)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	Err(err,0)
	s:=string(body)
	return s
}

func NodeSend2(cmode string,data string) string {
	arr:=IpCheck()
	reader:=strings.NewReader("cmode="+cmode+"&_token=9X1rK2Z2sofIeFpqg6VBXI5aUWsPOfGPGyzzztgu&data="+data+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&netname="+Configure.Network+"&netset="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&netport="+strconv.Itoa(Configure.Port)+"&ver="+Version)
	request,_:=http.NewRequest("POST","http://"+Configure.PoolServer+"/"+cmode, reader)
	request.Header.Add("content-type","application/x-www-form-urlencoded")
	request.Header.Add("cache-control","no-cache")
	client:=&http.Client{}
	res, err := client.Do(request)
	Err(err,0)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	Err(err,0)
	s:=string(body)
	return s
}

func NodeCube(cmode string,data string) string {
	arr:=IpCheck()
	reader:=strings.NewReader("cmode="+cmode+"&_token=9X1rK2Z2sofIeFpqg6VBXI5aUWsPOfGPGyzzztgu&data="+data+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&netname="+Configure.Network+"&netset="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&netport="+strconv.Itoa(Configure.Port)+"&ver="+Version)
	request,_:=http.NewRequest("POST","http://"+Configure.PosServer+"/"+cmode, reader)
	request.Header.Add("content-type","application/x-www-form-urlencoded")
	request.Header.Add("cache-control","no-cache")
	client:=&http.Client{}
	res, err := client.Do(request)
	Err(err,0)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	Err(err,0)
	s:=string(body)
	return s
}

func CubeDownload(cubeno int) string {
	decho("Download file")
	cubenum:=strconv.Itoa(cubeno)
	//filename:="test"+cubenum
    fileUrl:="http://"+Configure.Rpcip+"/download/"+cubenum
    filepath,err := DownloadFile(FilePath(cubeno)+filepathSeparator, fileUrl)
    if err != nil {
		decho(err)
    }
	return filepath
}

func DownloadFile(filepath string, url string) (string,error) {
    resp,err:=http.Get(url)
    if err != nil {
        return "",err
    }
	filename:=headerFilename(resp)
	if filename=="untitle.file" {
		return "",nil	
	}
	filepath+=headerFilename(resp)
    defer resp.Body.Close()
	out,err:=os.Create(filepath)
    if err!=nil {
        return "",err
    }
    defer out.Close()
    _, err=io.Copy(out, resp.Body)
    return filepath,err
}

func headerFilename(resp *http.Response) string {
	filename:="untitle.file"
	decho(resp)

	if resp.Header["Content-Length"][0]=="0" {
	} else if resp.Header["Content-Disposition"][0]>"" {
		filename=resp.Header["Content-Disposition"][0]
		filename=strings.Replace(filename,"attachment;", "",-1)
		filename=strings.Replace(filename," ", "",-1)
		filename=strings.Replace(filename,"filename=", "",-1)
		filename=strings.Replace(filename,"'", "",-1)
		//decho(filename)
	} else {
	}
	return filename
}


func decho(v interface{}) {
	if DebugMode {
		echo(v)
	}
}