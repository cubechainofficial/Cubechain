package core

import (
	"fmt"
	"crypto/sha256"
	"encoding/base64"
)

type BST struct {
	Root	*BSTNode
}

type BSTNode struct {
	Hash	string
	Val		interface{}
	Left	*BSTNode
	Right	*BSTNode
}

type depthNode struct {
	depth	int
	node	*BSTNode
	lr		string
}

func AddBST(v interface{}) *BST {
	tree:=&BST{}
	h:=setHashV(v)
	tree.Root=&BSTNode{Hash:h,Val:v}
	return tree
}

func (t *BST) Print() {
	q:=[]depthNode{}
	q=append(q,depthNode{depth:0,node:t.Root,lr:"Root"})
	currentDepth:=0
	for len(q)>0 {
		var first depthNode
		first,q=q[0],q[1:]
		if first.depth!=currentDepth {
			fmt.Println()
			currentDepth=first.depth
		}
		fmt.Println("[",first.depth," : ",first.lr,"] ",first.node.Hash," ",first.node.Val)
		if first.node.Left!=nil {
			q=append(q,depthNode{depth:currentDepth+1,node:first.node.Left,lr:"L"})
		}
		if first.node.Right!=nil {
			q=append(q,depthNode{depth:currentDepth+1,node:first.node.Right,lr:"R"})
		}
	}
}

func (t *BST) Merkle() string {
	var txhash []string
	q:=[]depthNode{}
	q=append(q,depthNode{depth:0,node:t.Root})
	currentDepth:=0
	for len(q)>0 {
		var first depthNode
		first,q=q[0],q[1:]
		if first.depth!=currentDepth {
			currentDepth=first.depth
		}
		txhash=append(txhash,first.node.Hash)
		if first.node.Left!=nil {
			q=append(q,depthNode{depth:currentDepth+1,node:first.node.Left,lr:"L"})
		}
		if first.node.Right!=nil {
			q=append(q,depthNode{depth:currentDepth+1,node:first.node.Right,lr:"R"})
		}
	}
	mhash:=MerkleHash(txhash)
	return mhash
}

func MerkleHash(elements []string) string {
	if len(elements) == 0 {
		return ""
	} else if len(elements) == 1 {
		return elements[0]
	}
	half := len(elements) / 2
	a := MerkleHash(elements[:half])
	b := MerkleHash(elements[half:])
	return MerkleElement(a, b)
}

func MerkleElement(a, b string) string {
	combined := a + b
	hash := sha256.New()
	hash.Write([]byte(combined))
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return sha
}

func (t *BST) Search(h string) (interface{},int) {
	return t.Root.Search(h, 1)
}

func (t *BST) Convert00(RtxArr *[]TxData) {
	t.Root.Convert(RtxArr)
}

func (t *BST) Convert(RtxArr *[]TxData) {
	q:=[]depthNode{}
	q=append(q,depthNode{depth:0,node:t.Root,lr:"Root"})
	currentDepth:=0
	for len(q)>0 {
		var first depthNode
		first,q=q[0],q[1:]
		if first.depth!=currentDepth {
			currentDepth=first.depth
		}
		if t.Root!=first.node {
			tData,_:=first.node.Val.(*TxData)
			*RtxArr=append(*RtxArr,*tData)
		}
		if first.node.Left!=nil {
			q=append(q,depthNode{depth:currentDepth+1,node:first.node.Left,lr:"L"})
		}
		if first.node.Right!=nil {
			q=append(q,depthNode{depth:currentDepth+1,node:first.node.Right,lr:"R"})
		}
	}
}


func (n *BSTNode) AddNode(v interface{},h string) *BSTNode {
	if h=="" {
		h=setHashV(v) 
	}
	if n.Hash==h {
		echo("Already exists Hash :",h," => ",v)
		return nil
	} else if n.Hash>h {
		if n.Left==nil {
			n.Left=&BSTNode{Hash:h,Val:v}
			return n.Left
		} else {
			return n.Left.AddNode(v,h)
		}
	} else {
		if n.Right==nil {
			n.Right=&BSTNode{Hash:h,Val:v}
			return n.Right
		} else {
			return n.Right.AddNode(v,h)
		}
	}
}

func (n *BSTNode) AddNodeHash(v interface{},h string) (*BSTNode,string) {
	if h=="" {
		h=setHashV(v) 
	}
	an:=n.AddNode(v,h)
	if an==nil {
		return nil,""
	} else {
		return an,h
	}
}


func (n *BSTNode) Search(h string,cnt int) (interface{},int) {
	if n.Hash==h {
		return n.Val,cnt
	} else if n.Hash>h {
		if n.Left!=nil {
			return n.Left.Search(h, cnt+1)
		}
		return nil,cnt
	} else {
		if n.Right!=nil {
			return n.Right.Search(h,cnt+1)
		}
		return nil,cnt
	}
}

func (n *BSTNode) Convert(txArr *[]TxData) {
	if n.Left!=nil {
		LeftData,_:=n.Left.Val.(TxData)
		*txArr=append(*txArr,LeftData)
		n.Left.Convert(txArr)
	} else if n.Right!=nil {
		RightData,_:=n.Right.Val.(TxData)
		*txArr=append(*txArr,RightData)
		n.Right.Convert(txArr)
	}
}

func formatString(arg interface{}) string {
	switch arg.(type) {
	case TxData:   
		p := arg.(TxData) 
		return p.String()
	case *TxData:              
		p := arg.(*TxData)
		return p.String()
	default:
		return "Error"
	}
}
