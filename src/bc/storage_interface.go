package bc

import (
	"storage"
)

//acts as the interface to the storage module
func readBlock(hash [32]byte) (b *Block) {

	encodedBlock := storage.ReadBlock(hash)
	if encodedBlock == nil {
		return nil
	}
	return decodeBlock(encodedBlock)
}

func writeBlock(b *Block) {

	storage.WriteBlock(b.Hash,encodeBlock(b))
}

func readOpenFundsTx(hash [32]byte) (tx *fundsTx) {

	encodedTx := storage.ReadOpenTx(hash)
	if encodedTx == nil {
		return nil
	}
	decodedTx := DecodeFundsTx(encodedTx)
	//this isn't very nice, but necessary to enrich the full hashes in the transaction
	decodedTx.verify()
	return decodedTx

}

func readClosedFundsTx(hash [32]byte) (tx *fundsTx) {

	encodedTx := storage.ReadClosedTx(hash)
	if encodedTx == nil {
		return nil
	}
	decodedTx := DecodeFundsTx(encodedTx)
	decodedTx.verify()
	return decodedTx
}

func readOpenAccTx(hash [32]byte) (tx *accTx) {

	encodedTx := storage.ReadOpenTx(hash)
	if encodedTx == nil {
		return nil
	}
	return DecodeAccTx(encodedTx)
}

func readClosedAccTx(hash [32]byte) (tx *accTx) {

	encodedTx := storage.ReadClosedTx(hash)
	if encodedTx == nil {
		return nil
	}
	return DecodeAccTx(encodedTx)
}

func writeOpenFundsTx(tx *fundsTx) {

	storage.WriteOpenTx(hashFundsTx(tx),EncodeFundsTx(tx))
}

func writeClosedFundsTx(tx *fundsTx) {

	storage.WriteClosedTx(hashFundsTx(tx),EncodeFundsTx(tx))
}

func writeOpenAccTx(tx *accTx) {

	storage.WriteOpenTx(hashAccTx(tx),EncodeAccTx(tx))
}

func writeClosedAccTx(tx *accTx) {

	storage.WriteClosedTx(hashAccTx(tx),EncodeAccTx(tx))
}

func deleteOpenFundsTx(hash [32]byte) {

	storage.WriteOpenTx(hash,nil)
}

//delete in the closed confirmation is needed as well, in case of block rollback
func deleteClosedFundsTx(hash [32]byte) {

	storage.WriteClosedTx(hash, nil)
}

func deleteOpenAccTx(hash [32]byte) {

	storage.WriteOpenTx(hash,nil)
}

func deleteClosedAccTx(hash [32]byte) {

	storage.WriteClosedTx(hash, nil)
}
