package core

import (
	"fmt"
	"strconv"
	"time"
)

func cfileName(cblock *CubeBlock) string {
	filename:=strconv.Itoa(cblock.Index) + "_" + cblock.Chash + ".cub"
	return filename	
}

func cblockFile(cblock CubeBlock) error {
	filename:=cfileName(&cblock)
	err := fileWrite2(filename, cblock)
	return err
}

func genesisCube() CubeBlock {
	var clen=CurrentHeight();
	var block Block
	var sblock [27]Block
	var cblock CubeBlock
	var chash string
	if clen>1 { 
		fmt.Printf("Invalid\n")
		return 	cblock
	}
	for i:=0;i<27;i++ {
		block=genesisBlock(i)
		sblock[i]=block
		chash+=block.Hash
	}
	cblock=CubeBlock{0,0,sblock,sethash2(chash)}
	cblock.Timestamp=int(time.Now().Unix())
	cblockFile(cblock)
	return cblock
}


