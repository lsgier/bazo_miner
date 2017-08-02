package storage

import (
	"github.com/boltdb/bolt"
	"github.com/lisgie/bazo_miner/protocol"
)

func ReadOpenBlock(hash [32]byte) (block *protocol.Block) {

	var encodedBlock []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("openblocks"))
		encodedBlock = b.Get(hash[:])
		return nil
	})

	if encodedBlock == nil {
		return nil
	}

	return block.Decode(encodedBlock)
}

func ReadClosedBlock(hash [32]byte) (block *protocol.Block) {

	var encodedBlock []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("closedblocks"))
		encodedBlock = b.Get(hash[:])
		return nil
	})

	if encodedBlock == nil {
		return nil
	}

	return block.Decode(encodedBlock)
}

func ReadOpenTx(hash [32]byte) (transaction protocol.Transaction) {

	return txMemPool[hash]
}

//needed for the miner to prepare a new block
func ReadAllOpenTxs() (allOpenTxs []protocol.Transaction) {

	for key := range txMemPool {
		allOpenTxs = append(allOpenTxs, txMemPool[key])
	}

	return
}

func ReadClosedTx(hash [32]byte) (transaction protocol.Transaction) {

	var encodedTx []byte
	var fundstx *protocol.FundsTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("closedfunds"))
		encodedTx = b.Get(hash[:])
		return nil
	})
	if encodedTx != nil {
		return fundstx.Decode(encodedTx)
	}

	var acctx *protocol.AccTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("closedaccs"))
		encodedTx = b.Get(hash[:])
		return nil
	})
	if encodedTx != nil {
		return acctx.Decode(encodedTx)
	}

	var configtx *protocol.ConfigTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("closedconfigs"))
		encodedTx = b.Get(hash[:])
		return nil
	})
	if encodedTx != nil {
		return configtx.Decode(encodedTx)
	}
	return nil
}
