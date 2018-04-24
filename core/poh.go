package core

import (
	"strconv"
	"../lib"
	"../config"
)

var Pratio=Pohr{50,30,20,30}

func (p POH) PowCubeHash(b Block,addr string) string {
	p=POH{b.Index,b.Cubeno,b.Hash,1,addr,Pratio.CHash,0}
	var str string
	var hash string
	l:=len(config.CubeCo[b.Cubeno])
	for _,v := range config.CubeCo[b.Cubeno] {
		str+=ReadCubeHash(b.Index-1,v)
	}
	for i:=0; i<=b.Nonce; i++ {
		hash=setCubeHash(str,l)
		if hash[:1]=="0" {
			break
		}
	}
	return hash
}

func (p POH) PowBlockHash(b Block,addr string) (string,Block) {
	p=POH{b.Index,b.Cubeno,b.Hash,2,addr,Pratio.BlockHash,0}
	var str string
	var hash string
	i:=0
	for i=0; i<=b.Nonce; i++ {
		str=strconv.Itoa(b.Index) + strconv.Itoa(b.Cubeno) + strconv.Itoa(b.Timestamp) + lib.ByteToStr(b.Data) + b.PrevHash + b.PrevCubeHash + strconv.Itoa(i)
		hash=setHash(str)
		if hash[:1]=="0" {
			break
		}
	}
	p.State=1
	b.Nonce=i
	b.Hash=hash
	return hash,b
}

func (p POH) PowCubing(c *CubeBlock,addr string) (string,*CubeBlock)  {
	p=POH{c.Index,2,c.Chash,3,addr,Pratio.Cubing,0}
	var bstr []byte
	var str string
	var hash string
	for i:=0; i<=100; i++ {
		for _,v := range c.Cube {
			bstr=Serialize(v)
		}
		str=strconv.Itoa(c.Index) + strconv.Itoa(c.Timestamp) + lib.ByteToStr(bstr) + strconv.Itoa(i)
		hash=sethash2(str)
		if hash[:1]=="0" {
			break
		}
	}
	p.State=1
	c.Chash=hash
	return hash,c

}

func (p POH) PosRun(b Block,addr string) (bool,[]PosWallet,int) {
	p=POH{b.Index,b.Cubeno,b.Hash,3,addr,Pratio.POS,0}
	var posw []PosWallet
	var total int
	if b.Cubeno!=Configure.Statistics {
		return false,posw,0
	} else {
		var iBlock Block
		var pAdd SBlock
		err:=blockRead(b.Index,b.Cubeno,iBlock)
		lib.Err(err,0)	
		pAdd=StaticDeserialize(iBlock.Data)
		for _,v := range pAdd.RuleArr {
			posw=append(posw,PosWallet{v.Address,v.Balance})
			total+=v.Balance
		}
	}
	if total>0 {
		return true,posw,total
	} else {
		return false,posw,0
	}
}

func ReadCubeHash(index int, cubeno int) string {
	var iBlock Block
	var hash string
	err:=blockRead(index,cubeno,iBlock)
	lib.Err(err,0)
	hash=iBlock.Hash
	return hash
}