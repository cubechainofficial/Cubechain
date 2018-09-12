package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cubechainofficial/Cubechain/lib"
)

var CubeSetNum [27]string

func CubenoSet() {
	k := 0
	for i := 0; i < 27; i++ {
		switch i + 1 {
		case Configure.Indexing:
			CubeSetNum[i] = "Indexing"
		case Configure.Statistics:
			CubeSetNum[i] = "Statistics"
		case Configure.Escrow:
			CubeSetNum[i] = "Escrow"
		case Configure.Format:
			CubeSetNum[i] = "Format"
		case Configure.Edit:
			CubeSetNum[i] = "Edit"
		default:
			k++
			CubeSetNum[i] = "Data" + strconv.Itoa(k)
		}
	}
}

func ChainSet(cubechain *Cubechain) {
	*cubechain = ChainFileRead()
}

func CubingToChain(idx int, cubechain *Cubechain) int {
	result := cubechain.Verify
	cubing := CubingFileRead2(idx)
	cubechain.Chain = append(cubechain.Chain, CChain{cubing.Index, cubing.Timestamp, cubing.Chash})
	return result
}

func ChainCheck(cubechain *Cubechain) int {
	result := 0
	ch := CurrentHeight()
	ch--
	if ch == cubechain.Verify {
		if len(cubechain.Chain) == cubechain.Verify {
			fmt.Println("Chain is verify.")
		} else {
			fmt.Println("Chain is not verify.")
		}
	} else {
		for ch > cubechain.Verify {
			cubechain.Verify++
			result = CubingToChain(cubechain.Verify, cubechain)
		}
	}
	ChainFileWrite(cubechain)
	return result
}

func chainFileName() string {
	filename := setHash(Configure.Network) + ".chn"
	return filename
}

func ChainFileWrite(cubechain *Cubechain) error {
	filename := chainFileName()
	path := Configure.Datafolder
	err := fileWrite2(path+string(filepath.Separator)+filename, cubechain)
	return err
}

func ChainFileRead() Cubechain {
	var cubechain Cubechain
	filename := chainFileName()
	path := Configure.Datafolder
	err := pathRead(path+string(filepath.Separator)+filename, &cubechain)
	if err != nil {
		fmt.Println(err)
	}
	return cubechain
}

func CubingSet(cblock *CubeBlock) Cubing {
	var cubing Cubing
	cubing.Index = cblock.Index
	cubing.Timestamp = cblock.Timestamp
	cubing.Chash = cblock.Chash
	for i := 0; i < 27; i++ {
		cubing.Hash1[i] = cblock.Cube[i].Hash
		cubing.Hash2[i] = cblock.Cube[i].PrevCubeHash
	}
	return cubing
}

func cubingFileName(cubing *Cubing) string {
	filename := cubing.Chash + ".cbi"
	return filename
}

func CubingFileWrite(cubing Cubing) error {
	filename := cubingFileName(&cubing)
	path := CubePath(cubing.Index)
	err := fileWrite2(path+string(filepath.Separator)+filename, cubing)
	return err
}

func CubingFileRead(index int, hash string) Cubing {
	var cubing Cubing
	filename := hash + ".cbi"
	path := CubePath(index)
	err := pathRead(path+string(filepath.Separator)+filename, &cubing)
	if err != nil {
		fmt.Println(err)
	}
	return cubing
}

func CubingFileRead2(index int) Cubing {
	var cubing Cubing
	path := CubePath(index)
	filename := fileSearch(path, ".cbi")
	err := pathRead(path+string(filepath.Separator)+filename, &cubing)
	if err != nil {
		fmt.Println(err)
	}
	return cubing
}

func cfileName(cblock *CubeBlock) string {
	filename := cblock.Chash + ".cub"
	return filename
}

func cblockFile(cblock *CubeBlock) error {
	filename := cfileName(cblock)
	path := CubePath(cblock.Index)
	err := fileWrite2(path+string(filepath.Separator)+filename, cblock)
	if err == nil {
		cubing := CubingSet(cblock)
		CubingFileWrite(cubing)
	}
	return err
}

