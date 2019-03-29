package core

import (
	"strings"
	"strconv"
)

type Cubing struct {
	Cubeno		int
	Timestamp	int
	Hash1		[27]string
	Hash2		[27]string
	PrevHash    string
	CHash       string
	Nonce       int
}

func (cubing *Cubing) FileName() string {
	filename:=cubing.CHash+".cbi"
	return filename
}

func CubingFileName(idx int) string {
	find:=".cbi"	
    dirname:=FilePath(idx)
	result:=FileSearch(dirname,find)
	return result
}

func CubingFileWrite(cubing Cubing) error {
	filename:=cubing.FileName()
	path:=FilePath(cubing.Cubeno)
	err:=FileWrite(path+filepathSeparator+filename,cubing)
	return err
}

func CubingFileRead(index int) Cubing {
	var cubing Cubing
	path:=FilePath(index)
	filename:=FileSearch(path,".cbi")
	err:=FileRead(path+filepathSeparator+filename,&cubing)
	if err!=nil {
		decho(err)
	}
	return cubing
}

func CubingFileRead2(index int,hash string) Cubing {
	var cubing Cubing
	filename:=hash+".cbi"
	path:=FilePath(index)
	err:=FileRead(path+filepathSeparator+filename,&cubing)
	if err!=nil {
		decho(err)
	}
	return cubing
}

func GetCubing(cubeno int) Cubing {
	var gcubing Cubing
	var hash3 [27]string
	var hash4 [27]string
	
	r:=NodeCube("getcubing","0&cubeno="+strconv.Itoa(cubeno))
	
	if r>"" {
		result:=strings.Split(r, "|")
		hash1:=strings.Split(result[1], ",")
		hash2:=strings.Split(result[2], ",")


		gcubing.Cubeno=cubeno
		gcubing.Timestamp,_=strconv.Atoi(result[0])
		

		for i:=0;i<27;i++ {
			hash3[i]=hash1[i]
			hash4[i]=hash2[i]
		}
		gcubing.Hash1=hash3
		gcubing.Hash2=hash4

		gcubing.PrevHash=result[3]
		gcubing.CHash=result[4]
		gcubing.Nonce,_=strconv.Atoi(result[5])
	}
	return gcubing
}


func CubenoSet()  {
	k:=0
	for i:=0;i<27;i++ {
		switch i+1 {
		case Configure.Indexing:
			CubeSetNum[i]="Indexing";
		case  Configure.Statistics:
			CubeSetNum[i]="Statistics";
		case  Configure.Escrow:
			CubeSetNum[i]="Escrow";
		case  Configure.Format:
			CubeSetNum[i]="Format";
		case  Configure.Edit:
			CubeSetNum[i]="Edit";
		default:
			k++
			CubeSetNum[i]="Data"+strconv.Itoa(k);
		}
	}
}
