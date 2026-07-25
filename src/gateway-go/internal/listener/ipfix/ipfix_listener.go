package ipfix

import (
	"context"
	"log"
	"net"

	"gateway-go/internal/listener"
)

type IPFIXListener struct {
	conn *net.UDPConn
}

func NewIPFIXListener() *IPFIXListener {
	return &IPFIXListener{}
}

func (l *IPFIXListener) Name() string {
	return "IPFIX-NetFlow-Listener"
}

func (l *IPFIXListener) Start(ctx context.Context, dataCh chan<- listener.TelemetryData) error {
	addr := net.UDPAddr {
		Port: 2055,
		IP: net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}
	l.conn = conn

	log.Printf("[%s] UDP socket listening on 0.0.0.0:2055 for IPFIX/NerFlow packets...", l.Name())

	go func() {
		defer conn.Close()

		go func() {
			<-ctx.Done()
			conn.Close()
		}()

		buf := make([]byte, 2048)
		for {
			n, srcAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				select {
				case <-ctx.Done():
					log.Printf("[%s] Stopping IPFIX Listener...", l.Name())
					return
				default:
					log.Printf("[%s] UDP Read error: %v", l.Name(), err)
					continue
				}
			}

			dataCh <- listener.TelemetryData {
				Timestamp: 0,
				Source: "ipfix",
				Metric: "netflow_packet_bytes",
				Value: float64(n),
				Tags: map[string]string {
					"exporter_ip": srcAddr.IP.String(),
				},
			}
		}
	}()

	return nil
}

func (l *IPFIXListener) Stop() error {
	if l.conn != nil {
		_ = l.conn.Close()
	}

	return nil
}

