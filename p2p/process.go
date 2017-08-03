package p2p

import (
	"encoding/binary"
	"strconv"
	"github.com/lisgie/bazo_miner/protocol"
	"github.com/lisgie/bazo_miner/storage"
)

//Process tx broadcasts from other miners We can't broadcast incoming messages directly, first check if
//the tx has already been broadcast before, whether it is a valid tx etc.
func processTxBrdcst(p *peer, payload []byte, brdcstType uint8) {

	var tx protocol.Transaction

	//Make sure the transaction can be properly decoded, verification is done at a later stage to reduce latency
	switch brdcstType {
	case FUNDSTX_BRDCST:
		var fTx *protocol.FundsTx
		fTx = fTx.Decode(payload)
		if fTx == nil {
			return
		}
		tx = fTx
	case ACCTX_BRDCST:
		var aTx *protocol.AccTx
		aTx = aTx.Decode(payload)
		if aTx == nil {
			return
		}
		tx = aTx
	case CONFIGTX_BRDCST:
		var cTx *protocol.ConfigTx
		cTx = cTx.Decode(payload)
		if cTx == nil {
			return
		}
		tx = cTx
	}
	if storage.ReadOpenTx(tx.Hash()) != nil {
		logger.Printf("Received transaction (%x) already in the mempool.\n", tx.Hash())
		return
	}
	if storage.ReadClosedTx(tx.Hash()) != nil {
		logger.Printf("Received transaction (%x) already validated.\n", tx.Hash())
		return
	}

	//Write to mempool and rebroadcast
	logger.Printf("Writing transaction (%x) in the mempool.\n", tx.Hash())
	storage.WriteOpenTx(tx)
	toBrdcst := BuildPacket(brdcstType, payload)
	brdcstMsg <- toBrdcst
}

func processTimeRes(p *peer, payload []byte) {

	time := int64(binary.BigEndian.Uint64(payload))
	//concurrent writes need to be protected
	//we use the same lock to prevent concurrent writes. It would be more efficient to use different locks
	//but the speedup is so marginal that it's not worth it
	p.l.Lock()
	defer p.l.Unlock()
	p.time = time
}

func processNeighborRes(p *peer, payload []byte) {

	//parse the incoming ipv4 addresses
	ipportList := _processNeighborRes(payload)

	for _, ipportIter := range ipportList {
		logger.Printf("IP/Port received: %v\n", ipportIter)
		iplistChan <- ipportIter
	}
}

//Decoupled for cleaner testing
func _processNeighborRes(payload []byte) (ipportList []string) {

	index := 0
	for cnt := 0; cnt < len(payload)/(IPV4ADDR_SIZE+PORT_SIZE); cnt++ {
		var addr string
		for singleAddr := index; singleAddr < index+IPV4ADDR_SIZE; singleAddr++ {
			tmp := int(payload[singleAddr])
			addr += strconv.Itoa(tmp)
			addr += "."
		}
		//remove trailing dot
		addr = addr[:len(addr)-1]
		addr += ":"
		//extract port number
		addr += strconv.Itoa(int(binary.BigEndian.Uint16(payload[index+4 : index+6])))

		//add ipaddr to the channel
		ipportList = append(ipportList, addr)
		index += IPV4ADDR_SIZE + PORT_SIZE
	}
	return ipportList
}

