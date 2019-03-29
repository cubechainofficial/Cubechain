package core

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"strconv"
)

func setHashV(str interface{}) string {
	h:=sha256.New()
	h.Write(GetBytes(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash(str string) string {
	h:=sha256.New()
	h.Write([]byte(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash2(str string) string {
	h:=sha512.New384()
	h.Write([]byte(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash3(str string) string {
	h:=sha512.New512_224()
	h.Write([]byte(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash4(str string) string {
	h:=sha512.New()
	h.Write([]byte(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func PowBlockHashing(pstr string,max int) (string,int) {
	var str,hash string
	i:=0
	for i=0;i<=max;i++ {
		str=pstr+strconv.Itoa(i)
		hash=setHash(str)
		if PHashVerify(hash) {
			break
		}
	}
	if i>max { i=max; }
	return hash,i
}

func PowCubeHashing(pstr string,max int) (string,int) {
	var str,hash string
	i:=0
	for i=0;i<=max;i++ {
		str=pstr+strconv.Itoa(i)
		hash=setHash2(str)
		if PHashVerify(hash) {
			break
		}
	}
	if i>max { i=max; }
	return hash,i
}

func PHashVerify(hash string) bool {
	result:=false
	vhash:=strings.Repeat(hash[:1], Configure.Zeronumber)
	if hash[:Configure.Zeronumber]==vhash {
		result=true
	}
	return result
}

func PatternHash(str string,cno int) string {
	switch cno {
		case 3 : return setHash3(str)
		case 4 : return setHash3(str)
		case 5 : return setHash3(str)
		case 6 : return setHash3(str)
		default : return setHash3(str)
	}
}

func BlockHash(str string) string {
	return setHash(str)
}

func CubingHash(str string) string {
	return setHash2(str)
}

func CallHash(str string,hno int) string {
	switch hno {
		case 2 : return setHash2(str)
		case 3 : return setHash3(str)
		case 4 : return setHash4(str)
		default : return setHash(str)
	}
}

