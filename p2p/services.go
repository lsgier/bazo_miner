package p2p

import (
	"time"
)

//this is not accessed concurrently, one single goroutine
func broadcastService() {
	for {
		select {
		//broadcasting all messages
		//Mutex for peers structure need not be claimed here, because
		//this is the only function that can actually add or reject connections (no race conditions
		case msg := <-brdcstMsg:
			for p := range peers.peerConns {
				p.ch <- msg
			}
		case p := <-register:
			peers.add(p)
			//peers.peerConns[p] = true
		case p := <-disconnect:
			peers.delete(p)
			//delete(peers.peerConns, p)
			close(p.ch)
		}
	}
}

//Belongs to the broadcast service
func peerBroadcast(p *peer) {
	for msg := range p.ch {
		sendData(p, msg)
	}
}

//Single goroutine that makes sure the system is well connected
func checkHealthService() {

	for {
		//time.Sleep(time.Minute)
		if len(peers.peerConns) >= MIN_MINERS {
			time.Sleep(2 * HEALTH_CHECK_INTERVAL * time.Minute)
			continue
		} else {
			//this delay is needed to prevent sending neighbor requests like a maniac
			time.Sleep(HEALTH_CHECK_INTERVAL * time.Second)
		}

	RETRY:
		select {
		case ipaddr := <-iplistChan:
			p, err := initiateNewMinerConnection(ipaddr)
			if p == nil || err != nil {
				goto RETRY
			}
			go minerConn(p)
			break
		default:
			neighborReq()
			break
		}
	}
}

func timeService() {

	//initialize system time
	systemTime = time.Now().Unix()
	go func() {
		for {
			time.Sleep(UPDATE_SYS_TIME * time.Second)
			writeSystemTime()
		}
	}()

	for {
		time.Sleep(TIME_BRDCST_INTERVAL * time.Second)
		packet := BuildPacket(TIME_BRDCST, getTime())
		brdcstMsg <- packet
	}
}
