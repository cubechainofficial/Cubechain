package core

import (
	"fmt"
	"bytes"
	"encoding/gob"
)

type TxBST struct {
	Coin		*BST
	Poh			*BST
	Token		*BST
	Data		*BST
	Contract	*BST
}

type tx interface {	txDefine() }
type txTree interface { BstNode() }

type txCoin struct {
}
type txPoh struct {
}
type txToken struct {
}
type txData struct {
}
type txContract struct {
}

func (td *txCoin) txDefine() {}
func (td *txPoh) txDefine() {}
func (td *txToken) txDefine() {}
func (td *txData) txDefine() {}
func (td *txContract) txDefine() {}

func (td *txCoin) BstNode(txb TxBST) {
	txtype:="Coin"
	txb.treeInsert(td,txtype)
	return
}

func (td *txPoh) BstNode(txb TxBST) {
	txtype:="Poh"
	txb.treeInsert(td,txtype)
	return
}

func (td *txToken) BstNode(txb TxBST) {
	txtype:="Token"
	txb.treeInsert(td,txtype)
	return
}

func (td *txData) BstNode(txb TxBST) {
	txtype:="Data"
	txb.treeInsert(td,txtype)
	return
}

func (td *txContract) BstNode(txb TxBST) {
	txtype:="Contract"
	txb.treeInsert(td,txtype)
	return
}

func (txbst *TxBST) Init() {
	if txbst.Coin==nil || txbst.Poh==nil || txbst.Token==nil || txbst.Data==nil || txbst.Contract==nil {
		txbst.Coin=AddBST("Coin")
		txbst.Poh=AddBST("Poh")
		txbst.Token=AddBST("Token")
		txbst.Data=AddBST("Data")
		txbst.Contract=AddBST("Contract")
	}
}

func (txbst *TxBST) treeInsert(treeData interface{},treeType string) {
	txbst.treeInsertNode(treeData,treeType)
}


func (txbst *TxBST) treeInsertNode(treeData interface{},treeType string) *BSTNode {
	hash:=""
	if treeType!="Data" &&  treeType!="Contract" {
		hash=treeData.(TxData).Hash
	}
	n:=&BSTNode{}
	if txbst.Coin==nil || txbst.Poh==nil || txbst.Token==nil || txbst.Data==nil || txbst.Contract==nil {
		txbst.Coin=AddBST("Coin")
		txbst.Poh=AddBST("Poh")
		txbst.Token=AddBST("Token")
		txbst.Data=AddBST("Data")
		txbst.Contract=AddBST("Contract")
	}
	switch treeType {
		case "Coin": n=txbst.Coin.Root.AddNode(treeData,hash)
		case "Poh": n=txbst.Poh.Root.AddNode(treeData,hash)
		case "Token": n=txbst.Token.Root.AddNode(treeData,hash)
		case "Data": n=txbst.Data.Root.AddNode(treeData,hash)
		case "Contract": n=txbst.Contract.Root.AddNode(treeData,hash)
		default : fmt.Println("Please check tx type.")
	}
	return n
}

func (txbst *TxBST) treeInsertHash(treeData interface{},treeType string) string {
	hash:=""
	if treeType!="Data" &&  treeType!="Contract" {
		hash=treeData.(TxData).Hash
	}
	h:=""
	if txbst.Coin==nil || txbst.Poh==nil || txbst.Token==nil || txbst.Data==nil || txbst.Contract==nil {
		txbst.Coin=AddBST("Coin")
		txbst.Poh=AddBST("Poh")
		txbst.Token=AddBST("Token")
		txbst.Data=AddBST("Data")
		txbst.Contract=AddBST("Contract")
	}
	switch treeType {
		case "Coin": _,h=txbst.Coin.Root.AddNodeHash(treeData,hash)
		case "Poh": _,h=txbst.Poh.Root.AddNodeHash(treeData,hash)
		case "Token": _,h=txbst.Token.Root.AddNodeHash(treeData,hash)
		case "Data": _,h=txbst.Data.Root.AddNodeHash(treeData,hash)
		case "Contract": _,h=txbst.Contract.Root.AddNodeHash(treeData,hash)
		default : fmt.Println("Please check tx type.")
	}
	return h
}

