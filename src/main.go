package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"bc"
	"fmt"
)

func main() {

	state := make(map[[64]byte]int64)
	privA, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privB, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return
	}

	accA := bc.Account{Nonce:0, Balance:15}
	copy(accA.Id[0:32], privA.PublicKey.X.Bytes())
	copy(accA.Id[32:64], privA.PublicKey.Y.Bytes())

	accB := bc.Account{Nonce:0, Balance:12}
	copy(accB.Id[0:32], privB.PublicKey.X.Bytes())
	copy(accB.Id[32:64], privB.PublicKey.Y.Bytes())

	state[accA.Id] = accA.Balance
	state[accB.Id] = accB.Balance

	b := bc.NewBlock([32]byte{}, state)

	tx, err := bc.ConstrTx(0, 2, accA, accB, privA)
	tx2, err := bc.ConstrTx(0, 3, accB, accA, privB)
	tx3, err := bc.ConstrTx(0, 1, accA, accB, privA)
	tx4, err := bc.ConstrTx(0, 4, accB, accA, privB)
	tx5, err := bc.ConstrTx(0, 3, accA, accB, privA)
	tx6, err := bc.ConstrTx(0, 1, accB, accA, privB)


	b.AddTx(&tx)
	b.AddTx(&tx2)
	b.AddTx(&tx3)
	b.AddTx(&tx4)
	b.AddTx(&tx5)
	b.AddTx(&tx6)

	if err != nil {
		return
	}

	b.FinalizeBlock()

	fmt.Printf("%x\n", b)

	/*var buf bytes.Buffer
	var rcvTx bc.Transaction
	enc := gob.NewEncoder(&buf)
	enc.Encode(tx)
	fmt.Printf("%d\n", len(buf.Bytes()))
	dec := gob.NewDecoder(&buf)
	dec.Decode(&rcvTx)
	fmt.Printf("%x\n", rcvTx)*/
}

