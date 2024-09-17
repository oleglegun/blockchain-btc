package node

type UTXO struct {
	Hash string
	// OutIndex is an index of the output in the transaction
	OutIndex int
	Amount   int64
	// Every UTXO is considered “unspent” until it is used as an input in a new transaction.
	// Once it is used, it is no longer a valid UTXO. The blockchain tracks all UTXOs to know what funds are available to be spent.
	IsSpent bool
}

func NewUTXO(hash string, outIndex int, amount int64) *UTXO {
	isSpent := false

	return &UTXO{
		Hash:     hash,
		OutIndex: outIndex,
		Amount:   amount,
		IsSpent:  isSpent,
	}
}
