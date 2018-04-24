package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/gob"
	"runtime"
	"fmt"
	"os"
	"strconv"
	"strings"
    "path/filepath"
)

func vaildCubeno(cno int) bool {
	result:=true
	if cno<1 || cno>27 {
		result=false	
	} 
	return result
}

func setHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func sethash2(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}


func ArgServer() {
    fmt.Println(os.Args)
}

func setCubeHash(str string, cno int) string {
	switch cno {
		case 3 : return setHash(str)
		case 4 : return setHash(str)
		case 5 : return setHash(str)
		case 6 : return setHash(str)
		default : return setHash(str)
	}
	return setHash(str)
}

func fileWrite(path string, object interface{}) error {
	datapath:="bdata"+ string(filepath.Separator) 
	file, err := os.Create(datapath+path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
 }

func fileWrite2(path string, object interface{}) error {
	datapath:="cdata"+ string(filepath.Separator) 
	file, err := os.Create(datapath+path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func fileRead(filename string, object interface{}) error {
	datapath:="bdata"+ string(filepath.Separator) 
	file, err := os.Open(datapath+filename)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func fileCheck(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}

func blockName(find string) string {
	result:=""
    dirname := "." + string(filepath.Separator) + "bdata"
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
            fstr:=fi.Name()
			fstr=fstr[0:len(find)]
			if fstr==find {
				result=fi.Name()
				return result
			}
        }
    }
	return result
}

func blockFinder(find string) bool {
	result:=false
    dirname := "." + string(filepath.Separator) + "bdata"
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
            fstr:=fi.Name()
			fstr=fstr[0:len(find)]
			if fstr==find {
				result=true
				return result
			}
        }
    }
	return result
}

func CurrentHeight() int {
	result:=0
    dirname := "." + string(filepath.Separator) + "bdata"
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	max:=0
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
			fstr:=strings.Split(fi.Name(),"_")
			fmt.Println(fi.Name(), fi.Size(), "bytes", fstr[0] ," ", fstr[1]," ", fstr[2])
			cstr:=fstr[0]
			cno,_:=strconv.Atoi(cstr)
			if cno>max {
				max=cno
			}
			result=max
        }
    }
	result++
	return result	
}

func blockRead(index int,cubeno int,object interface{}) error {
	err:=fileRead(blockName(strconv.Itoa(index)+"_"+strconv.Itoa(cubeno)+"_"),object)
	return err;
}

