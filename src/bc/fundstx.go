package bc

import (

	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
	"reflect"
	"fmt"
	"encoding/binary"
)

//when we broadcast transactions we need a way to distinguish with a type

type fundsTx struct {
	Header byte
	Amount [4]byte
	Fee [2]byte
	TxCnt [3]byte
	From [8]byte
	fromHash [32]byte
	To [8]byte
	toHash [32]byte
	Xored [24]byte
	Sig [40]byte
}

func constrFundsTx(header byte, amount [4]byte, fee [2]byte, txCnt [3]byte, from, to [32]byte, key *ecdsa.PrivateKey) (tx fundsTx, err error) {

	//avoid sending money to its own acc, doesn't make sense with account-based money
	txToHash := struct {
		Header byte
		Amount [4]byte
		Fee [2]byte
		TxCnt [3]byte
		From [32]byte
		To [32]byte
	} {
		header,
		amount,
		fee,
		txCnt,
		from,
		to,
	}

	sigHash := serializeHashContent(txToHash)

	fmt.Printf("%x\n", sigHash)

	r,s, err := ecdsa.Sign(rand.Reader, key, sigHash[:])

	var sig [64]byte
	copy(sig[:32],r.Bytes())
	copy(sig[32:],s.Bytes())

	tx.Header = header
	tx.Amount = amount
	tx.Fee = fee
	tx.TxCnt = txCnt

	copy(tx.From[0:8],from[0:8])
	copy(tx.To[0:8],to[0:8])

	for i := 0; i < 24; i++ {
		tx.Xored[i] = from[i+8] ^ to[i+8] ^ sig[i]
	}

	copy(tx.Sig[:],sig[24:64])

	return
}

//I believe sender balance check here is a bad idea. This limits to receive and send funds within the same block
//But if receiving and sending along funds within the same block, transaction ordering is important
func (tx *fundsTx) verify() bool {

	var sig [24]byte
	var concatSig [64]byte
	pub1,pub2 := new(big.Int), new(big.Int)
	r,s := new(big.Int), new(big.Int)

	//fundstx only makes sense if amount > 0
	if binary.BigEndian.Uint32(tx.Amount[:]) == 0 {
		return false
	}

	//check if accounts are present in the actual state
	for _,accFrom := range State[tx.From] {
		for _,accTo := range State[tx.To] {
			sig = [24]byte{}
			for cnt := 0; cnt < 24; cnt++ {
				sig[cnt] = tx.Xored[cnt] ^ accFrom.Hash[cnt+8] ^ accTo.Hash[cnt+8]
			}
			copy(concatSig[:24],sig[0:24])
			copy(concatSig[24:],tx.Sig[:])

			pub1.SetBytes(accFrom.Address[:32])
			pub2.SetBytes(accFrom.Address[32:])

			r.SetBytes(concatSig[:32])
			s.SetBytes(concatSig[32:])

			txHash := struct {
				Header byte
				Amount [4]byte
				Fee [2]byte
				TxCnt [3]byte
				From [32]byte
				To [32]byte
			} {
				tx.Header,
				tx.Amount,
				tx.Fee,
				tx.TxCnt,
				accFrom.Hash,
				accTo.Hash,
			}
			sigHash := serializeHashContent(txHash)

			pubKey := ecdsa.PublicKey{elliptic.P256(), pub1, pub2}
			if ecdsa.Verify(&pubKey,sigHash[:],r,s) == true && !reflect.DeepEqual(accFrom,accTo) {
				tx.fromHash = accFrom.Hash
				tx.toHash = accTo.Hash
				return true
			}
		}
	}

	return false
}


/*Header byte
Amount [4]byte
Fee [2]byte
TxCnt [3]byte
From [8]byte
fromHash [32]byte
To [8]byte
toHash [32]byte
Xored [24]byte
Sig [40]byte*/

//when we serialize the struct with binary.Write, unexported field get serialized as well, undesired
//behavior. Therefore, writing own encoder/decoder
func encodeFundsTx(tx fundsTx) (encodedTx []byte) {
	encodedTx = make([]byte,90)

	fmt.Printf("%x\n", tx.Fee[:])
	encodedTx[0] = tx.Header
	copy(encodedTx[1:5], tx.Amount[:])
	copy(encodedTx[5:7], tx.Fee[:])
	copy(encodedTx[7:10], tx.TxCnt[:])
	copy(encodedTx[10:18], tx.From[:])
	copy(encodedTx[18:26], tx.To[:])
	copy(encodedTx[26:50], tx.Xored[:])
	copy(encodedTx[50:90], tx.Sig[:])

	return encodedTx
}

func decodeFundsTx(encodedTx []byte) (tx *fundsTx) {
	tx = new(fundsTx)
	tx.Header = encodedTx[0]
	copy(tx.Amount[:], encodedTx[1:5])
	copy(tx.Fee[:], encodedTx[5:7])
	copy(tx.TxCnt[:], encodedTx[7:10])
	copy(tx.From[:], encodedTx[10:18])
	copy(tx.To[:], encodedTx[18:26])
	copy(tx.Xored[:], encodedTx[26:50])
	copy(tx.Sig[:], encodedTx[50:90])

	return tx
}


func (tx fundsTx) String() string {
	return fmt.Sprintf(
		"Header: %x\n" +
			"Amount: %v\n" +
			"Fee: %v\n" +
			"TxCnt: %v\n" +
			"From: %x\n" +
			"From Full Hash: %x\n" +
			"To: %x\n" +
			"To Full Hash: %x\n" +
			"Xored: %x\n" +
			"Sig: %x\n",
		tx.Header,
		tx.Amount,
		tx.Fee,
		tx.TxCnt,
		tx.From,
		tx.fromHash[0:12],
		tx.To,
		tx.toHash[0:12],
		tx.Xored[0:8],
		tx.Sig[0:8],
	)
}