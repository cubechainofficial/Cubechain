package core

import (
	"strconv"
	"strings"

	"github.com/cubechainofficial/Cubechain/config"
	"github.com/cubechainofficial/Cubechain/lib"
)

var Pratio = Pohr{25, 25, 25, 25}

func preHashVerify(hash string) bool {
	result := false
	vhash := strings.Repeat(hash[:1], Configure.Zeronumber)
	if hash[:Configure.Zeronumber] == vhash {
		result = true
	}
	return result
}

func PowBlockHashing(pstr string, max int) (string, int) {
	var str, hash string
	i := 0
	for i = 0; i <= max; i++ {
		str = pstr + strconv.Itoa(i)
		hash = setHash(str)
		if preHashVerify(hash) {
			break
		}
	}
	return hash, i
}

func PowBlockVerifyNonce(hash string, pstr string, nonce int) bool {
	result := false
	var str, vhash string
	str = pstr + strconv.Itoa(nonce)
	vhash = setHash(str)
	if vhash == hash {
		result = true
	}
	return result
}

func PowBlockVerify(hash string, pstr string, nonce int) bool {
	result := false
	vhash := strings.Repeat(hash[:1], Configure.Zeronumber)
	if hash[:Configure.Zeronumber] == vhash {
		result = PowBlockVerifyNonce(hash, pstr, nonce)
	}
	return result
}

func PowBlockHashing2(pstr string, max int) (string, int) {
	var str string
	var hash string
	i := 0
	for i = 0; i <= max; i++ {
		str = pstr + strconv.Itoa(i)
		hash = setHash(str)
		if preHashVerify(hash) {
			break
		}
	}
	return hash, i
}

func PowCubingHashing(pstr string, max int) (string, int) {
	var str, hash string
	i := 0
	for i = 0; i <= max; i++ {
		str = pstr + strconv.Itoa(i)
		hash = setHash2(str)
		if preHashVerify(hash) {
			break
		}
	}
	return hash, i
}

func PowCubingHash(pstr string) string {
	hash, _ := PowCubingHashing(pstr, Configure.Maxnonce)
	return hash
}

func (p POH) PowCubeHash(b *Block, addr string) string {
	p = POH{b.Index, b.Cubeno, b.Hash, 1, addr, Pratio.CHash, 0}
	var str, pstr, hash string
	l := len(config.CubeCo[b.Cubeno])
	for _, v := range config.CubeCo[b.Cubeno] {
		pstr += ReadCubeHash(b.Index-1, v-1)
	}
	for i := 0; i <= Configure.Maxnonce; i++ {
		str = pstr + strconv.Itoa(i)
		hash = setCubeHash(str, l)
		if preHashVerify(hash) {
			break
		}
	}
	return hash
}

func (p POH) PowBlockHash(b Block, addr string) (string, Block) {
	p = POH{b.Index, b.Cubeno, b.Hash, 2, addr, Pratio.BlockHash, 0}
	var str string
	var hash string
	i := 0
	p.State = 0
	pstr := strconv.Itoa(b.Index) + strconv.Itoa(b.Cubeno) + strconv.Itoa(b.Timestamp) + lib.ByteToStr(b.Data) + b.PrevHash + b.PrevCubeHash
	for i = 0; i <= Configure.Maxnonce; i++ {
		str = pstr + strconv.Itoa(i)
		hash = setHash(str)
		if preHashVerify(hash) {
			p.State = 1
			break
		}
	}
	b.Nonce = i
	b.Hash = hash
	return hash, b
}

func (p POH) PowCubing(c *CubeBlock, addr string) (string, *CubeBlock) {
	p = POH{c.Index, 2, c.Chash, 3, addr, Pratio.Cubing, 0}
	var bstr []byte
	var str string
	var hash string
	pstr := strconv.Itoa(c.Index) + strconv.Itoa(c.Timestamp) + lib.ByteToStr(bstr)
	p.State = 0
	for i := 0; i <= 100; i++ {
		for _, v := range c.Cube {
			bstr = Serialize(v)
		}
		str = pstr + strconv.Itoa(i)
		hash = setHash2(str)
		if preHashVerify(hash) {
			break
			p.State = 1
		}
	}
	c.Chash = hash
	return hash, c
}

func (p POH) PosRun(b Block, addr string) (bool, []PosWallet, int) {
	p = POH{b.Index, b.Cubeno, b.Hash, 3, addr, Pratio.POS, 0}
	var posw []PosWallet
	var total int
	if b.Cubeno != Configure.Statistics {
		return false, posw, 0
	} else {
		var iBlock Block
		var pAdd SBlock
		err := BlockRead(b.Index, b.Cubeno, &iBlock)
		lib.Err(err, 0)
		pAdd = StaticDeserialize(iBlock.Data)
		for _, v := range pAdd.RuleArr {
			posw = append(posw, PosWallet{v.Address, v.Balance})
			total += v.Balance
		}
	}
	if total > 0 {
		return true, posw, total
	} else {
		return false, posw, 0
	}
}

func ReadCubeHash(index int, cubeno int) string {
	var iBlock Block
	var hash string
	err := BlockRead(index, cubeno, &iBlock)
	lib.Err(err, 0)
	hash = iBlock.Hash
	return hash
}
