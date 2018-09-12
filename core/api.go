package core

import (
	"strings"

	"github.com/cubechainofficial/Cubechain/lib"
)

func GetBalance(addr string) int {
	c := CurrentHeight() - 1
	result := CubeBalance(addr, c)
	ci := GetIndexBlock(addr)
	for _, v := range ci {
		a := BlockScan(addr, v)
		vv := a.DataTx
		if lib.ByteToStr(vv.From) == addr {
			result = result - vv.Amount - vv.Fee
		} else if lib.ByteToStr(vv.To) == addr {
			result = result + vv.Amount
		}
	}
	return result
}

func GetTransactionCount(addr string) int {
	c := CurrentHeight() - 1
	ci := GetIndexBlock(addr)
	result := len(ci)
	cdata := CubeScan(addr, c)
	if cdata.DataType != "NULL" {
		result++
	}
	return result
}

func GetTransactionList(addr string) string {
	result := ""
	var res []string
	c := CurrentHeight() - 1
	ci := GetIndexBlock(addr)
	for _, v := range ci {
		a := BlockScan(addr, v)
		vv := a.DataTx
		if lib.ByteToStr(vv.From) == addr || lib.ByteToStr(vv.To) == addr {
			res = append(res, lib.ByteToStr(vv.Hash))
		}
	}
	cdata := CubeScan(addr, c)
	if cdata.DataType != "NULL" {
		res = append(res, lib.ByteToStr(cdata.DataTx.Hash))
	}
	result = strings.Join(res, ",")
	return result
}

func GetTransactionDetail(txhash string) (TransactionData, CubeIndex) {
	var tdata TransactionData
	c := CurrentHeight() - 1
	for i := 1; i <= c; i++ {
		cdata, j := CubeScanHash(txhash, c, "TX")
		if cdata.DataType != "NULL" {
			return cdata.DataTx, CubeIndex{i, j}
		}
	}
	return tdata, CubeIndex{0, 0}
}

func GetTransactionData(txhash string) (TransactionData, CubeIndex) {
	var tdata TransactionData
	c := CurrentHeight() - 1
	for i := 1; i <= c; i++ {
		cdata, j := CubeScanHash(txhash, c, "Data")
		if cdata.DataType != "NULL" {
			return cdata.DataTx, CubeIndex{i, j}
		}
	}
	return tdata, CubeIndex{0, 0}
}

func CubeScan(addr string, idx int) TxData {
	var iTxData TxData
	var ci CubeIndex
	for i := 0; i < 27; i++ {
		ci = CubeIndex{idx, i}
		iTxData = BlockScan(addr, ci)
		if iTxData.DataType != "NULL" {
			return iTxData
		}
	}
	return iTxData
}

func GetIndexBlock(addr string) []CubeIndex {
	var CSindex []CubeIndex
	Cindex := CubeIndex{0, 0}
	var iBlock Block
	c := CurrentHeight() - 1
	err := BlockRead(c, Configure.Indexing, iBlock)
	lib.Err(err, 0)
	iBlockData := IndexDeserialize(iBlock.Data)
	for _, v := range iBlockData.IndexAddress {
		if v.Address == addr {
			return v.Indexing
		} else if v.Address > addr {
			return append(CSindex, Cindex)
		}
	}
	return append(CSindex, Cindex)
}

func BlockScan(addr string, ci CubeIndex) TxData {
	var iBlock Block
	var iTxData TxData
	iTxData.DataType = "NULL"
	err := BlockRead(ci.Index, ci.CubeNum, iBlock)
	lib.Err(err, 0)
	iBlockData := Deserialize(iBlock.Data)
	for _, v := range iBlockData.Tdata {
		vv := v.DataTx
		if lib.ByteToStr(vv.From) == addr || lib.ByteToStr(vv.To) == addr {
			return v
		}
	}
	return iTxData
}

func CubeScanHash(txhash string, idx int, datatype string) (TxData, int) {
	var iTxData TxData
	var ci CubeIndex
	j := 0
	for i := 0; i < 27; i++ {
		ci = CubeIndex{idx, i}
		iTxData, j = BlockScanHash(txhash, ci, datatype)
		if j > 0 {
			return iTxData, j
		}
	}
	return iTxData, 0
}

func BlockScanHash(txhash string, ci CubeIndex, datatype string) (TxData, int) {
	var iBlock Block
	var iTxData TxData

	err := BlockRead(ci.Index, ci.CubeNum, iBlock)
	lib.Err(err, 0)
	iBlockData := Deserialize(iBlock.Data)
	for _, v := range iBlockData.Tdata {
		if datatype == "" || v.DataType == datatype {
			if v.DataType == "TX" {
				vv := v.DataTx
				if lib.ByteToStr(vv.Hash) == txhash {
					return v, ci.CubeNum
				}
			} else if v.DataType == "Data" {
				vv := v.DataTx
				if lib.ByteToStr(vv.Hash) == txhash {
					return v, ci.CubeNum
				}
			}
		}
	}
	return iTxData, 0
}
