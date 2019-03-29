package core

import (
)






/*
func ReadCubeHash(index int) string {
	var cBlock CubeBlock
	var hash string
	err:=CubeRead(index,&cBlock)
	Err(err,0)
	hash=cBlock.Chash
	return hash
}

func ReadCubeVerify(index int) (bool) {
	result:=false
	var cBlock CubeBlock
	hstr:=""
	err:=CubeRead(index,&cBlock)
	Err(err,0)
	for i:=0;i<27;i++ {
		hstr+=cBlock.Cube[i].Hash
	}
	if cBlock.Chash==PowCubingHash(hstr) {
		result=true
	}
	return result
}


func ReadBlockVerify(index int,cubeno int) (bool) {
	var txData TxData
	var iBlock Block
	err:=BlockRead(index,cubeno,&iBlock)
	Err(err,0)
	var p=new(POH)
	result:=PowBlockVerify(iBlock.Hash,BlockVerifyString(iBlock),iBlock.Nonce)
	if result {
		result=PowBlockVerifyChain(iBlock.PrevHash,index,cubeno)
	}
	if result {
		poh:=BlockTreeRoot(iBlock.Data,"Poh")
		if poh.Left!=nil {
			txData=TdDeserialize(GetBytes(poh.Left.Val))
			if txData.DataType=="POW" {
				if txData.DataTx.Amount-txData.DataTx.Fee!=Pratio.BlockHash {
					result=false
				} else if txData.DataTx.Fee>10 {
					result=false
				}
			}
		}
		if poh.Right!=nil {
			txData=TdDeserialize(GetBytes(poh.Right.Val))
			if txData.DataType=="POW" {
				if txData.DataTx.Amount-txData.DataTx.Fee!=Pratio.BlockHash {
					result=false
				} else if txData.DataTx.Fee>10 {
					result=false
				}
			}
		}
	}
	if result {
		phash:=p.PowCubeHash(&iBlock,"0")
		if iBlock.PrevCubeHash!=phash {
			result=false
		}
	}
	return result
}

func ReadBlockVerifySet(index int,cubeno int) (string,string,int) {
	var iBlock Block
	var hash string
	err:=BlockRead(index,cubeno,&iBlock)
	Err(err,0)
	hash=iBlock.Hash
	pstr:=BlockVerifyString(iBlock)
	nonce:=iBlock.Nonce
	return hash,pstr,nonce
}

*/