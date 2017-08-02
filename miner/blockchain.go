package miner

import (
	"fmt"
	"github.com/lisgie/bazo_miner/protocol"
	"github.com/lisgie/bazo_miner/storage"
	"log"
	"math/big"
	"os"
	"sync"
	"time"
)

var (
	logger               *log.Logger
	accA, accB, minerAcc *protocol.Account
	hashA, hashB         [32]byte
	blockValidation      = &sync.Mutex{}
	timestamp            []int64
	parameterSlice       []parameters
	activeParameters     *parameters
	tmpSlice             []parameters
	uptodate             bool
)

func Init() {

	//Initialize root key
	initRootKey()

	LogFile, _ := os.OpenFile("log/miner "+time.Now().String(), os.O_RDWR|os.O_CREATE, 0666)
	logger = log.New(LogFile, "", log.LstdFlags)

	//var tmpTimestamp []int64
	parameterSlice = append(parameterSlice, parameters{
		[32]byte{},
		1,
		1000,
		2016,
		60,
		0,
	})
	activeParameters = &parameterSlice[0]

	currentTargetTime = new(timerange)
	target = append(target, 14)

	localBlockCount = -1
	globalBlockCount = -1
	genesis := newBlock([32]byte{})
	collectStatistics(genesis)
	storage.WriteClosedBlock(genesis)

	logger.Println("Starting system, initializing state map")

	go incomingData()
	mining()
}

func mining() {
	currentBlock := newBlock([32]byte{})
	for {
		err := finalizeBlock(currentBlock)
		if err != nil {
			fmt.Printf("Mining failure: %v\n", err)
		} else {
			fmt.Println("Block mined.")
		}
		//else a block was received meanwhile that was added to the chain, all the effort was in vain :(
		//wait for lock here only
		if err != nil {
			logger.Printf("%v\n", err)
		} else {
			broadcastBlock(currentBlock)
			validateBlock(currentBlock)
		}

		//mining successful, construct new block out of mempool transactions
		blockValidation.Lock()
		nextBlock := newBlock(lastBlock.Hash)
		currentBlock = nextBlock
		prepareBlock(currentBlock)
		blockValidation.Unlock()
	}
}

func initRootKey() {

	var pubKey [64]byte

	pub1, _ := new(big.Int).SetString(INITROOTKEY1, 16)
	pub2, _ := new(big.Int).SetString(INITROOTKEY2, 16)

	copy(pubKey[:32], pub1.Bytes())
	copy(pubKey[32:], pub2.Bytes())

	rootHash := serializeHashContent(pubKey)

	rootAcc := protocol.Account{Address: pubKey}

	storage.State[rootHash] = &rootAcc
	storage.RootKeys[rootHash] = &rootAcc
}
