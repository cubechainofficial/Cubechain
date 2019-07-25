package core

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
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

func setHash5(str string) string {
	s:=setHash(str)
	result:=s[:20]+s[43:]+s[20:43]
	return result
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

func PowSpecialHashing(pstr string,max int) (string,int) {
	var str,hash string
	i:=0
	for i=0;i<=max;i++ {
		str=pstr+strconv.Itoa(i)
		hash=setHash5(str)
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
	if len(MineDifficulty)<8 {
		MineDifficulty=MineDifficultyBase
	}
	if hash[:len(MineDifficulty)]<=MineDifficulty {
		result=true
	}
	return result
}

func PatternHash(str string,cno int) string {
	result:=""
	ph:=""
	switch cno {
		case 3 : 
			ph=setHash3(str)
			result=ph[4:27]+ph[:4]+ph[27:]
		case 4 : 
			ph=setHash3(str)
			result=ph[:14]+ph[14:]
		case 5 : 
			ph=setHash3(str)
			result=ph[7:31]+ph[31:]+ph[:7]
		case 6 : 
			ph=setHash3(str)
			result=ph[8:30]+ph[51:]+ph[:8]+ph[30:51]
		default : 
			return setHash3(str)
	}
	return result
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
		case 5 : return setHash5(str)
		default : return setHash(str)
	}
}

