package transport

import (
	"vos"
)

type TransportUdp struct {
	listeningPoints []*net.UDPConn
	stop            bool
}

// event from tansport to upper

func (this *TransportUdp) Run() {
	// add local listen addresses
	localUdpAddrList := []string{"10.40.20.3:8000", "10.40.20.4:7000"}

	// start listen
	err := this.StartListen(localUdpAddrList)
	if err != nil {
		return
	}
}

func (this *TransportUdp) StartListen(localUdpAddrList []string) error {
	for _, v := range localUdpAddrList {
		err := this.Listen(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *TransportUdp) Listen(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	lp, err := net.ListenUDP("udp", addr)

	if err == nil {
		this.listeningPoints = append(udp.listeningPoints, lp)
		go this.listen(lp)
	}

	return err
}

func (this *TransportUdp) listen(conn *net.UDPConn) {
	buffer := make([]byte, 65535)
	for {
		num, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if this.stop {
				break
			} else {
				continue
			}
		}

		pkt := append([]byte(nil), buffer[:num]...)
		// TODO: send msg to upper layer
	}
}

func (this *TransportUdp) Stop() {
	this.stop = true
	for _, lp := range udp.listeningPoints {
		lp.Close()
	}
}