func BlockTree(BlockData []byte) TxBST{
	var tbst TxBST
	tbst=TreeDeserialize(BlockData)
	return tbst
}

func BlockBST(BlockData []byte,treeType string) BST {
	var result BST
	tbst:=BlockTree(BlockData)
	switch treeType {
		case "Coin": result=*tbst.Coin
		case "Poh": result=*tbst.Poh
		case "Token": result=*tbst.Token
		case "Data": result=*tbst.Data
		case "Contract": result=*tbst.Contract
		default : fmt.Println("Please check tx type.")
	}
	return result
}

func TreeToData(tbst TxBST) []TxData {
	var txbData []TxData
	if tbst.Poh!=nil {
		tbst.Poh.Convert(&txbData)
	}
	if tbst.Coin!=nil {
		tbst.Coin.Convert(&txbData)
	}
	if tbst.Token!=nil {
		tbst.Token.Convert(&txbData)
	}
	return txbData
}

func TreeToSpecial(tbst TxBST) map[string]string {
	var spData map[string]string
	if tbst.Data!=nil {
		spData,_=tbst.Data.Root.Val.(map[string]string)
	}
	return spData
}

func BlockTxData(BlockData []byte) []TxData{
	var tbst TxBST
	tbst=TreeDeserialize(BlockData)
	return TreeToData(tbst)
}

func BlockSpecialData(BlockData []byte) map[string]string{
	var tbst TxBST
	tbst=TreeDeserialize(BlockData)
	return TreeToSpecial(tbst)
}

func BlockTreeRoot(BlockData []byte,treeType string) BSTNode {
	var result BSTNode
	tbst:=BlockTree(BlockData)
	switch treeType {
		case "Coin": result=*tbst.Coin.Root
		case "Poh": result=*tbst.Poh.Root
		case "Token": result=*tbst.Token.Root
		case "Data": result=*tbst.Data.Root
		case "Contract": result=*tbst.Contract.Root
		default : fmt.Println("Please check tx type.")
	}
	return result
}

func (tb *TxBST) Insert(treeData interface{}) {
}


func TreeDeserialize(data []byte) TxBST {
	var object TxBST
	decoder:=gob.NewDecoder(bytes.NewReader(data))
	err:=decoder.Decode(&object)
	if err != nil {
		decho(err)
	}
	return object
}

func TreePrint() {
	var txb TxBST
	txb.treeInsert("tx1","Coin")
	txb.treeInsert("tx2","Coin")
	txb.treeInsert("tx3","Coin")
	txb.treeInsert([]string{"tx4","2234"},"Coin")
	txb.treeInsert("Token1","Token")
	txb.treeInsert("Token2","Token")
	txb.treeInsert("Token3","Token")

	txb.treeInsert("Data1","Data")
	txb.treeInsert("Data1","Data")
	txb.treeInsert("datererwwa1","Data")

	txb.treeInsert("Contract111","Contract")
	txb.treeInsert("Contract222","Contract")
	txb.treeInsert("Contract333","Contract")

	txb.treeInsert("Contract","ee")

	txb.TreePrint2()

	str,c:=txb.Coin.Search("453d7ae0dc1e658cf52e81e3488132025c0ff610da47c178a735ab09f058e8ea")

	decho(str)
	decho(c)
}

func (txbst TxBST) TreePrint2() {
	fmt.Println("")
	if txbst.Coin!=nil {
		fmt.Println("Print Coin.")
		txbst.Coin.Print()
	}
	if txbst.Poh!=nil {
		fmt.Println("Print Poh.")
		txbst.Poh.Print()
	}
	if txbst.Token!=nil {
		fmt.Println("")
		fmt.Println("Print Token.")
		txbst.Token.Print()
	}
	if txbst.Data!=nil {
		fmt.Println("")
		fmt.Println("Print Data.")
		txbst.Data.Print()
	}
	if txbst.Contract!=nil {
		fmt.Println("")
		fmt.Println("Print Contract.")
		txbst.Contract.Print()
	}
}