func CubePath(idx int) string {
	divn := idx / Configure.Datanumber
	divm := idx % Configure.Datanumber
	if divm > 0 {
		divn++
	} else if divm == 0 {
		divm = Configure.Datanumber
	}
	if divn == 0 {
		divn++
		divm = 1
	}
	nhex := fmt.Sprintf("%x", Configure.Datanumber)
	mcnt := len(nhex)
	nstr := fmt.Sprintf("%0.5x", divn)
	mstr := fmt.Sprintf("%0."+strconv.Itoa(mcnt)+"x", divm)
	dirname := Configure.Datafolder + string(filepath.Separator) + nstr + string(filepath.Separator) + mstr
	if dirExist(dirname) == false {
		if err := os.MkdirAll(dirname, os.FileMode(0755)); err != nil {
			return "Directory not found.\\1\\1"
		}
	}
	return dirname
}

func CubePathNum(path string) int {
	result := 0
	separator := string(filepath.Separator)
	split := strings.Split(path, separator)
	slen := len(split)
	nint, _ := strconv.ParseUint(split[slen-2], 16, 32)
	mint, _ := strconv.ParseUint(split[slen-1], 16, 32)
	result = (int(nint)-1)*Configure.Datanumber + int(mint)
	return result
}

func GenesisCube() CubeBlock {
	var clen = CurrentHeight()
	var block Block
	var sblock [27]Block
	var cblock CubeBlock
	var chash string
	if clen > 1 {
		fmt.Printf("Invalid genesis cube.\n")
		return cblock
	}
	for i := 0; i < 27; i++ {
		block = GenesisBlock(i)
		sblock[i] = block
		chash += block.Hash
	}
	cblock = CubeBlock{1, 0, sblock, PowCubingHash(chash)}
	cblock.Timestamp = int(time.Now().Unix())
	cblockFile(&cblock)
	return cblock
}

func AddCube(str string) CubeBlock {
	var clen = CurrentHeight()
	var block Block
	var sblock [27]Block
	var cblock CubeBlock
	var chash string
	if clen < 1 {
		clen = 1
	}
	for i := 0; i < 27; i++ {
		str = str + strconv.Itoa(i)
		block = addBlock([]byte(str), i)
		sblock[i] = block
		chash += block.Hash
	}
	cblock = CubeBlock{clen, 0, sblock, PowCubingHash(chash)}
	cblock.Timestamp = int(time.Now().Unix())
	cblockFile(&cblock)
	return cblock
}

func GetBalanceCheck(addr string) int {
	var amount = 0
	var iBlock [27]Block
	var Txd TransactionData
	c := CurrentHeight() - 1
	for i := 0; i < 27; i++ {
		err := BlockRead(c, i, iBlock[i])
		lib.Err(err, 0)
		if i == Configure.Indexing || i == Configure.Statistics || i == Configure.Escrow || i == Configure.Format || i == Configure.Edit {
		} else {
			iData := Deserialize(iBlock[i].Data)
			for _, v := range iData.Tdata {
				if v.DataType == "tx" {
					if lib.ByteToStr(Txd.From) == addr {
						amount += Txd.Amount * (-1)
					}
					if lib.ByteToStr(Txd.To) == addr {
						amount += Txd.Amount
					}
				}
			}
		}
	}
	return amount
}

func CubeBalance(addr string, c int) int {
	var amount = 0
	var iBlock [27]Block
	var Txd TransactionData
	if c <= 0 {
		c = CurrentHeight() - 1
	}
	for i := 0; i < 27; i++ {
		err := BlockRead(c, i, &iBlock[i])
		lib.Err(err, 0)
		if i == Configure.Indexing || i == Configure.Format || i == Configure.Edit {
		} else if i == Configure.Statistics {
		} else if i == Configure.Escrow {
		} else {
			iData := Deserialize(iBlock[i].Data)
			for _, v := range iData.Tdata {
				if v.DataType == "tx" {
					if lib.ByteToStr(Txd.From) == addr {
						amount += Txd.Amount * (-1)
					}
					if lib.ByteToStr(Txd.To) == addr {
						amount += Txd.Amount
					}
				}
			}
		}
	}
	return amount
}

func CubeCount(addr string, c int) (int, int) {
	var count = 0
	var ecount = 0
	var iBlock [27]Block
	if c <= 0 {
		c = CurrentHeight() - 1
	}
	for i := 0; i < 27; i++ {
		err := BlockRead(c, i, &iBlock[i])
		lib.Err(err, 0)
		if i == Configure.Indexing || i == Configure.Format || i == Configure.Edit {
		} else if i == Configure.Statistics {
		} else if i == Configure.Escrow {
			ecount++
		} else {
			count++
		}
	}
	return count, ecount
}
