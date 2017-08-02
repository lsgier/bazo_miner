package miner

import (
	"github.com/lisgie/bazo_miner/storage"
	"testing"
)

func TestGetBlockSequences(t *testing.T) {

	cleanAndPrepare()

	b := newBlock([32]byte{})
	createBlockWithTxs(b)
	finalizeBlock(b)
	validateBlock(b)

	b2 := newBlock(b.Hash)
	createBlockWithTxs(b2)
	finalizeBlock(b2)
	validateBlock(b2)

	b3 := newBlock(b2.Hash)
	createBlockWithTxs(b3)
	finalizeBlock(b3)

	rollback, validate := getBlockSequences(b3)

	if len(rollback) != 0 {
		t.Error("Rollback shouldn't execute here\n")
	}

	if len(validate) != 1 || validate[0].Hash != b3.Hash {
		t.Error("Wrong validation sequence\n")
	}

	c := newBlock([32]byte{})
	createBlockWithTxs(c)
	finalizeBlock(c)
	storage.WriteOpenBlock(c)

	c2 := newBlock(c.Hash)
	createBlockWithTxs(c2)
	finalizeBlock(c2)
	storage.WriteOpenBlock(c2)

	c3 := newBlock(c2.Hash)
	createBlockWithTxs(c3)
	finalizeBlock(c3)

	//Blockchain now: genesis <- b <- b2
	//New Blockchain of longer size: genesis <- c <- c2 <- c3
	rollback, validate = getBlockSequences(c3)

	//Rollback slice needs to include b2 and b (in that order)
	if len(rollback) != 2 ||
		rollback[0].Hash != b2.Hash ||
		rollback[1].Hash != b.Hash {
		t.Error("Rollback of current chain failed\n")
	}

	if len(validate) != 3 ||
		validate[0].Hash != c.Hash ||
		validate[1].Hash != c2.Hash ||
		validate[2].Hash != c3.Hash {
		t.Error("Validation failed\n")
	}

	cleanAndPrepare()
	//Make sure that another chain of equal length does not get activated

	b = newBlock([32]byte{})
	createBlockWithTxs(b)
	finalizeBlock(b)
	validateBlock(b)

	b2 = newBlock(b.Hash)
	createBlockWithTxs(b2)
	finalizeBlock(b2)
	validateBlock(b2)

	b3 = newBlock(b2.Hash)
	createBlockWithTxs(b3)
	finalizeBlock(b3)
	validateBlock(b3)

	//Blockchain now: genesis <- b <- b2 <- b3
	//Competing chain: genesis <- c <- c2 <- c3
	c = newBlock([32]byte{})
	createBlockWithTxs(c)
	finalizeBlock(c)
	storage.WriteOpenBlock(c)

	c2 = newBlock(c.Hash)
	createBlockWithTxs(c2)
	finalizeBlock(c2)
	storage.WriteOpenBlock(c2)

	c3 = newBlock(c2.Hash)
	createBlockWithTxs(c3)
	finalizeBlock(c3)

	//Make sure that the new blockchain of equal length does not get activated
	rollback, validate = getBlockSequences(c3)
	if rollback != nil || validate != nil {
		t.Error("Did not properly detect longest chain\n")
	}
}

func TestGetNewChain(t *testing.T) {

	cleanAndPrepare()
	b := newBlock([32]byte{})
	createBlockWithTxs(b)
	finalizeBlock(b)
	validateBlock(b)

	b2 := newBlock(b.Hash)
	createBlockWithTxs(b2)
	finalizeBlock(b2)

	ancestor, newChain := getNewChain(b2)

	if ancestor.Hash != b.Hash {
		t.Errorf("Hash mismatch: %x vs. %x\n", ancestor.Hash, b.Hash)
	}
	if len(newChain) != 1 || newChain[0].Hash != b2.Hash {
		t.Error("Wrong new chain\n")
	}

	//Blockchain now: genesis <- b <- b2
	//New chain: genesis <- c <- c2
	c := newBlock([32]byte{})
	createBlockWithTxs(c)
	finalizeBlock(c)
	storage.WriteOpenBlock(c)

	c2 := newBlock(c.Hash)
	createBlockWithTxs(c2)
	finalizeBlock(c2)

	ancestor, newChain = getNewChain(c2)

	if ancestor.Hash != [32]byte{} {
		t.Errorf("Hash mismatch")
	}

	if len(newChain) != 2 || newChain[0].Hash != c.Hash || newChain[1].Hash != c2.Hash {
		t.Error("Wrong new chain\n")
	}
}
