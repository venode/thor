package consensus

import (
	"github.com/vechain/thor/block"
	"github.com/vechain/thor/chain"
	"github.com/vechain/thor/tx"
)

type validator struct {
	block *block.Block
	chain *chain.Chain
}

func newValidator(blk *block.Block, chain *chain.Chain) *validator {
	return &validator{
		block: blk,
		chain: chain}
}

func (v *validator) validate() (*block.Header, error) {
	preHeader, err := v.chain.GetBlockHeader(v.block.ParentHash())
	if err != nil {
		if v.chain.IsNotFound(err) {
			return nil, errParentNotFound
		}
		return nil, err
	}

	if preHeader.Timestamp() >= v.block.Timestamp() {
		return nil, errTimestamp
	}

	header := v.block.Header()

	if header.TxsRoot() != v.block.Body().Txs.RootHash() {
		return nil, errTxsRoot
	}

	if header.GasUsed() > header.GasLimit() {
		return nil, errGasUsed
	}

	for _, transaction := range v.block.Transactions() {
		if !v.validateTransaction(transaction) {
			return nil, errTransaction
		}
	}

	return preHeader, nil
}

func (v *validator) validateTransaction(transaction *tx.Transaction) bool {
	if len(transaction.Clauses()) == 0 {
		return false
	}

	if transaction.TimeBarrier() > v.block.Timestamp() {
		return false
	}

	_, err := v.chain.LookupTransaction(v.block.ParentHash(), transaction.Hash())

	return v.chain.IsNotFound(err)
}
