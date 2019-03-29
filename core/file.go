package core

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func FilePath(idx int) string {
	divn:=idx/Configure.Datanumber
	divm:=idx%Configure.Datanumber
	if divm>0 {
		divn++
	} else if divm==0 {
		divm=Configure.Datanumber
	}
	if divn==0 {
		divn++
		divm=1
	}
	nhex:=fmt.Sprintf("%x",Configure.Datanumber)
	mcnt:=len(nhex)
	nstr:=fmt.Sprintf("%0.5x",divn)
	mstr:=fmt.Sprintf("%0."+strconv.Itoa(mcnt)+"x",divm)
	dirname:=Configure.Datafolder+filepathSeparator+nstr+filepathSeparator+mstr
	if DirExist(dirname)==false {
		if err:=os.MkdirAll(dirname, os.FileMode(0755)); err!=nil {
			return "Directory not found.\\1\\1"
		}	
	}	
	return dirname
}

func FileWrite(path string, object interface{}) error {
	file,err:=os.Create(path)
	if err==nil {
		encoder:=gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func FileRead(path string, object interface{}) error {
	/*
	var block Block
	var cube Cube
	var cubing Cubing
	
	gob.Register(block)  
	gob.Register(cube)
	gob.Register(cubing)  
	*/

	file,err:=os.Open(path)
	if err==nil {
		decoder:=gob.NewDecoder(file)
		err=decoder.Decode(object)
	}
	file.Close()
	return err
}

func FileSize(dirpath string) int64 {
	file, err := os.Open(dirpath) 
	if err != nil {
		echo (err)
	}
	fi, err := file.Stat()
	if err != nil {
		echo (err)
	}
	return fi.Size()
}

func FileSearch(dirname string,find string) string{
    result:=""
	d,err:=os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    file, err:=d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi:=range file {
        if fi.Mode().IsRegular() {
			fstr:=fi.Name()
			if strings.Index(fstr,find)>=0 {
				result=fi.Name()
				return result
			}
        }
    }
	return result
}

func DirExist(dirName string) bool{
	result:=true
	_,err:=os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			result=false
		}
	}
	return result
}

func MaxFind(dirpath string) string {
	find:="0"
    d, err:=os.Open(dirpath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
	fi, err:=d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi:=range fi {
        if fi.Mode().IsRegular() {
        } else {
  			if fi.Name()>find {
				find=fi.Name()
			}
		}
   }
   return find
}

func CubePathNum(path string) int {
	result:=0
	split:=strings.Split(path, filepathSeparator)
	slen:=len(split)
	nint,_:=strconv.ParseUint(split[slen-2],16,32)
	mint,_:=strconv.ParseUint(split[slen-1],16,32)
	result=(int(nint)-1)*Configure.Datanumber+int(mint)
	return result
}

func PathDelete(path string) error {
	err:=os.RemoveAll(path)
	os.MkdirAll(path,0755)
	return err
}



