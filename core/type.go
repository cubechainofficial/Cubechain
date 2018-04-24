package core

type CubeBlock struct {
	Index		int
	Timestamp	int
	Cube		[27]Block
	Chash       string
}

type Block struct {
	Index			int
	Cubeno			int
	Timestamp		int
	Data			[]byte
	Hash			string
	PrevHash		string
	PrevCubeHash	string
	Nonce			int
}

type TransactionData struct {
	Timestamp	int
	From		[]byte
	To			[]byte
	Amount		int
	Hash		[]byte
	Sign		[]byte
	Nonce		int
}

type TxData struct {
	DataType	string
	DataTx		TransactionData
}

type TxPool struct {
	Tdata	[]TxData
}

type IBlock struct {
	IndexAddress	[]IndexType
}

type IndexType struct {
	Address		string
	Indexing	[]CubeIndex
}

type CubeIndex struct {
	Index	int
	CubeNum	int
}

type StaticRule1 struct {
	Address string
	Balance	int
}

type SBlock struct {
	RuleArr	[]StaticRule1
}

type EscrowData struct {
	EscrowTx	TransactionData
	EscrowType  int
	EscrowKey	string
	EscrowTime	int 
	State		int
}

type OpenRule struct {
	HashId	    string
	Key2nd		string
	State		int    
}

type EBlock struct {
	Escrow		[]EscrowData 
}

type TxPoolTemp struct {
	Tdata	[]TxData
}

type POH struct {
	Index		int
	Cubeno		int
	Hash		string
	Type		int
	Address		string
	Amount		int
	State		int
}

type Pohr struct {
	BlockHash	int
	CHash		int
	Cubing		int
	POS			int
}

type PosWallet struct {
	Address		string
	Amount		int
}