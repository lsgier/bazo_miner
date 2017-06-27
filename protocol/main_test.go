package protocol

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"testing"
)

var accA, accB, minerAcc *Account
var PrivKeyA, PrivKeyB ecdsa.PrivateKey
var PubKeyA, PubKeyB ecdsa.PublicKey
var RootPrivKey ecdsa.PrivateKey
var MinerHash [32]byte
var MinerPrivKey *ecdsa.PrivateKey

func TestMain(m *testing.M) {

	addTestingAccounts()
	addRootAccounts()

	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func addTestingAccounts() {

	accA, accB, minerAcc = new(Account), new(Account), new(Account)

	puba1, _ := new(big.Int).SetString(PubA1, 16)
	puba2, _ := new(big.Int).SetString(PubA2, 16)
	priva, _ := new(big.Int).SetString(PrivA, 16)
	PubKeyA = ecdsa.PublicKey{
		elliptic.P256(),
		puba1,
		puba2,
	}
	PrivKeyA = ecdsa.PrivateKey{
		PubKeyA,
		priva,
	}

	pubb1, _ := new(big.Int).SetString(PubB1, 16)
	pubb2, _ := new(big.Int).SetString(PubB2, 16)
	privb, _ := new(big.Int).SetString(PrivB, 16)
	PubKeyB = ecdsa.PublicKey{
		elliptic.P256(),
		pubb1,
		pubb2,
	}
	PrivKeyB = ecdsa.PrivateKey{
		PubKeyB,
		privb,
	}

	copy(accA.Address[0:32], PrivKeyA.PublicKey.X.Bytes())
	copy(accA.Address[32:64], PrivKeyA.PublicKey.Y.Bytes())
	accAHash := serializeHashContent(accA.Address)

	//This one is just for testing purposes
	copy(accB.Address[0:32], PrivKeyB.PublicKey.X.Bytes())
	copy(accB.Address[32:64], PrivKeyB.PublicKey.Y.Bytes())
	accBHash := serializeHashContent(accB.Address)

	//just to bootstrap
	var shortHashA [8]byte
	var shortHashB [8]byte
	copy(shortHashA[:], accAHash[0:8])
	copy(shortHashB[:], accBHash[0:8])

	MinerPrivKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	var pubKey [64]byte
	var shortMiner [8]byte
	copy(pubKey[:32], MinerPrivKey.X.Bytes())
	copy(pubKey[32:], MinerPrivKey.Y.Bytes())
	MinerHash = serializeHashContent(pubKey)
	copy(shortMiner[:], MinerHash[0:8])
	minerAcc.Address = pubKey
}

func addRootAccounts() {

	var pubKey [64]byte

	pub1, _ := new(big.Int).SetString(RootPub1, 16)
	pub2, _ := new(big.Int).SetString(RootPub2, 16)
	priv, _ := new(big.Int).SetString(RootPriv, 16)
	PubKeyA = ecdsa.PublicKey{
		elliptic.P256(),
		pub1,
		pub2,
	}
	RootPrivKey = ecdsa.PrivateKey{
		PubKeyA,
		priv,
	}

	copy(pubKey[32-len(pub1.Bytes()):32], pub1.Bytes())
	copy(pubKey[64-len(pub2.Bytes()):], pub2.Bytes())

	rootHash := serializeHashContent(pubKey)

	var shortRootHash [8]byte
	copy(shortRootHash[:], rootHash[0:8])
}
