package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cubechainofficial/Cubechain/lib"
	"github.com/cubechainofficial/Cubechain/wallet"
)

var echo = fmt.Println

func GenesisBlock(Cubenum int) Block {
	var block = *new(Block)
	if !vaildCubeno(Cubenum) {
		return block
	}
	bFind := blockFinder(1, Cubenum)
	if bFind {
		return block
	}
	var txp = new(TxPool)
	var txd = new(TxData)
	txstr := "-1"
	amount := Cubenum * 1000
	w := *wallet.CreateWallet()
	a := w.GetAddress()
	addr := lib.ByteToStr(a)
	tx := txd.TxInput(w, txstr, addr, amount)
	txp.Tdata = append(txp.Tdata, tx)
	block.Index = 1
	block.Cubeno = Cubenum
	block.Timestamp = int(time.Now().Unix())
	block.Nonce = 0
	block.Data = Serialize(*txp)
	block.PrevCubeHash = cubeHash(&block)
	block.PrevHash = prvHash(Cubenum)
	block = calculateHash(block)
	return block
}

func addBlock(data []byte, Cubenum int) Block {
	var block Block
	var oblock Block
	ch := CurrentHeight()
	if ch < 1 {
		echo("Genesis block not create")
		return block
	}
	err := BlockRead(ch-1, Cubenum, &oblock)
	lib.Err(err, 0)
	block.Index = ch
	block.Timestamp = int(time.Now().Unix())
	block.Nonce = 0
	block.Cubeno = Cubenum
	block.Data = data
	block.PrevHash = oblock.Hash
	block.PrevCubeHash = cubeHash(&block)
	block = calculateHash(block)
	return block
}

func StringBlock(block Block) string {
	str := strconv.Itoa(block.Index) + "|" + strconv.Itoa(block.Cubeno) + "|" + strconv.Itoa(block.Timestamp) + "|" + TxpoolStr(Deserialize(block.Data)) + "|" + block.Hash + "|" + block.PrevHash + "|" + block.PrevCubeHash + "|" + strconv.Itoa(block.Nonce)
	return str
}

func calculateHash(block Block) Block {
	var p = new(POH)
	hash, blc := p.PowBlockHash(block, Configure.Address)
	if blc.Hash == hash {
		blockFile(blc)
	}
	return blc
}

func calculateHashVerify(block Block) string {
	str := strconv.Itoa(block.Index) + strconv.Itoa(block.Cubeno) + strconv.Itoa(block.Timestamp) + lib.ByteToStr(block.Data) + block.PrevHash + block.PrevCubeHash + strconv.Itoa(block.Nonce)
	return setHash(str)
}

func cubeHash(block *Block) string {
	var p = new(POH)
	hash := p.PowCubeHash(block, Configure.Address)
	return hash
}

func cubeHashVerify(cubenum int) string {
	return setHash2(strconv.Itoa(cubenum))
}

func prvHash(cubenum int) string {
	return setHash(strconv.Itoa(cubenum))
}

func fileName(block *Block) string {
	filename := strconv.Itoa(block.Index) + "_" + strconv.Itoa(block.Cubeno) + "_" + block.Hash + ".blc"
	return filename
}

func blockFile(block Block) error {
	filename := fileName(&block)
	datapath := CubePath(block.Index) + string(filepath.Separator)
	err := fileWrite2(datapath+string(filepath.Separator)+filename, block)
	return err
}

func blockStrFile(block Block) error {
	filename := fileName(&block)
	str := strconv.Itoa(block.Index) + strconv.Itoa(block.Cubeno) + strconv.Itoa(block.Timestamp) + lib.ByteToStr(block.Data) + block.PrevHash + block.PrevCubeHash + strconv.Itoa(block.Nonce)
	bytes := []byte(str)
	err := ioutil.WriteFile(CubePath(block.Index)+string(filepath.Separator)+filename, bytes, 0)
	return err
}

func StrFile(str string, idx int) error {
	filename := "str" + strconv.Itoa(idx) + ".blc"
	bytes := []byte(str)
	err := ioutil.WriteFile(CubePath(idx)+string(filepath.Separator)+filename, bytes, 0)
	return err
}

func PrintBlockHead(block Block) {
	echo("Index =", block.Index)
	echo("Cubeno =", block.Cubeno)
	echo("Timestamp =", block.Timestamp)
	echo("Hash =", block.Hash)
	echo("PrevHash =", block.PrevHash)
	echo("PrevCubeHash =", block.PrevCubeHash)
	echo("Nonce =", block.Nonce)
}
