package core

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func validCubeno(cno int) bool {
	result := true
	if cno < 0 || cno > 26 {
		result = false
	}
	return result
}

func setHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash2(str string) string {
	h := sha512.New384()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash3(str string) string {
	h := sha512.New512_224()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash4(str string) string {
	h := sha512.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CallHash(str string, hno int) string {
	switch hno {
	case 2:
		return setHash2(str)
	case 3:
		return setHash3(str)
	case 4:
		return setHash4(str)
	default:
		return setHash(str)
	}
}

func setCubeHash(str string, cno int) string {
	switch cno {
	case 3:
		return setHash3(str)
	case 4:
		return setHash3(str)
	case 5:
		return setHash3(str)
	case 6:
		return setHash3(str)
	default:
		return setHash3(str)
	}
}

func fileWrite(path string, object interface{}) error {
	idx := object.(Block).Index
	datapath := CubePath(idx) + string(filepath.Separator)
	file, err := os.Create(datapath + path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(&object)
	}
	file.Close()
	return err
}

func fileWrite2(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func fileRead(filename string, object interface{}) error {
	idx := object.(Block).Index
	datapath := CubePath(idx) + string(filepath.Separator)
	file, err := os.Open(datapath + filename)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&object)
	}
	file.Close()
	return err
}

func FileRead(filename string, object interface{}) error {
	idx := object.(Block).Index
	datapath := CubePath(idx) + string(filepath.Separator)
	file, err := os.Open(datapath + filename)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func pathRead(path string, object interface{}) error {
	file, err := os.Open(path)
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

func BlockName(idx int, cno int) string {
	find := strconv.Itoa(idx) + "_" + strconv.Itoa(cno-1) + "_"
	dirname := CubePath(idx)
	result := fileSearch(dirname, find)
	return result
}

func blockFinder(idx int, cno int) bool {
	result := false
	bf := BlockName(idx, cno)
	if bf > "" {
		result = true
	}
	return result
}

func fileSearch(dirname string, find string) string {
	result := ""
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
			fstr := fi.Name()
			if strings.Index(fstr, find) >= 0 {
				result = fi.Name()
				return result
			}
		}
	}
	return result
}

func MaxFind(dirpath string) string {
	find := "0"
	d, err := os.Open(dirpath)
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
		} else {
			if fi.Name() > find {
				find = fi.Name()
			}
		}
	}
	return find
}

func BlockRead(index int, cubeno int, object interface{}) error {
	filename := BlockName(index, cubeno)
	datapath := CubePath(index) + string(filepath.Separator)
	file, err := os.Open(datapath + filename)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func CurrentHeight() int {
	result := 0
	f := MaxFind(Configure.Datafolder + string(filepath.Separator))
	if f == "0" {
		return 1
	}
	f2 := MaxFind(Configure.Datafolder + string(filepath.Separator) + f)
	if f2 == "0" {
		return 1
	}
	nint, _ := strconv.ParseUint(f, 16, 32)
	mint, _ := strconv.ParseUint(f2, 16, 32)
	result = (int(nint)-1)*Configure.Datanumber + int(mint)
	if fileSearch(CubePath(result), ".cub") > "" {
		result++
	}
	return result
}

func dirExist(dirName string) bool {
	result := true
	_, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			result = false
		}
	}
	return result
}

func PathDelete(path string) error {
	err := os.RemoveAll(path)
	os.MkdirAll(path, 0755)
	return err
}
