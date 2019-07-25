package core

import (
	"time"
	"strings"
	"strconv"
)

type Pohr struct {
	BlockHash	float64
	CHash		float64
	Cubing		float64
	POW			float64
	POS			float64
}

type PoolMining struct {
	Timestamp	int
	Cubeno		int
	Blockno		int
	HashStr		string
	ResultHash	string
	ResultNonce	int
	Reward		float64
	Addr		string
	Success		bool
}

func PohSet(cubeno int) PoolMining {
	var pm PoolMining
	pm.Timestamp=int(time.Now().Unix())
	pm.Addr=Configure.Address
	pm.Success=false

	return pm
}

func (pm *PoolMining) CubeHeight() {
	pm.Cubeno,_=strconv.Atoi(GetCubeHeight())
	pm.Cubeno++
}

func (pm *PoolMining) HashString() {
	result:=NodeSend("hashstring","=0&cubeno="+strconv.Itoa(pm.Cubeno)+"&blockno="+strconv.Itoa(pm.Blockno))
	pm.HashStr=result
}

func (pm *PoolMining) Result(rcnt int) {
	bstr:=NodeSend2("pool_result","=0&cubeno="+strconv.Itoa(pm.Cubeno)+"&blockno="+strconv.Itoa(pm.Blockno)+"&hashstr="+pm.HashStr)
	result:=strings.Split(bstr,"|")
	pm.ResultHash=result[0]
	pm.ResultNonce,_=strconv.Atoi(result[1])
	if pm.ResultHash>"" && pm.ResultNonce>0 {
		h:=CallHash(pm.HashStr+result[1],1)
		if pm.ResultHash!=h || PHashVerify(h)==false {
			rf:=NodeSend2("pool_fault","0&cubeno="+strconv.Itoa(pm.Cubeno)+"&blockno="+strconv.Itoa(pm.Blockno)+"&resulthash="+pm.ResultHash+"&resultnonce="+strconv.Itoa(pm.ResultNonce))
			pm.ResultHash=""
			pm.ResultNonce=0
			echo (rf)
			echo ("Mining Fault")
		} else {
			echo ("Mining Success")
			return
		}
	} else {
	}
	rcnt++
	if rcnt<11 {
		time.Sleep(2*time.Second)
		pm.Result(rcnt)
	}
}

func PowTx() TxData {
	var txd TxData
	txd.Datatype="POW"
	txd.Timestamp=int(time.Now().Unix())
	txd.From="0"
	txd.To=Configure.Address
	txd.Amount=Pratio.BlockHash
	txd.Fee=0.0
	txd.Nonce=0
	txd.Hash=setHash(strconv.Itoa(txd.Timestamp)+txd.From+txd.To+strconv.FormatFloat(txd.Amount,'f',-1,64)+strconv.Itoa(txd.Nonce))
	txd.Sign=setHash(txd.Hash)
	if Configure.MiningMode>"" {
		txd.Message=Configure.MiningMode
	} else {
		txd.Message=""
	}
	return txd
}


